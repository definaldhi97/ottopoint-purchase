package controllers

import (
	"errors"
	"fmt"

	"ottopoint-purchase/models"
	"ottopoint-purchase/services"
	"ottopoint-purchase/utils"
	"time"

	"ottopoint-purchase/db"
	opl "ottopoint-purchase/hosts/opl/host"
	modelsopl "ottopoint-purchase/hosts/opl/models"
	token "ottopoint-purchase/hosts/redis_token/host"
	signature "ottopoint-purchase/hosts/signature/host"

	"github.com/astaxie/beego/logs"

	"github.com/gin-gonic/gin"
	zaplog "github.com/opentracing-contrib/go-zap/log"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	ottologer "ottodigital.id/library/logger"
	utilsgo "ottodigital.id/library/utils"

	"net/http"
)

func UseVouhcerUVController(ctx *gin.Context) {
	req := models.UseVoucherUVReq{}
	res := models.Response{}

	sugarLogger := ottologer.GetLogger()
	namectrl := "[UseVouhcerUVController]"

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

	// jsonSignature, _ := json.Marshal(req)

	ValidateSignature, errSignature := signature.Signature(req, header)
	if errSignature != nil || ValidateSignature.ResponseCode == "" {
		sugarLogger.Info("[ValidateSignature]-[UseVouhcerUVController]")
		sugarLogger.Info(fmt.Sprintf("Error when validation request header"))

		logs.Info("[ValidateSignature]-[UseVouhcerUVController]")
		logs.Info(fmt.Sprintf("Error when validation request header"))

		res = utils.GetMessageResponse(res, 400, false, errors.New("Silahkan login kembali"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	dataToken, errToken := token.CheckToken(header)
	if errToken != nil || dataToken.ResponseCode != "00" {
		sugarLogger.Info("[ValidateToken]-[UseVouhcerUVController]")
		sugarLogger.Info(fmt.Sprintf("Error when validation request header"))

		logs.Info("[ValidateToken]-[UseVouhcerUVController]")
		logs.Info(fmt.Sprintf("Error when validation request header"))

		res = utils.GetMessageResponse(res, 400, false, errors.New("Silahkan login kembali"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	getData, errData := db.GetUltraVoucher(req.VoucherCode, req.AccountNumber)
	if errData != nil || getData.CampaignID == "" {
		logs.Info("Internal Server Error : ", errData)
		logs.Info("[UseVouhcerUVController]-[GetUltraVoucher]")
		logs.Info("[Failed Redeem Voucher]-[Get Data User]")

		// sugarLogger.Info("Internal Server Error : ", errredeem)
		sugarLogger.Info("[UseVouhcerUVController]-[GetUltraVoucher]")
		sugarLogger.Info("[Failed Redeem Voucher]-[Get Data User]")

		res = utils.GetMessageResponse(res, 500, false, errors.New("User belum Eligible"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	spanid := utilsgo.GetSpanId(span)
	sugarLogger.Info("REQUEST:", zap.String("SPANID", spanid), zap.String("CTRL", namectrl),
		zap.Any("BODY", req),
		zap.Any("HEADER", ctx.Request.Header))

	usevoucher := services.UseVoucherUVServices{
		General: models.GeneralModel{
			ParentSpan: span,
			OttoZaplog: sugarLogger,
			SpanId:     spanid,
			Context:    context,
		},
	}

	cekVoucher, errVoucher := opl.VoucherDetail(getData.CampaignID)
	if errVoucher != nil || cekVoucher.CampaignID == "" {
		sugarLogger.Info("[VoucherDetail]-[UseVouhcerUVController]")
		sugarLogger.Info(fmt.Sprintf("Error : ", errVoucher))

		logs.Info("[VoucherDetail]-[UseVouhcerUVController]")
		logs.Info(fmt.Sprintf("Error : ", errVoucher))

		res = utils.GetMessageResponse(res, 422, false, errors.New("Internal Server Error"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	data := switchDataUV(cekVoucher)

	logs.Info("SupplierID : ", data.SupplierID)
	logs.Info("producrType : ", data.ProductType)

	sugarLogger.Info("=== SupplierID ===")
	sugarLogger.Info(data.SupplierID)

	sugarLogger.Info("=== producrType ===")
	sugarLogger.Info(data.ProductType)

	param := models.Params{
		AccountNumber: req.AccountNumber,
		MerchantID:    dataToken.MerchantID,
		InstitutionID: header.InstitutionID,
		SupplierID:    data.SupplierID,
		ProductType:   data.ProductType,
		ProductCode:   data.ProductCode,
		NamaVoucher:   data.NamaVoucher,
		Category:      data.Category,
		CouponID:      getData.CouponID,
		Point:         data.Point,
		CustID:        getData.AccountId,
	}

	res = usevoucher.UseVoucherUV(req, param, getData.CampaignID)

	sugarLogger.Info("RESPONSE:", zap.String("SPANID", spanid), zap.String("CTRL", namectrl),
		zap.Any("BODY", res))

	datalog := utils.LogSpanMax(res)
	zaplog.InfoWithSpan(span, namectrl,
		zap.Any("RESP", datalog),
		zap.Duration("backoff", time.Second))

	defer span.Finish()
	ctx.JSON(http.StatusOK, res)

}

func switchDataUV(data modelsopl.VoucherDetailResp) models.Params {

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

	res = models.Params{
		ProductType: data.BrandName,
		ProductCode: coupon,
		SupplierID:  supplierID,
		// CouponID:    couponId,
		NamaVoucher: data.Name,
		ExpDate:     data.CampaignActivity.ActiveTo,
		Point:       data.CostInPoints,
		Category:    data.BrandName,
	}

	return res
}
