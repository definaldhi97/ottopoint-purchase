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
		fmt.Println("[db]-[GetConfig]")
		fmt.Println(fmt.Sprintf("Failed to connect database transaction %v", err))

		return res, err
	}

	res.TransaksiPPOB = res.TransaksiPPOB / 100
	res.TransaksiPayQR = res.TransaksiPayQR / 100
	res.TransaksiMerchant = res.TransaksiMerchant / 100

	return res, nil
}

func GetData(trxID, InstitutionID string) (dbmodels.DeductTransaction, error) {
	res := dbmodels.DeductTransaction{}

	err := DbCon.Where("trx_id = ? and institution_id = ?", trxID, InstitutionID).First(&res).Error
	if err != nil {

		fmt.Println("[EEROR-DATABASE]")
		fmt.Println("[db]-[GetData]")
		fmt.Println(fmt.Sprintf("Failed to connect database transaction %v", err))

		return res, err
	}

	return res, nil
}

func GetVoucherUV(phone, couponID string) (dbmodels.UserMyVocuher, error) {
	res := dbmodels.UserMyVocuher{}

	err := DbCon.Where("phone = ? and coupon_id = ?", phone, couponID).First(&res).Error
	if err != nil {

		fmt.Println("[EEROR-DATABASE]")
		fmt.Println("[db]-[GetVoucherUV]")
		fmt.Println(fmt.Sprintf("Failed to connect database Voucher UV %v", err))

		return res, err
	}

	return res, nil
}

func ParamData(code string) (dbmodels.MParameters, error) {
	res := dbmodels.MParameters{}

	err := DbCon.Where("code = ?", code).First(&res).Error
	if err != nil {

		fmt.Println("[EEROR-DATABASE]")
		fmt.Println("[db]-[GetVoucherUV]")
		fmt.Println(fmt.Sprintf("Failed to connect database Voucher UV %v", err))

		return res, err
	}

	return res, nil
}
