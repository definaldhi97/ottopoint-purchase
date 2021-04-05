package check_status

import (
	"encoding/json"
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	kafka "ottopoint-purchase/hosts/publisher/host"
	sepulsa "ottopoint-purchase/hosts/sepulsa/host"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"

	"ottopoint-purchase/services/v2.1/Trx"

	"github.com/sirupsen/logrus"
)

func CheckStatusSepulsaServices(trxid string) error {
	// res := models.SchedulerCheckStatusData{}

	nameservice := "[PackageCheckStatus]-[CheckStatusSepulsaServices]"

	logrus.Info(nameservice)

	// check status ke sepulsa
	resp, errStatus := sepulsa.EwalletDetailTransaction(trxid)
	if errStatus != nil {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[EwalletDetailTransaction]-[Error : %v]", errStatus))
		logrus.Println("TransactionID : ", trxid)

		return errStatus
	}

	// Get Spending By TransactionID and OrderID
	spending, err := db.GetSpendingSepulsa(resp.TransactionID, resp.OrderID)
	if err != nil {
		fmt.Println("[GetSpendingSepulsa] : ", err.Error())
		logrus.Error("[ Failed Get SpendingSepulsa ] : ", err.Error())
	}

	responseCode := models.GetErrorMsg(resp.ResponseCode)

	logrus.Info("[HandleSchedulerSepulsa] - [ResponseCode] : ", resp.ResponseCode)
	logrus.Info("[HandleSchedulerSepulsa] - [ResponseDesc] : ", responseCode)

	param := models.Params{
		InstitutionID: spending.Institution,
		NamaVoucher:   spending.Voucher,
		AccountId:     spending.AccountId,
		AccountNumber: spending.AccountNumber,
		RRN:           spending.RRN,
		TrxID:         utils.GenTransactionId(),
		RewardID:      *spending.MRewardID,
		Point:         spending.Point,
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

	if (responseCode != "Success") && (responseCode != "Pending") {

		resultReversal := Trx.V21_Adding_PointVoucher(param, spending.Point, 1, header)
		fmt.Println(resultReversal)

		fmt.Println("[ >>>>>>>>>>>>>>>>>>>>>>> Send Publisher <<<<<<<<<<<<<<<<<<<< ]")

		pubreq := models.NotifPubreq{
			Type:           constants.CODE_REVERSAL_POINT,
			NotificationTo: spending.AccountNumber,
			Institution:    spending.Institution,
			ReferenceId:    spending.RRN,
			TransactionId:  spending.TransactionId,
			Data: models.DataValue{
				RewardValue: "point",
				Value:       fmt.Sprint(spending.Point),
			},
		}

		bytePub, _ := json.Marshal(pubreq)

		kafkaReq := kafka.PublishReq{
			Topic: constants.TOPIC_PUSHNOTIF_GENERAL,
			Value: bytePub,
		}

		kafkaRes, err := kafka.SendPublishKafka(kafkaReq)
		if err != nil {
			logrus.Error("Gagal Send Publisher : ", err)
		}
		logrus.Info("[ Response Publisher ] : ", kafkaRes)

	}

	responseSepulsa, _ := json.Marshal(resp)

	// Update TSpending
	go db.UpdateVoucherSepulsa(responseCode, resp.ResponseCode, string(responseSepulsa), resp.TransactionID, resp.OrderID)

	return nil

}
