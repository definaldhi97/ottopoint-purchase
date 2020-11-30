package schedulers

import (
	"encoding/json"
	"fmt"
	"log"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	kafka "ottopoint-purchase/hosts/publisher/host"
	sepulsa "ottopoint-purchase/hosts/sepulsa/host"
	signature "ottopoint-purchase/hosts/signature/host"
	voucherAg "ottopoint-purchase/hosts/voucher_aggregator/host"
	voucherModel "ottopoint-purchase/hosts/voucher_aggregator/models"
	"ottopoint-purchase/models"
	"ottopoint-purchase/services/v2_migrate"
	"reflect"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
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
		InstitutionID: spendings[0].Institution,
		DeviceID:      "ottopoint-scheduler",
		Geolocation:   "-",
		ChannelID:     "H2H",
		AppsID:        "-",
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

		if orderStatus.ResponseCode == "09" || orderStatus.ResponseCode == "01" {

			spending := spendings[0]
			totalPoint := int(spending.Amount) * count
			// transactionID := utils.GenTransactionId()
			// text := trxID + spending.Institution + constants.CodeReversal + "#" + "OP009 - Reversal point cause transaction " + spending.Voucher + " is failed"

			param := models.Params{
				InstitutionID: spending.Institution,
				NamaVoucher:   spending.Voucher,
				AccountNumber: spending.AccountNumber,
				TrxID:         trxID,
				AccountId:     spending.AccountId,
				RewardID:      spending.RewardID,
			}

			header := models.RequestHeader{
				DeviceID:      "ottopoint-scheduler",
				InstitutionID: spending.Institution,
				Geolocation:   "-",
				ChannelID:     "H2H",
				AppsID:        "-",
				Timestamp:     "-",
				Authorization: "-",
				Signature:     "-",
			}

			resultReversal := v2_migrate.Adding_PointVoucher(param, totalPoint, count, param.TrxID, header)
			logrus.Info(resultReversal)

			// schedulerData := dbmodels.TSchedulerRetry{
			// 	Code:          constants.CodeScheduler,
			// 	TransactionID: utils.Before(text, "#"),
			// 	Count:         0,
			// 	IsDone:        false,
			// 	CreatedAT:     time.Now(),
			// }

			// reversal, errReversal := opl.TransferPoint(spending.AccountId, fmt.Sprint(totalPoint), text)

			// statusEarning := constants.Success
			// msgEarning := constants.MsgSuccess

			// if errReversal != nil || reversal.PointsTransferId == "" {

			// 	statusEarning = constants.TimeOut
			// 	statusEarning = constants.TimeOut

			// 	for _, val1 := range reversal.Form.Children.Customer.Errors {
			// 		if val1 != "" {
			// 			msgEarning = val1
			// 			statusEarning = constants.Failed
			// 		}
			// 	}

			// 	for _, val2 := range reversal.Form.Children.Points.Errors {
			// 		if val2 != "" {
			// 			msgEarning = val2
			// 			statusEarning = constants.Failed
			// 		}
			// 	}

			// 	if reversal.Message != "" {
			// 		msgEarning = reversal.Message
			// 		statusEarning = constants.Failed
			// 	}

			// 	if reversal.Error.Message != "" {
			// 		msgEarning = reversal.Error.Message
			// 		statusEarning = constants.Failed
			// 	}

			// 	if statusEarning == constants.TimeOut {
			// 		errSaveScheduler := db.DbCon.Create(&schedulerData).Error
			// 		if errSaveScheduler != nil {

			// 			fmt.Println("===== Gagal SaveScheduler ke DB =====")
			// 			fmt.Println(fmt.Sprintf("Error : %v", errSaveScheduler))
			// 			fmt.Println(fmt.Sprintf("===== Phone : %v || RRN : %v =====", spending.AccountNumber, spending.RRN))

			// 		}

			// 	}

			// }

			// expired := t.ExpiredPointService()

			// saveReversal := dbmodels.TEarning{
			// 	ID:               utils.GenerateTokenUUID(),
			// 	PartnerId:        spending.Institution,
			// 	TransactionId:    transactionID,
			// 	AccountNumber:    spending.AccountNumber,
			// 	Point:            int64(totalPoint),
			// 	Status:           statusEarning,
			// 	StatusMessage:    msgEarning,
			// 	PointsTransferId: reversal.PointsTransferId,
			// 	TransType:        constants.CodeReversal,
			// 	AccountId:        spending.AccountId,
			// 	ExpiredPoint:     expired,
			// 	TransactionTime:  time.Now(),
			// }

			// errSaveReversal := db.DbCon.Create(&saveReversal).Error
			// if errSaveReversal != nil {
			// 	fmt.Println(fmt.Sprintf("[Failed Save Reversal to DB]-[Error : %v]", errSaveReversal))
			// 	fmt.Println("[PackageServices]-[SaveEarning]")

			// 	fmt.Println(">>> Save CSV <<<")
			// 	name := jodaTime.Format("dd-MM-YYYY", time.Now()) + ".csv"
			// 	go utils.CreateCSVFile(saveReversal, name)
			// }

			fmt.Println("[ >>>>>>>>>>>>>>>>> Send Publisher Notification <<<<<<<<<<<<<<<< ]")
			pubreq := models.NotifPubreq{
				Type:           constants.CODE_REVERSAL_POINT,
				NotificationTo: spending.AccountNumber,
				Institution:    spending.Institution,
				ReferenceId:    spending.RRN,
				TransactionId:  trxID,
				Data: models.DataValue{
					RewardValue: "point",
					Value:       fmt.Sprint(totalPoint),
				},
			}

			bytePub, _ := json.Marshal(pubreq)
			kafkaReq := kafka.PublishReq{
				Topic: constants.TOPIC_PUSHNOTIF_GENERAL,
				Value: bytePub,
			}

			kafkaRes, err := kafka.SendPublishKafka(kafkaReq)
			if err != nil {
				fmt.Println("Gagal Send Publisher")
				fmt.Println("Error : ", err)
			}

			fmt.Println("Response Publisher : ", kafkaRes)

		}

		for _, v := range spendings {
			go db.UpdateVoucherAgSecond(orderStatus.ResponseDesc, orderStatus.ResponseCode, v.ID)
		}

	}

	return nil

}
