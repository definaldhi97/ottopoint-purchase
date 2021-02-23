package db

import (
	"fmt"
	"ottopoint-purchase/models/dbmodels"

	"github.com/sirupsen/logrus"
)

func CheckUser(phone string) (dbmodels.User, error) {
	res := dbmodels.User{}

	// err := Dbcon.Exec(`select * from users where phone = ?, status = true`, phone).Scan(&res).Error
	err := DbCon.Where("phone = ? and status = true", phone).First(&res).Error
	if err != nil {
		logrus.Info("Failed to Checking from database", err)
		return res, err
	}
	logrus.Info("data eligible :", res)

	return res, nil
}

func UserWithInstitution(phone, institution string) (dbmodels.User, error) {
	res := dbmodels.User{}

	err := DbCon.Raw(`select * from users as a join users_link as b on a.id = b.user_id where a.phone = ? and a.status = true and b.institution_id = ? and b.is_link = true`, phone, institution).Scan(&res).Error
	if err != nil {

		fmt.Println(fmt.Sprintf("[Failed get User from DB][Error : %v]", err))
		fmt.Println("[PackageDB][UserWithInstitution]")

		return res, err
	}
	logrus.Info("data eligible :", res)

	return res, nil
}
