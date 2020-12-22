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
	"ottopoint-purchase/services/v2.1/Trx"
	"ottopoint-purchase/utils"
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
		Timestamp:     utils.GetTimeFormatMillisecond(),
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
				RewardID:      spending.MRewardID,
			}

			header := models.RequestHeader{
				DeviceID:      "ottopoint-purchase",
				InstitutionID: spending.Institution,
				Geolocation:   "-",
				ChannelID:     "H2H",
				AppsID:        "-",
				Timestamp:     utils.GetTimeFormatYYMMDDHHMMSS(),
				Authorization: "-",
				Signature:     "-",
			}

			resultReversal := Trx.V21_Adding_PointVoucher(param, totalPoint, count, header)
			logrus.Info(resultReversal)

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
