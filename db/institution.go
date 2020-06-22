package db

import (
	"fmt"
	"ottopoint-purchase/models/dbmodels"
)

func GetIdInstitution(name string) (dbmodels.MInstution, error) {
	res := dbmodels.MInstution{}

	fmt.Println("[PackageDB]-[GetIdInstitution]")

	err := DbCon.Raw(`select * from m_institution where partner_id = ?`, name).Scan(&res).Error
	if err != nil {

		fmt.Println(fmt.Sprintf("[GetIdInstitution]-[Error : %v]", err))
		fmt.Println("[GetIdInstitution]-[Failed GetIdInstitution]")

		return res, err
	}
	fmt.Println("response GetIdInstitution :", res)

	return res, nil
}
