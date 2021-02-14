package vouchers

import (
	"fmt"
	"ottopoint-purchase/constants"
	db "ottopoint-purchase/db"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"

	"github.com/sirupsen/logrus"
)

func ViewVoucherServices(accountNumber, couponId string) models.Response {

	nameservice := "[PackageVouchers]-[ViewVoucherServices]"
	logReq := fmt.Sprintf("[AccountNumber : %v || CouponId : %v]", accountNumber, couponId)

	logrus.Info(nameservice)

	resp := models.Response{Meta: utils.ResponseMetaOK()}

	// get voucher
	getVouc, errGetVouc := db.GetVoucher(accountNumber, couponId)
	if errGetVouc != nil || getVouc.AccountNumber == "" {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[GetVoucher]-[Error : %v]", errGetVouc))
		logrus.Println(logReq)

		resp = utils.GetMessageFailedErrorNew(resp, constants.RC_VOUCHER_NOTFOUND, constants.RD_VOUCHER_NOTFOUND)
		return resp
	}

	// decrypt voucher code
	a := []rune(getVouc.CouponId)
	key32 := string(a[0:32])
	key := []byte(key32)
	chiperText := []byte(getVouc.VoucherCode)
	plainText, errDec := utils.DecryptAES(chiperText, key)
	plainTextVoucCod := fmt.Sprintf("%s", plainText)
	fmt.Println("voucher code : ")
	fmt.Println(plainText)
	if errDec != nil {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[Error Decrytp]-[Error : %v]", errDec))
		logrus.Println(logReq)

		resp = utils.GetMessageFailedErrorNew(resp, constants.RC_FAILED_DECRYPT_VOUCHER, constants.RD_FAILED_DECRYPT_VOUCHER)
		return resp
	}

	// get path product brand
	pathpathImg, errPath := db.GetPathImageProduct(getVouc.ProductType)
	if errPath != nil {
		resp = utils.GetMessageFailedErrorNew(resp, 500, "Internal Server Error")
		return resp
	}
	patahProductBrand := utils.UrlImage + pathpathImg.Path

	xpd := getVouc.ExpDate.Format("2006-01-02 15:04:05")

	//resp view voucher
	dataVouch := models.ViewVocuherVidio{
		VoucherName: getVouc.Voucher,
		ExpiredDate: xpd,
		VoucherCode: plainTextVoucCod,
		ImageUrl:    patahProductBrand,
	}

	resp.Data = dataVouch
	return resp
}
