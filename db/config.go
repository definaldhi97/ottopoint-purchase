package db

import (
	"fmt"
	"ottopoint-purchase/models/dbmodels"
)

func GetConfig() (dbmodels.Configs, error) {
	res := dbmodels.Configs{}

	err := DbCon.Raw(`SELECT * FROM public.configs`).Scan(&res).Error
	if err != nil {

		fmt.Println("[EEROR-DATABASE]")
		fmt.Println("[db]")
		fmt.Println("[GetConfig]")
		fmt.Println(fmt.Sprintf("Failed to connect database transaction %v", err))

		return res, err
	}

	return res, nil
}
