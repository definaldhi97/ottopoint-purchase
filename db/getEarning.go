package db

import (
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
