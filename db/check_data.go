package db

import (
	"fmt"
	"ottopoint-purchase/models/dbmodels"

	"github.com/sirupsen/logrus"
)

func CheckTrxbyTrxID(trxId string) error {
	res := dbmodels.TSpending{}

	err := DbCon.Raw(`select * from public.t_spending where cummulative_ref = ?`, trxId).Scan(&res).Error
	if err != nil {

		logrus.Error("[PackageDB]-[CheckTrxbyTrxID]")
		logrus.Error(fmt.Sprintf("[Failed get Data by TrxID : %v from TSpending]-[Error : %v]", trxId, err))

		return err
	}

	return nil
}
