package schedulers

import (
	"fmt"
	"log"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	"ottopoint-purchase/hosts/opl/host"
	opl "ottopoint-purchase/hosts/opl/host"
	sepulsa "ottopoint-purchase/hosts/sepulsa/host"
	signature "ottopoint-purchase/hosts/signature/host"
	voucherAg "ottopoint-purchase/hosts/voucher_aggregator/host"
	voucherModel "ottopoint-purchase/hosts/voucher_aggregator/models"
	"ottopoint-purchase/models"
	"ottopoint-purchase/models/dbmodels"
	"ottopoint-purchase/utils"
	"reflect"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/vjeantet/jodaTime"
)

func (t SchedulerCheckStatusService) CheckStatusSepulsaServices(trxid string) error {
	// res := models.SchedulerCheckStatusData{}

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info(">>> [Start]-[CheckStatusSepulsaServices] <<<")

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[CheckStatusSepulsaServices]")
	defer span.Finish()

	fmt.Println(">>> [Start]-[CheckStatusSepulsaServices] <<<")

	// check status ke sepulsa
	_, errStatus := sepulsa.EwalletDetailTransaction(trxid)
	if errStatus != nil {

		fmt.Println(fmt.Sprintf("[Error from EwalletDetailTransaction]-[Error : %v]", errStatus))
		fmt.Println("[PackageServices]-[EwalletDetailTransaction]")

		sugarLogger.Info(fmt.Sprintf("[Error from EwalletDetailTransaction]-[Error : %v]", errStatus))
		sugarLogger.Info("[PackageServices]-[EwalletDetailTransaction]")

		return errStatus
	}

	return nil

}

func (t SchedulerCheckStatusService) CheckStatusVoucherAgService(trxID string) error {

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info(">>> [Start]-[CheckStatusVoucherAgService] <<<")

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[CheckStatusVoucherAgService]")
	defer span.Finish()

	fmt.Println(">>> [Start]-[CheckStatusVoucherAgService] <<<")

	// Get TSpending By OrderID
	spendings, err := db.GetVoucherAgSpendingSecond(trxID)
	if err != nil {
		fmt.Println(fmt.Sprintf("[Error from VoucherAgTransaction]-[Error : %v]", err))
		fmt.Println("[PackageServices]-[VoucherAgTransaction]")

		sugarLogger.Info(fmt.Sprintf("[Error from VoucherAgTransaction]-[Error : %v]", err))
		sugarLogger.Info("[PackageServices]-[VoucherAgTransaction]")
	}

	head := models.RequestHeader{
		InstitutionID: "PSM0001",
		DeviceID:      "12341414",
		Geolocation:   "45453452, 25235235",
		ChannelID:     "H2H",
		AppsID:        "422432432435",
		Timestamp:     "1579666534",
	}

	count := len(spendings)
	voucherReq := voucherModel.RequestCheckOrderStatus{
		OrderID:       trxID,
		RecordPerPage: fmt.Sprintf("%d", count),
		CurrentPage:   "1",
	}
	log.Println("VOUCHER AGGREGATOR: ", voucherReq)
	sign, err := signature.Signature(voucherReq, head)
	if err != nil {
		fmt.Println(fmt.Sprintf("[Error from VoucherAgTransaction]-[Error : %v]", err))
		fmt.Println("[PackageServices]-[VoucherAgTransaction]")

		sugarLogger.Info(fmt.Sprintf("[Error from VoucherAgTransaction]-[Error : %v]", err))
		sugarLogger.Info("[PackageServices]-[VoucherAgTransaction]")
	}

	s := reflect.ValueOf(sign.Data)
	for _, k := range s.MapKeys() {
		head.Signature = fmt.Sprintf("%s", s.MapIndex(k))
	}

	// Get Order Status Voucher Aggregator
	orderStatus, errStatus := voucherAg.CheckStatusOrder(voucherReq, head)
	if errStatus != nil {

		fmt.Println(fmt.Sprintf("[Error from VoucherAgTransaction]-[Error : %v]", errStatus))
		fmt.Println("[PackageServices]-[VoucherAgTransaction]")

		sugarLogger.Info(fmt.Sprintf("[Error from VoucherAgTransaction]-[Error : %v]", errStatus))
		sugarLogger.Info("[PackageServices]-[VoucherAgTransaction]")

		return errStatus
	}

	if orderStatus != nil {

		if orderStatus.ResponseCode == "09" ||
			orderStatus.ResponseCode == "01" {

			spending := spendings[0]
			totalPoint := int(spending.Amount) * count
			trxID := utils.GenTransactionId()
			text := trxID + spending.Institution + constants.CodeReversal + "#" + "OP009 - Reversal point cause transaction " + spending.Voucher + " is failed"

			schedulerData := dbmodels.TSchedulerRetry{
				Code:          constants.CodeScheduler,
				TransactionID: utils.Before(text, "#"),
				Count:         0,
				IsDone:        false,
				CreatedAT:     time.Now(),
			}

			reversal, errReversal := opl.TransferPoint(spending.AccountId, fmt.Sprint(totalPoint), text)

			statusEarning := constants.Success
			msgEarning := constants.MsgSuccess

			if errReversal != nil || reversal.PointsTransferId == "" {

				statusEarning = constants.TimeOut
				statusEarning = constants.TimeOut

				for _, val1 := range reversal.Form.Children.Customer.Errors {
					if val1 != "" {
						msgEarning = val1
						statusEarning = constants.Failed
					}
				}

				for _, val2 := range reversal.Form.Children.Points.Errors {
					if val2 != "" {
						msgEarning = val2
						statusEarning = constants.Failed
					}
				}

				if reversal.Message != "" {
					msgEarning = reversal.Message
					statusEarning = constants.Failed
				}

				if reversal.Error.Message != "" {
					msgEarning = reversal.Error.Message
					statusEarning = constants.Failed
				}

				if statusEarning == constants.TimeOut {
					errSaveScheduler := db.DbCon.Create(&schedulerData).Error
					if errSaveScheduler != nil {

						fmt.Println("===== Gagal SaveScheduler ke DB =====")
						fmt.Println(fmt.Sprintf("Error : %v", errSaveScheduler))
						fmt.Println(fmt.Sprintf("===== Phone : %v || RRN : %v =====", spending.AccountNumber, spending.RRN))

					}

				}

			}

			expired := t.ExpiredPointService()

			saveReversal := dbmodels.TEarning{
				ID:               utils.GenerateTokenUUID(),
				PartnerId:        spending.Institution,
				TransactionId:    trxID,
				AccountNumber:    spending.AccountNumber,
				Point:            int64(totalPoint),
				Status:           statusEarning,
				StatusMessage:    msgEarning,
				PointsTransferId: reversal.PointsTransferId,
				TransType:        constants.CodeReversal,
				AccountId:        spending.AccountId,
				ExpiredPoint:     expired,
				TransactionTime:  time.Now(),
			}

			errSaveReversal := db.DbCon.Create(&saveReversal).Error
			if errSaveReversal != nil {
				fmt.Println(fmt.Sprintf("[Failed Save Reversal to DB]-[Error : %v]", errSaveReversal))
				fmt.Println("[PackageServices]-[SaveEarning]")

				fmt.Println(">>> Save CSV <<<")
				name := jodaTime.Format("dd-MM-YYYY", time.Now()) + ".csv"
				go utils.CreateCSVFile(saveReversal, name)
			}

		}

		for _, v := range spendings {
			go db.UpdateVoucherAgSecond(orderStatus.ResponseDesc, orderStatus.ResponseCode, v.ID)
		}

	}

	return nil

}

func (t SchedulerCheckStatusService) ExpiredPointService() string {

	fmt.Println(">>> ExpiredPointService <<<")

	get, err := host.SettingsOPL()
	if err != nil || get.Settings.ProgramName == "" {

		fmt.Println(fmt.Sprintf("[Error : %v]", err))
		fmt.Println("[PackageBulkService]-[ExpiredPointService]")

	}

	data := get.Settings.PointsDaysActiveCount + 1

	expired := utils.FormatTimeString(time.Now(), 0, 0, data)

	return expired

}
