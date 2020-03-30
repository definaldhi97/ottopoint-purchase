package controllers

import (
	"errors"
	"fmt"

	"ottopoint-purchase/constants"
	"ottopoint-purchase/models"
	"ottopoint-purchase/services"
	"ottopoint-purchase/utils"
	"time"

	opl "ottopoint-purchase/hosts/opl/host"
	modelsopl "ottopoint-purchase/hosts/opl/models"
	token "ottopoint-purchase/hosts/redis_token/host"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	zaplog "github.com/opentracing-contrib/go-zap/log"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	ottologer "ottodigital.id/library/logger"
	utilsgo "ottodigital.id/library/utils"

	"net/http"
)

func UseVouhcerController(ctx *gin.Context) {
	req := models.UseVoucherReq{}
	res := models.Response{}

	sugarLogger := ottologer.GetLogger()
	namectrl := "[UseVouhcerController]"

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

	//validate request
	header, resultValidate := ValidateRequest(ctx, true, req)
	if !resultValidate.Meta.Status {
		ctx.JSON(http.StatusOK, resultValidate)
		return
	}

	dataToken, _ := token.CheckToken(header)

	spanid := utilsgo.GetSpanId(span)
	sugarLogger.Info("REQUEST:", zap.String("SPANID", spanid), zap.String("CTRL", namectrl),
		zap.Any("BODY", req),
		zap.Any("HEADER", ctx.Request.Header))

	usevoucher := services.UseVoucherServices{
		General: models.GeneralModel{
			ParentSpan: span,
			OttoZaplog: sugarLogger,
			SpanId:     spanid,
			Context:    context,
		},
	}

	cekVoucher, errVoucher := opl.HistoryVoucherCustomer(dataToken.Data, "")

	data := switchData(cekVoucher.Campaigns, req.CampaignID, req.Category)

	if errVoucher != nil || data.NamaVoucher == "" {
		sugarLogger.Info("[HistoryVoucherCustomer]-[UseVouhcerController]")
		sugarLogger.Info(fmt.Sprintf("Error : ", errVoucher))

		logs.Info("[HistoryVoucherCustomer]-[UseVouhcerController]")
		logs.Info(fmt.Sprintf("Error : ", errVoucher))

		res = utils.GetMessageResponse(res, 422, false, errors.New("Gagal Get History Voucher Customer"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	logs.Info("SupplierID : ", data.SupplierID)
	logs.Info("producrType : ", data.ProductType)

	sugarLogger.Info("=== SupplierID ===")
	sugarLogger.Info(data.SupplierID)

	sugarLogger.Info("=== producrType ===")
	sugarLogger.Info(data.ProductType)

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
		Point:         data.Point,
	}

	switch data.SupplierID {
	case constants.UV:
		res = usevoucher.UseVoucherUV(req, param)
	case constants.OttoAG:
		res = usevoucher.UseVoucherOttoAG(req, param)
	}

	sugarLogger.Info("RESPONSE:", zap.String("SPANID", spanid), zap.String("CTRL", namectrl),
		zap.Any("BODY", res))

	datalog := utils.LogSpanMax(res)
	zaplog.InfoWithSpan(span, namectrl,
		zap.Any("RESP", datalog),
		zap.Duration("backoff", time.Second))

	defer span.Finish()
	ctx.JSON(http.StatusOK, res)

}

func switchData(data []modelsopl.CampaignsDetail, CampaignID, product string) models.Params {
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
	var point int
	for _, valco := range resp {
		nama = valco.Name
		point = valco.CostInPoints
		couponId = valco.Coupon.ID
		couponCode = valco.Coupon.Code
		expDate = valco.ActiveTo
	}

	supplierid := couponCode[:2]
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
		ProductCode: couponCode,
		SupplierID:  supplierID,
		CouponID:    couponId,
		NamaVoucher: nama,
		ExpDate:     expDate,
		Point:       point,
	}

	return res
}
