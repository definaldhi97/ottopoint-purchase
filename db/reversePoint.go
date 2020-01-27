package db

import (
	"fmt"
	"ottopoint-purchase/models/dbmodels"
)

func GetDataDeduct(trxID string) (dbmodels.ReversePoint, error) {
	res := dbmodels.ReversePoint{}

	err := Dbcon.Raw(`SELECT * FROM public.t_deduct_transaction where trx_id = '` + trxID + `' and status = '00'`).Scan(&res).Error
	if err != nil {

		fmt.Println("[EEROR-DATABASE]")
		fmt.Println("[db]")
		fmt.Println("[GetConfig]")
		fmt.Println(fmt.Sprintf("Failed to connect database transaction %v", err))

		return res, err
	}

	return res, nil
}

func UpdateDataDeduct(trxID string) (string, error) {
	err := Dbcon.Raw(`UPDATE public.t_deduct_transaction set status = '21' where trx_id = '` + trxID + `' and status = '00'`).Error
	if err != nil {
		fmt.Println("[EEROR-DATABASE]")
		fmt.Println("[db]")
		fmt.Println("[GetConfig]")
		fmt.Println(fmt.Sprintf("Failed to connect database transaction %v", err))

		return "Update status deduction failed", err
	}
	fmt.Println(fmt.Sprint("Update status deduction Success"))
	return "Update status deduction Success", nil
}
