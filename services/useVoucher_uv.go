package services

import (
	"errors"
	"ottopoint-purchase/db"
	opl "ottopoint-purchase/hosts/opl/host"
	uv "ottopoint-purchase/hosts/ultra_voucher/host"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"

	"github.com/astaxie/beego/logs"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

func (t UseVoucherServices) UseVoucherUV(req models.UseVoucherReq, param models.Params) models.Response {
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

	// save to redis usedAt
	// req.Date = (time.Now().Local().Add(time.Hour * time.Duration(7))).Format("2006-01-02T15:04:05-0700")
	// keyTimeVoucher := fmt.Sprintf("usedVoucherAt-%s-%s", param.CouponID, param.AccountNumber)
	// go redis.SaveRedis(keyTimeVoucher, req.Date)

	get, errGet := db.GetVoucherUV(param.AccountNumber, param.CouponID)
	if errGet != nil {
		logs.Info("Internal Server Error : ", errGet)
		logs.Info("[UseVoucherUV-Servcies]-[GetVoucherUV]")
		logs.Info("[Failed get data from DB]")

		// sugarLogger.Info("Internal Server Error : ", errGet)
		sugarLogger.Info("[UseVoucherUV-Servcies]-[GetVoucherUV]")
		sugarLogger.Info("[Failed get data from DB]")

		res = utils.GetMessageResponse(res, 400, false, errors.New("Internal Server Error"))
		return res
	}

	// get to UV
	useUV, errUV := uv.UseVoucherUV(param.AccountNumber, get.VoucherCode)
	if errUV != nil {
		logs.Info("Internal Server Error : ", errUV)
		logs.Info("[UseVoucherUV-Servcies]-[UseVoucherUV]")
		logs.Info("[Failed Use Voucher UV]-[Gagal Use Voucher UV]")

		// sugarLogger.Info("Internal Server Error : ", errUV)
		sugarLogger.Info("[UseVoucherUV-Servcies]-[UseVoucherUV]")
		sugarLogger.Info("[Failed Use Voucher UV]-[Gagal Use Voucher UV]")

		res = utils.GetMessageResponse(res, 400, false, errors.New("Gagal Redeem Voucher, Harap coba lagi"))
		return res
	}

	if useUV.ResponseCode == "" {
		// Use Voucher to Openloyalty
		_, err2 := opl.CouponVoucherCustomer(req.CampaignID, param.CouponID, param.ProductCode, dataUser.CustID, 0)
		if err2 != nil {
			res = utils.GetMessageResponse(res, 400, false, errors.New("Gagal Redeem Voucher, Harap coba lagi"))
			return res
		}
		res = utils.GetMessageResponse(res, 129, false, errors.New("Voucher gagal digunakan, silahkan coba beberapa saat lagi"))
		return res
	}

	if useUV.ResponseCode == "00" {
		res = models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: models.UseVoucherUVResp{
				Voucher: param.NamaVoucher,
				Link:    useUV.Data.Link,
			},
		}
		return res
	}

	return res
}
