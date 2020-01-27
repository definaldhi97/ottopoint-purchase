package db

import (
	"ottopoint-purchase/models/dbmodels"

	"github.com/astaxie/beego/logs"
)

func CheckUser(phone string) (dbmodels.User, error) {
	res := dbmodels.User{}

	// err := Dbcon.Exec(`select * from users where phone = ?, status = true`, phone).Scan(&res).Error
	err := DbCon.Where("phone = ? and status = true", phone).First(&res).Error
	if err != nil {
		logs.Info("Failed to Checking from database", err)
		return res, err
	}
	logs.Info("data eligible :", res)

	return res, nil
}
