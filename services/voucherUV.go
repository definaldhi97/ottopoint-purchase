package services

import (
	"errors"
	"fmt"
	opl "ottopoint-purchase/hosts/opl/host"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"

	"github.com/astaxie/beego/logs"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type UseVoucherUVServices struct {
	General models.GeneralModel
}

func (t UseVoucherUVServices) UseVoucherUV(req models.UseVoucherUVReq, param models.Params, campaignID string) models.Response {
	var res models.Response

	logs.Info("=== UseVoucherUV ===")
	fmt.Sprintf("=== UseVoucherUV ===")

	var useUV interface{}
	var reqUV interface{}

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[UseVoucherUV]",
		zap.String("AccountNumber : ", param.AccountNumber), zap.String("InstitutionID : ", param.InstitutionID),
		zap.String("category : ", param.Category), zap.String("campaignId : ", campaignID),
		zap.String("VoucherCode : ", req.VoucherCode))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[UseVoucherUV]")
	defer span.Finish()

	logs.Info("Campaign : ", campaignID)
	logs.Info("CouponID : ", param.CouponID)
	logs.Info("ProductCode : ", param.CouponCode)
	logs.Info("CustID : ", param.CustID)

	// Use Voucher to Openloyalty
	use, err2 := opl.CouponVoucherCustomer(campaignID, param.CouponID, param.CouponCode, param.CustID, 1)

	var useErr string
	for _, value := range use.Coupons {
		useErr = value.CouponID
	}

	if err2 != nil || useErr == "" {

		logs.Info(fmt.Sprintf("[Error : %v]", err2))
		logs.Info(fmt.Sprintf("[Response : %v]", use))
		logs.Info("[Error from OPL]-[CouponVoucherCustomer]")

		go SaveTransactionUV(param, useUV, reqUV, req, "Payment", "01", "")

		res = utils.GetMessageResponse(res, 400, false, errors.New("Gagal Use Voucher, Harap coba lagi"))
		return res
	}

	go SaveTransactionUV(param, useUV, reqUV, req, "Payment", "00", "00")

	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: models.UseVoucherUVResp{
			Voucher: param.NamaVoucher,
		},
	}
	return res

}
