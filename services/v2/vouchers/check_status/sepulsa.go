package check_status

import (
	"fmt"
	sepulsa "ottopoint-purchase/hosts/sepulsa/host"

	"github.com/sirupsen/logrus"
)

func CheckStatusSepulsaServices(trxid string) error {
	// res := models.SchedulerCheckStatusData{}

	nameservice := "[PackageCheckStatus]-[CheckStatusSepulsaServices]"

	logrus.Info(nameservice)

	// check status ke sepulsa
	_, errStatus := sepulsa.EwalletDetailTransaction(trxid)
	if errStatus != nil {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[EwalletDetailTransaction]-[Error : %v]", errStatus))
		logrus.Println("TransactionID : ", trxid)

		return errStatus
	}

	return nil

}
