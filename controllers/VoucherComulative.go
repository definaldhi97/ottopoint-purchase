package controllers

import (
	"errors"
	"fmt"

	"ottopoint-purchase/constants"
	opl "ottopoint-purchase/hosts/opl/host"
	modelsopl "ottopoint-purchase/hosts/opl/models"
	"ottopoint-purchase/models"
	"ottopoint-purchase/services"
	"ottopoint-purchase/utils"
	"time"

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
		sugarLogger.Info("[ValidateSignature]-[VoucherComulative-Controller]")
		sugarLogger.Info(fmt.Sprintf("Error when validation request header"))

		logs.Info("[ValidateSignature]-[VoucherComulative-Controller]")
		logs.Info(fmt.Sprintf("Error when validation request header"))

		res = utils.GetMessageResponse(res, 400, false, errors.New("Signature salah"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	dataToken, errToken := token.CheckToken(header)
	if errToken != nil || dataToken.ResponseCode != "00" {
		sugarLogger.Info("[ValidateToken]-[VoucherComulative-Controller]")
		sugarLogger.Info(fmt.Sprintf("Error when validation request header"))
		sugarLogger.Info(fmt.Sprintf("Error : ", errToken))

		logs.Info("[ValidateToken]-[VoucherComulative-Controller]")
		logs.Info(fmt.Sprintf("Error when validation request header"))
		logs.Info(fmt.Sprintf("Error : ", errToken))

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

	cekVoucher, errVoucher := opl.HistoryVoucherCustomer(dataToken.Data, "")
	if errVoucher != nil || cekVoucher.Campaigns[0].CampaignID == "" {
		sugarLogger.Info("[HistoryVoucherCustomer]-[VoucherComulative-Controller]")
		sugarLogger.Info(fmt.Sprintf("Error : ", errVoucher))

		logs.Info("[HistoryVoucherCustomer]-[VoucherComulative-Controller]")
		logs.Info(fmt.Sprintf("Error : ", errVoucher))

		res = utils.GetMessageResponse(res, 422, false, errors.New("Internal Server Error"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	data := SwitchCheckData(cekVoucher.Campaigns, req.Category, req.CampaignID)

	logs.Info("SupplierID : ", data.SupplierID)
	logs.Info("producrType : ", data.ProductType)

	sugarLogger.Info("SupplierID : ", data.SupplierID)
	sugarLogger.Info("producrType : ", data.ProductType)

	param := models.Params{
		AccountNumber: dataToken.Data,
		MerchantID:    dataToken.MerchantID,
		InstitutionID: header.InstitutionID,
		SupplierID:    data.SupplierID,
		ProductType:   data.ProductType,
		ProductCode:   data.ProductCode,
		NamaVoucher:   data.NamaVoucher,
		Category:      req.Category,
		CouponID:      data.CouponID,
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

func SwitchCheckData(data []modelsopl.CampaignsDetail, product, CampaignID string) models.Params {
	res := models.Params{}

	resp := []models.CampaignsDetail{}
	for _, val := range data {
		if val.CampaignID == CampaignID && val.CanBeUsed == true {
			a := models.CampaignsDetail{
				Name:       val.Campaign.Name,
				CampaignID: val.CampaignID,
				ActiveTo:   val.ActiveTo,
				Coupon: models.CouponDetail{
					Code: val.Coupon.Code,
					ID:   val.Coupon.ID,
				},
			}

			resp = append(resp, a)
		}
	}

	var couponId, couponCode, nama, expDate string
	for _, valco := range resp {
		nama = valco.Name
		couponId = valco.Coupon.ID
		couponCode = valco.Coupon.Code
		expDate = valco.ActiveTo
	}

	supplierid := couponCode[2:]
	var supplierID string
	if supplierid == "UV" {
		supplierID = "Ultra Voucher"
	} else {
		supplierID = "OttoAG"
	}

	var producrType string
	switch product {
	case constants.CategoryPulsa:
		producrType = "Pulsa"
	case constants.CategoryFreeFire, constants.CategoryMobileLegend:
		producrType = "Game"
	case constants.CategoryToken:
		producrType = "PLN"
	}

	res = models.Params{
		ProductType: producrType,
		SupplierID:  supplierID,
		NamaVoucher: nama,
		CouponID:    couponId,
		ExpDate:     expDate,
	}

	return res
}
