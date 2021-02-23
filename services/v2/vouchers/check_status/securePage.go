package check_status

import (
	"fmt"
	lp "ottopoint-purchase/hosts/landing_page/host"

	"github.com/sirupsen/logrus"
)

func CheckStatusSecurePageServices(trxid string) error {

	nameservice := "[PackageCheckStatus]-[CheckStatusSepulsaServices]"

	logrus.Info(nameservice)

	// check status ke sepulsa
	_, errStatus := lp.CheckStatusLandingPage(trxid)
	if errStatus != nil {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[EwalletDetailTransaction]-[Error : %v]", errStatus))
		logrus.Println("TransactionID : ", trxid)

		return errStatus
	}

	return nil

}
