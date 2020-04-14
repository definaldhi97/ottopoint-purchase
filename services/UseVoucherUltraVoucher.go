package services

import (
	"errors"
	"ottopoint-purchase/db"
	"ottopoint-purchase/hosts/opl/host"
	uv "ottopoint-purchase/hosts/ultra_voucher/host"
	"ottopoint-purchase/models"
	"ottopoint-purchase/models/dbmodels"
	"ottopoint-purchase/utils"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type UseVoucherUltraVoucher struct {
	General models.GeneralModel
}

func (t UseVoucherUltraVoucher) UltraVoucherServices(req models.VoucherComultaiveReq, param models.Params) models.Response {
	var res models.Response

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[UltraVoucherServices]",
		zap.String("NameVoucher : ", param.NamaVoucher), zap.Int("Jumlah : ", req.Jumlah),
		zap.String("CampaignID : ", req.CampaignID), zap.String("CampaignID : ", req.CampaignID),
		zap.String("CustID2 : ", req.CustID2), zap.String("ProductCode : ", param.ProductCode),
		// zap.Int("Point : ", req.Point),
		zap.String("AccountNumber : ", param.AccountNumber), zap.String("InstitutionID : ", param.InstitutionID))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[UltraVoucherServices]")
	defer span.Finish()

	// var err bool
	// success := 0
	// failed := 0
	couponOPL := []models.CouponsRedeem{}

	// for i := req.Jumlah; i >= 1; i-- {

	dataorder := DataParameterOrder()

	param.Reffnum = utils.GenTransactionId()
	param.ExpDate = dataorder.Expired

	total := strconv.Itoa(req.Jumlah)

	// redeem to opl (potong point)
	redeem, errredeem := host.RedeemVoucherCumulative(req.CampaignID, param.CustID, total)
	if errredeem != nil || redeem.Error != "" {
		logs.Info("Internal Server Error : ", errredeem)
		logs.Info("[UltraVoucherServices]-[RedeemVoucher]")
		logs.Info("[Failed Redeem Voucher]-[Gagal Redeem Voucher]")

		// sugarLogger.Info("Internal Server Error : ", errredeem)
		sugarLogger.Info("[UltraVoucherServices]-[RedeemVoucher]")
		sugarLogger.Info("[Failed Redeem Voucher]-[Gagal Redeem Voucher]")

		res = utils.GetMessageResponse(res, 500, false, errors.New("Gagal Redeem Voucher"))
		return res
	}

	// if redeem.Error != "" {
	// 	sugarLogger.Info("[UltraVoucherServices]-[RedeemVoucher]")
	// 	sugarLogger.Info(redeem.Error)

	// 	logs.Info("[UltraVoucherServices]-[RedeemVoucher]-[Error : %v]", redeem.Error)
	// }

	// order to u
	order, errOrder := uv.OrderVoucher(param, req.Jumlah, dataorder.Email, dataorder.Phone, dataorder.Nama)
	if errOrder != nil || order.ResponseCode != "00" {
		logs.Info("Internal Server Error : ", errOrder)
		logs.Info("ResponseCode : ", order.ResponseCode)
		logs.Info("[UltraVoucherServices]-[OrderVoucher]")
		logs.Info("[Failed Order Voucher]-[Gagal Order Voucher]")

		// sugarLogger.Info("Internal Server Error : ", errOrder)
		sugarLogger.Info("[UltraVoucherServices]-[OrderVoucher]")
		sugarLogger.Info("[Failed Order Voucher]-[Gagal Order Voucher]")

		Text := "OP009 - " + "Reversal point cause transaction " + param.NamaVoucher + " is failed"
		point := strconv.Itoa(param.Point)
		_, errReversal := host.TransferPoint(param.CustID, point, Text)
		if errReversal != nil {
			logs.Info("Internal Server Error : ", errReversal)
			logs.Info("[UltraVoucherServices]-[TransferPoint]")
			logs.Info("[Failed Order Voucher]-[Gagal Reversal Point]")
		}

		res = models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: models.UltraVoucherResp{
				Success: 0,
				Failed:  req.Jumlah,
				Total:   req.Jumlah,
				Voucher: param.NamaVoucher,
			},
		}

		return res
	}

	// opl
	var coupon, code string
	for _, val := range redeem.Coupons {
		coupon = val.Id
		a := models.CouponsRedeem{
			Code: val.Code,
			ID:   val.Id,
		}
		// coupon = val.

		couponOPL = append(couponOPL, a)
	}

	// uv
	for _, value := range order.Data.VouchersCode {
		code = value.Code
	}

	id := utils.GenerateTokenUUID()
	go SaveDB(id, param.InstitutionID, coupon, code, param.AccountNumber)
	// success++
	// }

	// if err == true {
	// 	res = models.Response{
	// 		Meta: models.MetaData{
	// 			Status:  true,
	// 			Message: "Point Tidak Cukup",
	// 			Code:    201,
	// 		},
	// 		Data: models.UltraVoucherResp{
	// 			Success: success,
	// 			Failed:  failed,
	// 			Total:   req.Jumlah,
	// 			Voucher: param.NamaVoucher,
	// 		},
	// 	}
	// 	return res
	// }

	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: models.UltraVoucherResp{
			Success: req.Jumlah,
			Failed:  0,
			Total:   req.Jumlah,
			Voucher: param.NamaVoucher,
		},
	}

	return res
}

func SaveDB(id, institution, coupon, vouchercode, phone string) {
	save := dbmodels.UserMyVocuher{
		ID:            id,
		InstitutionID: institution,
		CouponID:      coupon,
		VoucherCode:   vouchercode,
		Phone:         phone,
	}

	err := db.DbCon.Create(&save).Error
	if err != nil {
		logs.Info("[Failed Save to DB ]", err)
		logs.Info("[Package-Services]-[UltraVoucherServices]")
		// return err
	}
}

func DataParameterOrder() models.ParamUV {
	res := models.ParamUV{}

	nama := "" // nama
	email := "UV_EMAIL_ORDER"
	phone := "UV_PHONE_ORDER"
	expired := "UV_EXPIRED_VOUCHER"

	datanama, errnama := db.ParamData(nama)
	if errnama != nil {
		logs.Info("[Error get data from Db m_paramaters]")
		logs.Info("Error : ", errnama)
		logs.Info("Code :", nama)
	}

	dataemail, erremail := db.ParamData(email)
	if erremail != nil {
		logs.Info("[Error get data from Db m_paramaters]")
		logs.Info("Error : ", erremail)
		logs.Info("Code :", email)
	}

	dataphone, errphone := db.ParamData(phone)
	if errphone != nil {
		logs.Info("[Error get data from Db m_paramaters]")
		logs.Info("Error : ", errphone)
		logs.Info("Code :", phone)
	}

	dataexpired, errexpired := db.ParamData(expired)
	if errexpired != nil {
		logs.Info("[Error get data from Db m_paramaters]")
		logs.Info("Error : ", errexpired)
		logs.Info("Code :", expired)
	}

	res = models.ParamUV{
		Nama:    datanama.Value,
		Email:   dataemail.Value,
		Phone:   dataphone.Value,
		Expired: dataexpired.Value,
	}

	return res
}
