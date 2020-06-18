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

func GetCheckStatusEarning(reff, institution string) (dbmodels.TEarning, error) {
	res := dbmodels.TEarning{}

	err := DbCon.Where("reference_id = ? and m_institution_id = ?", reff, institution).First(&res).Error
	if err != nil {
		logs.Info("Failed to get GetCheckStatusEarning from database", err)
		return res, err
	}
	logs.Info("Data GetCheckStatusEarning :", res)

	return res, nil
}
