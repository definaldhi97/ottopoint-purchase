package db

import (
	"fmt"
	"ottopoint-purchase/models/dbmodels"
)

func SaveEarning(earning dbmodels.TEarning) error {

	err := DbCon.Create(&earning).Error
	if err != nil {
		fmt.Println(fmt.Sprintf("[Error : %v]", err))
		fmt.Println("[Failed SaveEarning to DB]")
		fmt.Println("[Package DB]-[SaveEarning]")
		return err

	}

	return nil

}
