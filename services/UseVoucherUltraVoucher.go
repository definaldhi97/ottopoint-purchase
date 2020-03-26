package services

import (
	"ottopoint-purchase/db"
	"ottopoint-purchase/hosts/opl/host"
	uv "ottopoint-purchase/hosts/ultra_voucher/host"
	"ottopoint-purchase/models"
	"ottopoint-purchase/models/dbmodels"
	"ottopoint-purchase/utils"

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
	sugarLogger.Info("[VoucherComulative-Services]",
		zap.String("NameVoucher : ", param.NamaVoucher), zap.Int("Jumlah : ", req.Jumlah),
		zap.String("CampaignID : ", req.CampaignID), zap.String("CampaignID : ", req.CampaignID),
		zap.String("CustID2 : ", req.CustID2), zap.String("ProductCode : ", param.ProductCode),
		zap.Int("Point : ", req.Point),
		zap.String("AccountNumber : ", param.AccountNumber), zap.String("InstitutionID : ", param.InstitutionID))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[UltraVoucherServices]")
	defer span.Finish()

	success := 0
	failed := 0
	couponOPL := []models.CouponsRedeem{}
	for i := req.Jumlah; i >= 1; i-- {

		dataorder := DataParameterOrder()

		param.Reffnum = utils.GenTransactionId()
		param.ExpDate = dataorder.Expired

		// order to u
		order, errOrder := uv.OrderVoucher(param, req.Jumlah, dataorder.Email, dataorder.Phone, dataorder.Nama)
		if errOrder != nil {
			logs.Info("Internal Server Error : ", errOrder)
			logs.Info("[UltraVoucherServices]-[OrderVoucher]")
			logs.Info("[Failed Order Voucher]-[Gagal Order Voucher]")

			sugarLogger.Info("[UltraVoucherServices]-[OrderVoucher]")
			sugarLogger.Info("[Failed Order Voucher]-[Gagal Order Voucher]")
			// sugarLogger.Info("Internal Server Error : ", errOrder)

			// res = utils.GetMessageResponse(res, 500, false, errors.New("Gagal Order Voucher"))
			// return res
			failed++
			continue
		}

		redeem, errredeem := host.RedeemVoucher(req.CampaignID, param.AccountNumber)
		if errredeem != nil {
			logs.Info("Internal Server Error : ", errredeem)
			logs.Info("[UltraVoucherServices]-[RedeemVoucher]")
			logs.Info("[Failed Redeem Voucher]-[Gagal Redeem Voucher]")

			sugarLogger.Info("[UltraVoucherServices]-[RedeemVoucher]")
			sugarLogger.Info("[Failed Redeem Voucher]-[Gagal Redeem Voucher]")
			// sugarLogger.Info("Internal Server Error : ", errredeem)

			// res = utils.GetMessageResponse(res, 500, false, errors.New("Gagal Redeem Voucher"))
			// return res
			failed++
			continue
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
	}

	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: models.UltraVoucherResp{
			Success: success,
			Failed:  failed,
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
		Nama:    datanama.Desc,
		Email:   dataemail.Desc,
		Phone:   dataphone.Desc,
		Expired: dataexpired.Desc,
	}

	return res
}
