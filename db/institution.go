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

type InstitutionKey struct {
	Apikey string `gorm:"apikey"`
	PubKey string `gorm:"pub_key"`
}

func GetInstitutionKey(institution string) (InstitutionKey, error) {

	res := InstitutionKey{}

	err := DbCon.Raw(`select apikey, pub_key from m_institution_key as a join m_institution as b on a.institution_id = b.partner_id where b.partner_id = ?`, institution).Scan(&res).Error
	if err != nil {

		fmt.Println(fmt.Sprintf("[GetIdInstitution]-[Error : %v]", err))
		fmt.Println("[GetIdInstitution]-[Failed GetIdInstitution]")

		return res, err
	}

	fmt.Println("GetInstitutionKey : ", res)

	return res, nil
}
