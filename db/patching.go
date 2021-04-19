package db

import (
	"fmt"
	"ottopoint-purchase/models/dbmodels"

	"github.com/sirupsen/logrus"
)

func GeDataPatching() ([]dbmodels.TSpending, error) {
	res := []dbmodels.TSpending{}

	err := DbCon.Raw(`select * from t_spending where invoice_number = '' and created_at between '2021-04-16 00:00:00' and '2021-04-19 00:00:00'`).Scan(&res).Error
	if err != nil {
		logrus.Error("[PackageDB]-[GeDataPatching]")
		logrus.Error(fmt.Sprintf("[Error : %v]", err))

		return res, err
	}

	return res, nil
}

func UpdateDataPatching(invoiceNum, id string) error {
	res := dbmodels.TSpending{}

	err := DbCon.Raw(`update t_spending set invoice_number = ? where id = ?`, invoiceNum, id).Scan(&res).Error
	strerror := fmt.Sprintf("%v", err)
	if strerror != "record not found" {
		logrus.Error("[PackageDB]-[UpdateDataPatching]")
		logrus.Error(fmt.Sprintf("[Error : %v]-[TrxId : %v]", err, id))

		return err
	}

	return nil
}
