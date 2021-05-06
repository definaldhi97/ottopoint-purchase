package db

import (
	"fmt"
	"ottopoint-purchase/models/dbmodels"

	"github.com/sirupsen/logrus"
)

func CheckTrxbyCumReff(cumReff string) (dbmodels.TSpending, error) {
	res := dbmodels.TSpending{}

	err := DbCon.Where("cummulative_ref = ?", cumReff).First(&res).Error
	// err := DbCon.Raw(`select * from public.t_spending where cummulative_ref = ?`, trxId).Scan(&res).Error
	if err != nil {

		logrus.Error("[PackageDB]-[CheckTrxbyCumReff]")
		logrus.Error(fmt.Sprintf("[Failed get Data by CumReff : %v from TSpending]-[Error : %v]", cumReff, err))

		return res, err
	}

	return res, nil
}

func CheckTrxbyTrxID(trxId string) (dbmodels.TSpending, error) {
	res := dbmodels.TSpending{}

	err := DbCon.Where("transaction_id = ?", trxId).First(&res).Error
	if err != nil {

		logrus.Error("[PackageDB]-[CheckTrxbyTrxID]")
		logrus.Error(fmt.Sprintf("[Failed get Data by TrxID : %v from TSpending]-[Error : %v]", trxId, err))

		return res, err
	}

	return res, nil
}
