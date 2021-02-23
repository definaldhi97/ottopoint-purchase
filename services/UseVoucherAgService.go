package services

import (
	"errors"
	db "ottopoint-purchase/db"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

func (t UseVoucherServices) UseVoucherAggregator(req models.UseVoucherReq, param models.Params) models.Response {
	var res models.Response

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[UseVoucherAggregator-Services]",
		zap.String("AccountNumber : ", param.AccountNumber), zap.String("InstitutionID : ", param.InstitutionID),
		zap.String("category : ", param.Category), zap.String("campaignId : ", req.CampaignID),
		zap.String("cust_id : ", req.CustID), zap.String("cust_id2 : ", req.CustID2),
		zap.String("product_code : ", param.ProductCode))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[GetVoucherUV]")
	defer span.Finish()

	get, err := db.GetVoucherAg(param.AccountId, param.CouponID)
	if err != nil {
		logrus.Info("Internal Server Error : ", err)
		logrus.Info("[UseVoucherAggregator]-[GetVoucherUV]")
		logrus.Info("[Failed get data from DB]")

		sugarLogger.Info("[UseVoucherAggregator]-[GetVoucherUV]")
		sugarLogger.Info("[Failed get data from DB]")

		res = utils.GetMessageResponse(res, 422, false, errors.New("Voucher Tidak Ditemukan"))
		return res
	}

	spend, err := db.GetVoucherSpending(get.AccountId, get.CouponID)
	if err != nil {
		logrus.Info("Internal Server Error : ", err)
		logrus.Info("[UseVoucherAggregator]-[GetVoucherUV]")
		logrus.Info("[Failed get data from DB]")

		sugarLogger.Info("[UseVoucherAggregator]-[GetVoucherUV]")
		sugarLogger.Info("[Failed get data from DB]")

		res = utils.GetMessageResponse(res, 422, false, errors.New("Terjadi Kesalahan"))
		return res
	}

	// Update Status Voucher
	// timeUse := jodaTime.Format("dd-MM-YYYY HH:mm:ss", time.Now())
	timeUse := time.Now()
	go db.UpdateVoucher(timeUse, spend.CouponId)

	codeVoucher := t.decryptVoucherCode(spend.VoucherCode, spend.CouponId)

	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: models.GetVoucherAgResp{
			Voucher:     param.NamaVoucher,
			VoucherCode: codeVoucher,
			Link:        spend.VoucherLink,
		},
	}
	return res
}

func (t UseVoucherServices) decryptVoucherCode(voucherCode, couponID string) string {

	var codeVoucher string
	if voucherCode == "" {
		return voucherCode
	}

	a := []rune(couponID)
	key32 := string(a[0:32])
	secretKey := []byte(key32)
	codeByte := []byte(voucherCode)
	chiperText, _ := utils.DecryptAES(codeByte, secretKey)
	codeVoucher = string(chiperText)

	return codeVoucher

}
