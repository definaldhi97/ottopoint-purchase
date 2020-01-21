package services

import (
	"errors"
	"ottopoint-purchase/hosts/opl/host"
	hostopl "ottopoint-purchase/hosts/opl/host"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"

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

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[RedeemVoucher]")
	defer span.Finish()

	// make sure that minimal one request
	var total int
	if req.Jumlah <= 0 {
		total = 1
	}

	dataVoucher, errVoucher := host.VoucherDetail2(req.CampaignID)
	if errVoucher != nil {
		res = utils.GetMessageResponse(res, 500, false, errors.New("Internal Server Error"))
		return res
	}

	var voucher string
	coupon := []models.CouponsRedeem{}
	for i := total; i >= 1; i-- {
		data, err := host.RedeemVoucher(req.CampaignID, AccountNumber)
		if err != nil {
			res = utils.GetMessageResponse(res, 500, false, errors.New("Internal Server Error"))
			return res
		}

		for _, val := range data.Coupons {
			a := models.CouponsRedeem{
				Code:    val.Code,
				Voucher: voucher,
			}
			for _, value := range dataVoucher {
				if value.Coupons[0].Coupon == val.Code {
					voucher = value.Name
				}
			}

			a.Voucher = voucher
			coupon = append(coupon, a)
		}
	}

	// check if no data founded
	if len(coupon) == 0 {
		res = utils.GetMessageResponse(res, 422, false, errors.New("Anda mencapai batas maksimal pembelian voucher"))
		return res
	}

	// get data rewardValue(harga voucher)
	dataCoupon, errCoupon := hostopl.VoucherDetail(req.CampaignID)
	if errCoupon != nil {
		res = utils.GetMessageResponse(res, 422, false, errors.New("Internal Server Error"))
		return res
	}

	if dataCoupon.Name == "" {
		res = utils.GetMessageResponse(res, 422, false, errors.New("Voucher tidak ditemukan"))
		return res
	}

	// go redis.SaveRedis(fmt.Sprintf("Harga-Voucher-%s-%s :", req.CampaignID, AccountNumber), dataCoupon.CostInPoints)

	res = models.Response{
		Meta: resMeta,
		Data: models.RedeemResp{
			CodeVoucher: coupon,
		},
	}
	return res
}
