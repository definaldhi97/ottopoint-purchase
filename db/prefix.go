package db

import (
	"fmt"
	"ottopoint-purchase/models/dbmodels"
)

// GetOperatorCodebyPrefix ..
func GetOperatorCodebyPrefix(prefix string) (dbmodels.OperatorPrefixes, error) {
	res := dbmodels.OperatorPrefixes{}

	data := prefix[0:4]

	err := Dbcon.Where("prefix = ?", data).First(&res).Error
	if err != nil {
		fmt.Println("Failed to connect database OperatorPPOBPrefixes %v", err)
		return res, err
	}

	fmt.Println("[RESPONSE-LISTPRODUCT]-[GetOperatorCodebyPrefix]")
	fmt.Println("Data Prefix = %s", res)

	return res, nil
}
