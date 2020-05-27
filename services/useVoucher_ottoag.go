package services

import (
	"errors"
	"ottopoint-purchase/constants"
	opl "ottopoint-purchase/hosts/opl/host"
	"ottopoint-purchase/models"
	redeem "ottopoint-purchase/services/voucher"
	"ottopoint-purchase/utils"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type UseVoucherServices struct {
	General models.GeneralModel
}

func (t UseVoucherServices) UseVoucherOttoAG(req models.UseVoucherReq, param models.Params) models.Response {
	var res models.Response

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[UseVoucherOttoAG-Services]",
		zap.String("AccountNumber : ", param.AccountNumber), zap.String("InstitutionID : ", param.InstitutionID),
		zap.String("category : ", param.Category), zap.String("campaignId : ", req.CampaignID),
		zap.String("cust_id : ", req.CustID), zap.String("cust_id2 : ", req.CustID2),
		zap.String("product_code : ", param.ProductCode))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[UseVoucherOttoAG]")
	defer span.Finish()

	// Use Voucher to Openloyalty
	_, err2 := opl.CouponVoucherCustomer(req.CampaignID, param.CouponID, param.ProductCode, param.AccountId, 1)
	if err2 != nil {
		res = utils.GetMessageResponse(res, 400, false, errors.New("Gagal Redeem Voucher, Harap coba lagi"))
		return res
	}

	category := strings.ToLower(param.Category)

	resRedeem := models.UseRedeemResponse{}

	reqRedeem := models.UseRedeemRequest{
		AccountNumber: param.AccountNumber,
		CustID:        req.CustID,
		CustID2:       req.CustID2,
		ProductCode:   param.ProductCode,
	}

	switch category {
	case constants.CategoryPulsa:
		resRedeem = redeem.RedeemPulsa(reqRedeem, req, param)
	case constants.CategoryPLN:
		resRedeem = redeem.RedeemPLN(reqRedeem, req, param)
	case constants.CategoryMobileLegend, constants.CategoryFreeFire:
		resRedeem = redeem.RedeemGame(reqRedeem, req, param)
	}

	if resRedeem.Msg == "Prefix Failed" {
		logs.Info("[Prefix Failed]")
		logs.Info("[UseVoucherOttoAG]")

		res = utils.GetMessageResponse(res, 500, false, errors.New("Transaksi Gagal"))
		return res
	}

	if resRedeem.Msg == "Inquiry Failed" {
		logs.Info("[Inquiry Failed]")
		logs.Info("[UseVoucherOttoAG]")

		logs.Info("[Reversal Voucher")
		_, erv := opl.CouponVoucherCustomer(req.CampaignID, param.CouponID, param.ProductCode, param.AccountId, 0)
		if erv != nil {
			res = utils.GetMessageResponse(res, 500, false, errors.New("Gagal Reversal Voucher"))
			return res
		}

		res = utils.GetMessageResponse(res, 500, false, errors.New("Transaksi Gagal"))
		return res
	}

	if resRedeem.Msg == "Payment Failed" {
		logs.Info("[Payment Failed]")
		logs.Info("[UseVoucherOttoAG]")

		logs.Info("[Reversal Voucher")
		_, erv := opl.CouponVoucherCustomer(req.CampaignID, param.CouponID, param.ProductCode, param.AccountId, 0)
		if erv != nil {
			res = utils.GetMessageResponse(res, 500, false, errors.New("Gagal Reversal Voucher"))
			return res
		}

		res = utils.GetMessageResponse(res, 500, false, errors.New("Transaksi Gagal"))
		return res
	}

	if resRedeem.Msg == "Request in progress" {
		logs.Info("[Prefix Failed]")
		logs.Info("[UseVoucherOttoAG]")

		res = utils.GetMessageResponse(res, 500, false, errors.New("Transaksi Pending"))
		return res
	}

	if resRedeem.Msg == "SUCCESS" {
		if resRedeem.Category == "PLN" {
			res = models.Response{
				Data: models.ResponseUseVoucherPLN{
					Voucher:     param.NamaVoucher,
					CustID:      resRedeem.CustID,
					CustID2:     resRedeem.CustID2,
					ProductCode: resRedeem.ProductCode,
					Amount:      resRedeem.Amount,
					Token:       resRedeem.Data.Tokenno,
				},
				Meta: utils.ResponseMetaOK(),
			}
			return res
		}

		res = models.Response{
			Data: models.ResponseUseVoucher{
				Voucher:     param.NamaVoucher,
				CustID:      resRedeem.CustID,
				CustID2:     resRedeem.CustID2,
				ProductCode: resRedeem.ProductCode,
				Amount:      resRedeem.Amount,
			},
			Meta: utils.ResponseMetaOK(),
		}
		return res
	}

	return res
}
