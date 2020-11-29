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

func GetUltraVoucher(voucherCode, accountId string) (dbmodels.UserMyVocuher, error) {
	res := dbmodels.UserMyVocuher{}

	err := DbCon.Where("voucher_code = ? and account_id = ?", voucherCode, accountId).First(&res).Error
	if err != nil {

		fmt.Println("[EEROR-DATABASE]")
		fmt.Println("[db]-[GetUltraVoucher]")
		fmt.Println(fmt.Sprintf("Failed to connect database GetUltraVoucher %v", err))

		return res, err
	}

	return res, nil
}

func CheckCouponUV(phone, campaign, couponId string) (dbmodels.UserMyVocuher, error) {
	res := dbmodels.UserMyVocuher{}

	err := DbCon.Raw(`SELECT * FROM public.user_myvoucher where phone = ? and campaign_id = ? and coupon_id = ?`, phone, campaign, couponId).Scan(&res).Error
	if err != nil {

		fmt.Println("[EEROR-DATABASE]")
		fmt.Println("[db]-[GetCouponUV]")
		fmt.Println(fmt.Sprintf("Failed to connect database GetCouponUV %v", err))

		return res, err
	}

	return res, nil
}

func ParamData(group, code string) (dbmodels.MParameters, error) {
	res := dbmodels.MParameters{}

	// err := DbCon.Where("code = ?", code).First(&res).Error
	err := DbCon.Raw(`select * from public.m_parameters mp where mp."group" = ? and  mp.code = ?`, group, code).Scan(&res).Error
	if err != nil {

		fmt.Println("[EEROR-DATABASE]")
		fmt.Println("[db]-[GetVoucherUV]")
		fmt.Println(fmt.Sprintf("Failed to connect database Voucher UV %v", err))

		return res, err
	}

	return res, nil
}

func GetDataInstitution(Institution string) (dbmodels.MInstution, error) {
	result := dbmodels.MInstution{}

	err := DbCon.Raw(`select * from public.m_institution where partner_id = ?`, Institution).Scan(&result).Error
	if err != nil {
		fmt.Println("[EEROR-DATABASE]")
		fmt.Println("[db]")
		fmt.Println("[Get Instution / Issuer in database]")
		fmt.Println(fmt.Sprintf("Failed to connect database %v", err))

		return result, err
	}
	return result, nil

}

func GetSpendingSepulsa(transactionID, orderID string) (dbmodels.TSpending, error) {
	res := dbmodels.TSpending{}

	err := DbCon.Where("transaction_id = ? and rrn = ?", orderID, transactionID).First(&res).Error
	if err != nil {

		fmt.Println("[EEROR-DATABASE]")
		fmt.Println("[db]-[GetSpendingSepulsa]")
		fmt.Println(fmt.Sprintf("Failed to connect database GetSpendingSepulsa %v", err))

		return res, err
	}

	return res, nil
}
