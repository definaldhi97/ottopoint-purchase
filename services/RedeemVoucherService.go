package services

import (
	"errors"
	"ottopoint-purchase/hosts/opl/host"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"

	"github.com/astaxie/beego/logs"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type VoucherRedeemServices struct {
	General models.GeneralModel
}

func (t VoucherRedeemServices) VoucherRedeem(req models.RedeemReq, AccountNumber string) models.Response {
	var res models.Response

	resMeta := models.MetaData{
		Code:    200,
		Status:  true,
		Message: "Succesful",
	}

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[RedeemVoucher]", zap.String("CampaignID", req.CampaignID), zap.Int("Jumlah", req.Jumlah))
	logs.Info("[Request-RedeemVoucher] : ", req)
	logs.Info("[AccountNumber] : ", AccountNumber)

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[RedeemVoucher]")
	defer span.Finish()

	// make sure that minimal one request
	if req.Jumlah <= 0 {
		req.Jumlah = 1
	}

	// get data rewardValue(harga voucher)
	dataCoupon, errCoupon := host.VoucherDetail(req.CampaignID)
	if errCoupon != nil || dataCoupon.Name == "" {

		logs.Info("Internal Server Error : ", errCoupon)
		logs.Info("[VoucherRedeem-Services]")
		logs.Info("[Get VoucherDetail]")

		// sugarLogger.Info("Internal Server Error : ", errCoupon)
		sugarLogger.Info("[VoucherRedeem-Services]")
		sugarLogger.Info("[Get VoucherDetail]")

		res = utils.GetMessageResponse(res, 422, false, errors.New("Voucher tidak ditemukan"))
		return res
	}

	var voucher string
	coupon := []models.CouponsRedeem{}
	for i := req.Jumlah; i >= 1; i-- {
		data, err := host.RedeemVoucher(req.CampaignID, AccountNumber)
		if err != nil {
			logs.Info("Internal Server Error : ", err)
			logs.Info("[VoucherRedeem-Services]")
			logs.Info("[Failed Redeem Voucher]")

			// sugarLogger.Info("Internal Server Error : ", err)
			sugarLogger.Info("[VoucherRedeem-Services]")
			sugarLogger.Info("[Failed Redeem Voucher]")

			res = utils.GetMessageResponse(res, 500, false, errors.New("Gagal Redeem Voucher"))
			return res
		}

		for _, val := range data.Coupons {
			a := models.CouponsRedeem{
				Code: val.Code,
			}

			if dataCoupon.CampaignID == req.CampaignID {
				voucher = dataCoupon.Name
			}

			// for _, value := range dataCoupon {
			// 	if value.CampaignID == req.CampaignID {
			// 		voucher = value.Name
			// 	}
			// }
			a.Voucher = voucher
			coupon = append(coupon, a)
		}
	}

	logs.Info("Voucher :", coupon)
	// check if no data founded
	if len(coupon) == 0 {
		logs.Info("[VoucherRedeem-Services]")
		logs.Info("[Voucher Kosong]")

		sugarLogger.Info("[VoucherRedeem-Services]")
		sugarLogger.Info("[Voucher Kosong]")

		res = utils.GetMessageResponse(res, 422, false, errors.New("Anda mencapai batas maksimal pembelian voucher"))
		return res
	}

	res = models.Response{
		Meta: resMeta,
		Data: models.RedeemResp{
			CodeVoucher: coupon,
		},
	}
	return res
}
