package services

import (
	"errors"
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

	var useUV interface{}
	var reqUV interface{}

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[GetVoucherUV-Services]",
		zap.String("AccountNumber : ", param.AccountNumber), zap.String("InstitutionID : ", param.InstitutionID),
		zap.String("category : ", param.Category), zap.String("campaignId : ", campaignID),
		zap.String("VoucherCode : ", req.VoucherCode))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[GetVoucherUV]")
	defer span.Finish()

	logs.Info("Campaign : ", campaignID)
	logs.Info("CouponID : ", param.CouponID)
	logs.Info("ProductCode : ", param.ProductCode)
	logs.Info("CustID : ", param.CustID)

	// Use Voucher to Openloyalty
	use, err2 := opl.CouponVoucherCustomer(campaignID, param.CouponID, param.CouponCode, param.CustID, 1)
	if err2 != nil || use.Coupons[0].CouponID == "" {

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
