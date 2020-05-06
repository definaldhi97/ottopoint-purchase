package services

import (
	"errors"
	"fmt"
	"ottopoint-purchase/db"
	"ottopoint-purchase/hosts/opl/host"
	opl "ottopoint-purchase/hosts/opl/host"
	uv "ottopoint-purchase/hosts/ultra_voucher/host"
	uvmodels "ottopoint-purchase/hosts/ultra_voucher/models"
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
	// couponOPL := []models.CouponsRedeem{}

	// for i := req.Jumlah; i >= 1; i-- {

	dataorder := DataParameterOrder()

	param.Reffnum = utils.GenTransactionId()
	param.ExpDate = dataorder.Expired

	total := strconv.Itoa(req.Jumlah)

	param.Amount = int64(param.Point)

	// redeem to opl (potong point)
	redeem, errredeem := host.RedeemVoucherCumulative(req.CampaignID, param.AccountNumber, total)

	if redeem.Message == "Invalid JWT Token" {
		logs.Info("Error : ", errredeem)
		logs.Info("[UltraVoucherServices]-[RedeemVoucher]")
		logs.Info("[Internal Server Error]-[Gagal Redeem Voucher]")

		// sugarLogger.Info("Internal Server Error : ", errredeem)
		sugarLogger.Info("[UltraVoucherServices]-[RedeemVoucher]")
		sugarLogger.Info("[Internal Server Error]-[Gagal Redeem Voucher]")

		res = utils.GetMessageResponse(res, 422, false, errors.New("Internal Server Error"))
		// res.Data = "Transaksi Gagal"

		return res
	}

	if redeem.Error == "Not enough points" {
		logs.Info("Error : ", errredeem)
		logs.Info("[UltraVoucherServices]-[RedeemVoucher]")
		logs.Info("[Not enough points]-[Gagal Redeem Voucher]")

		// sugarLogger.Info("Internal Server Error : ", errredeem)
		sugarLogger.Info("[UltraVoucherServices]-[RedeemVoucher]")
		sugarLogger.Info("[Not enough points]-[Gagal Redeem Voucher]")

		res = utils.GetMessageResponse(res, 500, false, errors.New("Point Tidak Cukup"))
		// res.Data = "Transaksi Gagal"

		return res
	}

	if redeem.Error == "Limit exceeded" {
		logs.Info("Error : ", errredeem)
		logs.Info("[UltraVoucherServices]-[RedeemVoucher]")
		logs.Info("[Limit exceeded]-[Gagal Redeem Voucher]")

		// sugarLogger.Info("Internal Server Error : ", errredeem)
		sugarLogger.Info("[UltraVoucherServices]-[RedeemVoucher]")
		sugarLogger.Info("[Limit exceeded]-[Gagal Redeem Voucher]")

		res = utils.GetMessageResponse(res, 500, false, errors.New("Voucher Sudah Limit"))
		// res.Data = "Transaksi Gagal"

		return res
	}

	if errredeem != nil || redeem.Error != "" {
		logs.Info("Error : ", errredeem)
		logs.Info("[UltraVoucherServices]-[RedeemVoucher]")
		logs.Info("[Failed Redeem Voucher]-[Gagal Redeem Voucher]")

		// sugarLogger.Info("Internal Server Error : ", errredeem)
		sugarLogger.Info("[UltraVoucherServices]-[RedeemVoucher]")
		sugarLogger.Info("[Failed Redeem Voucher]-[Gagal Redeem Voucher]")

		res = utils.GetMessageResponse(res, 500, false, errors.New("Gagal Redeem Voucher"))
		// res.Data = "Transaksi Gagal"

		return res
	}

	nama := "OTTOPOINT"
	expired, _ := strconv.Atoi(param.ExpDate)
	reqOrder := uvmodels.OrderVoucherReq{
		Sku:               param.ProductCode,
		Qty:               req.Jumlah,
		AccountID:         param.CustID,
		InstitutionRefno:  param.Reffnum,
		ExpireDateVoucher: expired,
		ReceiverName:      nama,
		ReceiverEmail:     dataorder.Email,
		ReceiverPhone:     dataorder.Phone,
	}

	// order to u
	order, errOrder := uv.OrderVoucher(reqOrder, param.InstitutionID)
	if errOrder != nil || order.ResponseCode == "" || order.ResponseCode == "01" {
		logs.Info("Error : ", errOrder)
		logs.Info("ResponseCode : ", order.ResponseCode)
		logs.Info("[UltraVoucherServices]-[OrderVoucher]")
		logs.Info("[Failed Order Voucher]-[Gagal Order Voucher]")

		// sugarLogger.Info("Internal Server Error : ", errOrder)
		sugarLogger.Info("[UltraVoucherServices]-[OrderVoucher]")
		sugarLogger.Info("[Failed Order Voucher]-[Gagal Order Voucher]")

		for i := req.Jumlah; i > 0; i-- {

			logs.Info("[Line : %v]", i)

			t := i - 1
			CouponID := redeem.Coupons[t].Id
			couponCode := redeem.Coupons[t].Code

			_, err2 := opl.CouponVoucherCustomer(req.CampaignID, CouponID, couponCode, param.CustID, 1)
			if err2 != nil {
				// res = utils.GetMessageResponse(res, 400, false, errors.New("Gagal Redeem Voucher, Harap coba lagi"))
				// return res

				logs.Info("[UltraVoucherServices]-[CouponVoucherCustomer]")
				logs.Info("[UltraVoucherServices]-[Error : %v]", err2)
				sugarLogger.Info("[UltraVoucherServices]-[CouponVoucherCustomer]")
			}
		}

		Text := "OP009 - " + "Reversal point cause transaction " + param.NamaVoucher + " is failed"
		point := param.Point * req.Jumlah
		totalPoint := strconv.Itoa(point)
		_, errReversal := host.TransferPoint(param.CustID, totalPoint, Text)
		if errReversal != nil {
			logs.Info("Internal Server Error : ", errReversal)
			logs.Info("[UltraVoucherServices]-[TransferPoint]")
			logs.Info("[Failed Order Voucher]-[Gagal Reversal Point]")
		}

		res = utils.GetMessageResponse(res, 500, false, errors.New("Gagal Redeem Voucher"))
		// res.Data = "Transaksi Gagal"

		return res
	}

	if order.ResponseCode == "02" {
		logs.Info("Internal Server Error : ", errOrder)
		logs.Info("ResponseCode : ", order.ResponseCode)
		logs.Info("[UltraVoucherServices]-[OrderVoucher]")
		logs.Info("[Stock not Available]-[Gagal Order Voucher]")

		// sugarLogger.Info("Internal Server Error : ", errOrder)
		sugarLogger.Info("[UltraVoucherServices]-[OrderVoucher]")
		sugarLogger.Info("[Stock not Available]-[Gagal Order Voucher]")

		for i := req.Jumlah; i > 0; i-- {
			logs.Info("[Line : %v]", i)

			t := i - 1
			CouponID := redeem.Coupons[t].Id
			couponCode := redeem.Coupons[t].Code

			_, err2 := opl.CouponVoucherCustomer(req.CampaignID, CouponID, couponCode, param.CustID, 1)
			if err2 != nil {
				// res = utils.GetMessageResponse(res, 400, false, errors.New("Gagal Redeem Voucher, Harap coba lagi"))
				// return res

				logs.Info("[UltraVoucherServices]-[CouponVoucherCustomer]")
				logs.Info("[UltraVoucherServices]-[Error : %v]", err2)
				sugarLogger.Info("[UltraVoucherServices]-[CouponVoucherCustomer]")
			}
		}

		Text := "OP009 - " + "Reversal point cause transaction " + param.NamaVoucher + " is failed"
		point := param.Point * req.Jumlah
		totalPoint := strconv.Itoa(point)
		_, errReversal := host.TransferPoint(param.CustID, totalPoint, Text)
		if errReversal != nil {
			logs.Info("Internal Server Error : ", errReversal)
			logs.Info("[UltraVoucherServices]-[TransferPoint]")
			logs.Info("[Failed Order Voucher]-[Gagal Reversal Point]")
		}

		res = utils.GetMessageResponse(res, 145, false, errors.New(fmt.Sprintf("Voucher yg tersedia %v", order.Data.VouchersAvailable)))
		// res.Data = "Stok Tidak Tersedia"

		return res
	}

	for i := req.Jumlah; i > 0; i-- {

		logs.Info("[Line Save DB : %v]", i)

		t := i - 1
		coupon := redeem.Coupons[t].Id
		code := order.Data.VouchersCode[t].Code

		id := utils.GenerateTokenUUID()
		go SaveDB(id, param.InstitutionID, coupon, code, param.AccountNumber, param.CustID)
	}

	// go SaveTransactionUV(param, order, reqOrder, req, "Payment", "00", order.ResponseCode)

	logs.Info("ResponseCode : ", order.ResponseCode)
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

func SaveDB(id, institution, coupon, vouchercode, phone, custIdOPL string) {
	logs.Info("[SaveDB]-[UltraVoucherServices]")
	save := dbmodels.UserMyVocuher{
		ID:            id,
		InstitutionID: institution,
		CouponID:      coupon,
		VoucherCode:   vouchercode,
		Phone:         phone,
		AccountId:     custIdOPL,
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
