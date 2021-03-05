package db

import (
	"fmt"
	"ottopoint-purchase/models/dbmodels"

	"github.com/sirupsen/logrus"
)

func GetEarningCode(code string) (dbmodels.MEarningRule, error) {
	res := dbmodels.MEarningRule{}

	err := DbCon.Where("code = ?", code).First(&res).Error
	if err != nil {
		logrus.Info("Failed to Checking from database", err)
		return res, err
	}
	logrus.Info("Data MEarning :", res)

	return res, nil
}

func GetEarningCodebyProductCode(productCode string) (dbmodels.MEarningRule, error) {
	res := dbmodels.MEarningRule{}

	fmt.Println("[Select from GeneralSpending]")

	include := "%" + productCode + "%"
	exclude := "%" + productCode + "%"
	// query := fmt.Sprintf("select * from m_earning_rule where code like '%GSR%' and included_skus like %v or excluded_skus like %v", include, exclude, productCode, productCode)

	err := DbCon.Raw(`select * from m_earning_rule where code like '%GSR%' and (included_skus like ? or excluded_skus like ? and active = true)`, include, exclude).Scan(&res).Error
	if err != nil {

		fmt.Println("[PackageDB]-[GetEarningCodebyProductCode]")
		fmt.Println(fmt.Sprintf("[Failed to Get EarningCode from GeneralSpending]-[Error : %v]", err.Error()))

		fmt.Println("[Select from CustoomeEventRule]")
		err = DbCon.Raw(`select * from m_earning_rule where code like '%CER%' and event_name = ? and active = true`, productCode).Scan(&res).Error
		if err != nil {

			fmt.Println("[PackageDB]-[GetEarningCodebyProductCode]")
			fmt.Println(fmt.Sprintf("[Failed to Get EarningCode from CustoomeEventRule]-[Error : %v]", err.Error()))

			return res, err
		}

	}

	return res, nil
}

func GetCheckStatusEarning(reff string, institution string) (dbmodels.TEarning, error) {
	res := dbmodels.TEarning{}

	err := DbCon.Where("reference_id = ? and partner_id = ?", reff, institution).First(&res).Error
	if err != nil {
		logrus.Info("Failed to get GetCheckStatusEarning from database", err)
		return res, err
	}
	logrus.Info("Data GetCheckStatusEarning :", res)

	return res, nil
}
