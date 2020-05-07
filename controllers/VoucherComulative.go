package controllers

import (
	"errors"
	"fmt"

	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	opl "ottopoint-purchase/hosts/opl/host"
	modelsopl "ottopoint-purchase/hosts/opl/models"
	token "ottopoint-purchase/hosts/redis_token/host"
	signature "ottopoint-purchase/hosts/signature/host"
	"ottopoint-purchase/models"
	"ottopoint-purchase/services"
	"ottopoint-purchase/utils"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	zaplog "github.com/opentracing-contrib/go-zap/log"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	ottologer "ottodigital.id/library/logger"
	utilsgo "ottodigital.id/library/utils"

	"net/http"
)

func VoucherComulativeController(ctx *gin.Context) {
	req := models.VoucherComultaiveReq{}
	res := models.Response{}

	sugarLogger := ottologer.GetLogger()
	namectrl := "[VoucherComulative-Controller]"

	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.Meta.Code = 03
		res.Meta.Message = "Error, Unmarshall Body Request"
		ctx.JSON(http.StatusBadRequest, res)
		go sugarLogger.Error("Error, body Request", zap.Error(err))
		return
	}

	span := TracingFirstControllerCtx(ctx, req, namectrl)
	c := ctx.Request.Context()
	context := opentracing.ContextWithSpan(c, span)

	header := models.RequestHeader{
		DeviceID:      ctx.Request.Header.Get("DeviceId"),
		InstitutionID: ctx.Request.Header.Get("InstitutionId"),
		Geolocation:   ctx.Request.Header.Get("Geolocation"),
		ChannelID:     ctx.Request.Header.Get("ChannelId"),
		AppsID:        ctx.Request.Header.Get("AppsId"),
		Timestamp:     ctx.Request.Header.Get("Timestamp"),
		Authorization: ctx.Request.Header.Get("Authorization"),
		Signature:     ctx.Request.Header.Get("Signature"),
	}

	ValidateSignature, errSignature := signature.Signature(req, header)
	if errSignature != nil || ValidateSignature.ResponseCode == "" {
		sugarLogger.Info("[ValidateSignature]-[DeductSplitBillController]")
		sugarLogger.Info(fmt.Sprintf("Error when validation request header"))

		logs.Info("[ValidateSignature]-[DeductSplitBillController]")
		logs.Info(fmt.Sprintf("Error when validation request header"))

		res = utils.GetMessageResponse(res, 400, false, errors.New("Silahkan login kembali"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	dataToken, errToken := token.CheckToken(header)
	if errToken != nil || dataToken.ResponseCode != "00" {
		sugarLogger.Info("[ValidateToken]-[DeductSplitBillController]")
		sugarLogger.Info(fmt.Sprintf("Error when validation request header"))

		logs.Info("[ValidateToken]-[DeductSplitBillController]")
		logs.Info(fmt.Sprintf("Error when validation request header"))

		res = utils.GetMessageResponse(res, 400, false, errors.New("Silahkan login kembali"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	spanid := utilsgo.GetSpanId(span)
	sugarLogger.Info("REQUEST : ", zap.String("SPANID", spanid), zap.String("CTRL", namectrl),
		zap.Any("BODY : ", req),
		zap.Any("HEADER : ", ctx.Request.Header))

	voucherComulative := services.VoucherComulativeService{
		General: models.GeneralModel{
			ParentSpan: span,
			OttoZaplog: sugarLogger,
			SpanId:     spanid,
			Context:    context,
		},
	}

	ultraVoucher := services.UseVoucherUltraVoucher{
		General: models.GeneralModel{
			ParentSpan: span,
			OttoZaplog: sugarLogger,
			SpanId:     spanid,
			Context:    context,
		},
	}

	cekVoucher, errVoucher := opl.VoucherDetail(req.CampaignID)
	if errVoucher != nil || cekVoucher.CampaignID == "" {
		sugarLogger.Info("[HistoryVoucherCustomer]-[VoucherComulative-Controller]")
		sugarLogger.Info(fmt.Sprintf("Error : ", errVoucher))

		logs.Info("[HistoryVoucherCustomer]-[VoucherComulative-Controller]")
		logs.Info(fmt.Sprintf("Error : ", errVoucher))

		res = utils.GetMessageResponse(res, 422, false, errors.New("Internal Server Error"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	dataUser, errUser := db.CheckUser(dataToken.Data)
	if errUser != nil || dataUser.CustID == "" {
		logs.Info("Internal Server Error : ", errUser)
		logs.Info("[UltraVoucherServices]-[CheckUser]")
		logs.Info("[Failed Redeem Voucher]-[Get Data User]")

		// sugarLogger.Info("Internal Server Error : ", errredeem)
		sugarLogger.Info("[UltraVoucherServices]-[CheckUser]")
		sugarLogger.Info("[Failed Redeem Voucher]-[Get Data User]")

		res = utils.GetMessageResponse(res, 500, false, errors.New("User belum Eligible"))
	}

	data := SwitchCheckData(cekVoucher)

	logs.Info("SupplierID : ", data.SupplierID)
	logs.Info("producrType : ", data.ProductType)

	// sugarLogger.Info("SupplierID : ", data.SupplierID)
	// sugarLogger.Info("producrType : ", data.ProductType)

	param := models.Params{
		AccountNumber: dataToken.Data,
		MerchantID:    dataToken.MerchantID,
		InstitutionID: header.InstitutionID,
		CustID:        dataUser.CustID,
		SupplierID:    data.SupplierID,
		ProductType:   data.ProductType,
		ProductCode:   data.ProductCode,
		NamaVoucher:   data.NamaVoucher,
		Point:         data.Point,
		Category:      data.Category,
	}

	switch data.SupplierID {
	case constants.UV:
		res = ultraVoucher.UltraVoucherServices(req, param)
	case constants.OttoAG:
		res = voucherComulative.VoucherComulative(req, param)
	}

	sugarLogger.Info("RESPONSE : ", zap.String("SPANID", spanid), zap.String("CTRL", namectrl),
		zap.Any("BODY : ", res))

	datalog := utils.LogSpanMax(res)
	zaplog.InfoWithSpan(span, namectrl,
		zap.Any("RESP : ", datalog),
		zap.Duration("backoff : ", time.Second))

	defer span.Finish()
	ctx.JSON(http.StatusOK, res)

}

func SwitchCheckData(data modelsopl.VoucherDetailResp) models.Params {
	res := models.Params{}

	coupon := data.Coupons[0]

	supplierid := coupon[:2]
	var supplierID string
	if supplierid == "UV" {
		supplierID = "Ultra Voucher"
		coupon = coupon[3:]
	} else {
		supplierID = "OttoAG"
	}

	var producrType string
	switch data.BrandName {
	case constants.CategoryPulsa:
		producrType = "Pulsa"
	case constants.CategoryFreeFire, constants.CategoryMobileLegend:
		producrType = "Game"
	case constants.CategoryPLN:
		producrType = "PLN"
	default:
		producrType = data.BrandName
	}

	res = models.Params{
		ProductType: producrType,
		ProductCode: coupon,
		SupplierID:  supplierID,
		NamaVoucher: data.Name,
		Point:       data.CostInPoints,
		Category:    data.BrandName,
		// ExpDate:     data.CampaignActivity.ActiveTo,
	}

	return res
}
