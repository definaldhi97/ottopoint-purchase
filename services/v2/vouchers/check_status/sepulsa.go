package check_status

import (
	"encoding/json"
	"fmt"
	"ottopoint-purchase/db"
	sepulsa "ottopoint-purchase/hosts/sepulsa/host"
	"ottopoint-purchase/models"

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
	logrus.Infof("[DetailTransactionSpending] %v\n", spending)

	if (responseCode != "Success") && (responseCode != "Pending") {

	}

	responseSepulsa, _ := json.Marshal(resp)

	// Update TSpending
	go db.UpdateVoucherSepulsa(responseCode, resp.ResponseCode, string(responseSepulsa), resp.TransactionID, resp.OrderID)

	return nil

}
