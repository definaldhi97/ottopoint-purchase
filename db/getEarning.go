package db

import (
	"fmt"
	"ottopoint-purchase/models/dbmodels"

	"github.com/astaxie/beego/logs"
)

func GetEarningCode(code string) (dbmodels.MEarningRule, error) {
	res := dbmodels.MEarningRule{}

	err := DbCon.Where("code = ?", code).First(&res).Error
	if err != nil {
		logs.Info("Failed to Checking from database", err)
		return res, err
	}
	logs.Info("Data MEarning :", res)

	return res, nil
}

func GetEarningCodebyProductCode(productCode string) (dbmodels.MEarningRule, error) {
	res := dbmodels.MEarningRule{}

	fmt.Println("[Select from GeneralSpending]")
	err := DbCon.Raw(`select * from m_earning_rule where lower(sku_ids) like '%?%' = ?`, productCode).Scan(&res).Error
	if err != nil {

		fmt.Println(fmt.Sprintf("[Failed to Get EarningCode from GeneralSpending]-[Error : %v]", err.Error()))
		fmt.Println("[PackageDB]-[GetEarningCodebyProductCode]")

		fmt.Println("[Select from CustoomeEventRule]")
		err = DbCon.Raw(`select * from m_earning_rule where event_name = ?`, productCode).Scan(&res).Error
		if err != nil {

			fmt.Println(fmt.Sprintf("[Failed to Get EarningCode from CustoomeEventRule]-[Error : %v]", err.Error()))
			fmt.Println("[PackageDB]-[GetEarningCodebyProductCode]")

			return res, err
		}

	}

	return res, nil
}

func GetCheckStatusEarning(reff string, institution string) (dbmodels.TEarning, error) {
	res := dbmodels.TEarning{}

	err := DbCon.Where("reference_id = ? and partner_id = ?", reff, institution).First(&res).Error
	if err != nil {
		logs.Info("Failed to get GetCheckStatusEarning from database", err)
		return res, err
	}
	logs.Info("Data GetCheckStatusEarning :", res)

	return res, nil
}
