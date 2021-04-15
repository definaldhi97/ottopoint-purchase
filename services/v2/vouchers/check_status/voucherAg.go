package check_status

import (
	"encoding/json"
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	kafka "ottopoint-purchase/hosts/publisher/host"
	signature "ottopoint-purchase/hosts/signature/host"
	voucherAg "ottopoint-purchase/hosts/voucher_aggregator/host"
	voucherModel "ottopoint-purchase/hosts/voucher_aggregator/models"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"
	"reflect"

	"ottopoint-purchase/services/v2.1/Trx"

	"github.com/sirupsen/logrus"
)

func CheckStatusVoucherAgService(trxID string) error {

	nameservice := "[PackageCheckStatus]-[CheckStatusVoucherAgService]"

	logrus.Info(nameservice)

	// Get TSpending By OrderID
	spendings, err := db.GetVoucherAgSpendingSecond(trxID)
	if err != nil {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[GetVoucherAgSpendingSecond]-[Error : %v]", err))
		logrus.Println("TransactionID : ", trxID)
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

	logrus.Info("VOUCHER AGGREGATOR: ", voucherReq)

	sign, errSign := signature.Signature(voucherReq, head)
	if errSign != nil {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[Signature]-[Error : %v]", errSign))
		logrus.Println("TransactionID : ", trxID)

	}

	s := reflect.ValueOf(sign.Data)
	for _, k := range s.MapKeys() {
		head.Signature = fmt.Sprintf("%s", s.MapIndex(k))
	}

	// Get Order Status Voucher Aggregator
	orderStatus, errStatus := voucherAg.CheckStatusOrder(voucherReq, head)
	if errStatus != nil {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[CheckStatusOrder]-[Error : %v]", errStatus))
		logrus.Println("TransactionID : ", trxID)

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
				TrxID:         utils.GenTransactionId(),
				AccountId:     spending.AccountId,
				RewardID:      *spending.MRewardID,
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

			_, errKafka := kafka.SendPublishKafka(kafkaReq)
			if errKafka != nil {

				logrus.Error(nameservice)
				logrus.Error(fmt.Sprintf("[SendPublishKafka]-[Error : %v]", errKafka))
				logrus.Println("TransactionID : ", trxID)

			}

		}

		for _, v := range spendings {
			go db.UpdateVoucherAgSecond(orderStatus.ResponseDesc, orderStatus.ResponseCode, v.ID)
		}

	}

	return nil

}
