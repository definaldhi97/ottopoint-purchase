package services

import (
	"fmt"
	"ottopoint-purchase/constants"
	db "ottopoint-purchase/db"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"

	"github.com/astaxie/beego/logs"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type ViewVoucherService struct {
	General models.GeneralModel
}

func (t ViewVoucherService) ViewVoucher(accountNumber, couponId string) models.Response {
	fmt.Println("[ >>>>>>>>>>>>>>>>>>>>> View Voucher Service <<<<<<<<<<<<<<<<<< ]")

	resp := models.Response{Meta: utils.ResponseMetaOK()}

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[ViewVoucher-Services]",
		zap.String("Coupon Id : ", couponId))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[ViewVoucher]")
	defer span.Finish()

	// get voucher
	getVouc, errGetVouc := db.GetVoucher(accountNumber, couponId)
	if errGetVouc != nil || getVouc.AccountNumber == "" {
		logs.Info("Internal server error")
		logs.Info("[Failed get data voucher from DB]")

		// sugarLogger.Info("Internal Server Error : ", errGet)
		sugarLogger.Info("[GetVoucher-Servcies]-[GetVoucher]")
		sugarLogger.Info("[Failed get data from DB]")
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
