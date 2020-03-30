package services

import (
	"errors"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
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
	sugarLogger.Info("[UseVoucher-Services]",
		zap.String("AccountNumber : ", param.AccountNumber), zap.String("InstitutionID : ", param.InstitutionID),
		zap.String("category : ", req.Category), zap.String("campaignId : ", req.CampaignID),
		zap.String("cust_id : ", req.CustID), zap.String("cust_id2 : ", req.CustID2),
		zap.String("product_code : ", param.ProductCode))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[RedeemVoucher]")
	defer span.Finish()

	// get CustID
	dataUser, errUser := db.CheckUser(param.AccountNumber)
	if errUser != nil {
		res = utils.GetMessageResponse(res, 422, false, errors.New("User belum Eligible"))
		return res
	}

	// Use Voucher to Openloyalty
	_, err2 := opl.CouponVoucherCustomer(req.CampaignID, param.CouponID, param.ProductCode, dataUser.CustID, 1)
	if err2 != nil {
		res = utils.GetMessageResponse(res, 400, false, errors.New("Gagal Redeem Voucher, Harap coba lagi"))
		return res
	}

	category := strings.ToLower(req.Category)

	resRedeem := models.UseRedeemResponse{}

	reqRedeem := models.UseRedeemRequest{
		AccountNumber: param.AccountNumber,
		CustID:        req.CustID,
		CustID2:       req.CustID2,
		ProductCode:   param.ProductCode,
	}

	switch category {
	case constants.CategoryPulsa, constants.CategoryPaketData:
		resRedeem = redeem.RedeemPulsa(reqRedeem, req, param)
	case constants.CategoryToken:
		resRedeem = redeem.RedeemPLN(reqRedeem, req, param)
	case constants.CategoryMobileLegend, constants.CategoryFreeFire:
		resRedeem = redeem.RedeemGame(reqRedeem, req, param)
	}

	if resRedeem.Msg == "Prefix Failed" {
		logs.Info("[Prefix Failed]")
		logs.Info("[Services-Voucher-UserVoucher]")

		res = utils.GetMessageResponse(res, 500, false, errors.New("Transaksi Gagal"))
		return res
	}

	if resRedeem.Msg == "Inquiry Failed" {
		logs.Info("[Inquiry Failed]")
		logs.Info("[Services-Voucher-UserVoucher]")

		logs.Info("[Reversal Voucher")
		_, erv := opl.CouponVoucherCustomer(req.CampaignID, param.CouponID, param.ProductCode, dataUser.CustID, 0)
		if erv != nil {
			res = utils.GetMessageResponse(res, 500, false, errors.New("Gagal Reversal Voucher"))
			return res
		}

		res = utils.GetMessageResponse(res, 500, false, errors.New("Transaksi Gagal"))
		return res
	}

	if resRedeem.Msg == "Payment Failed" {
		logs.Info("[Payment Failed]")
		logs.Info("[Services-Voucher-UserVoucher]")

		logs.Info("[Reversal Voucher")
		_, erv := opl.CouponVoucherCustomer(req.CampaignID, param.CouponID, param.ProductCode, dataUser.CustID, 0)
		if erv != nil {
			res = utils.GetMessageResponse(res, 500, false, errors.New("Gagal Reversal Voucher"))
			return res
		}

		res = utils.GetMessageResponse(res, 500, false, errors.New("Transaksi Gagal"))
		return res
	}

	if resRedeem.Msg == "Request in progress" {
		logs.Info("[Prefix Failed]")
		logs.Info("[Services-Voucher-UserVoucher]")

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
