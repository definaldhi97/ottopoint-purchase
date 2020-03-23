package controllers

import (
	"errors"
	"fmt"

	"ottopoint-purchase/constants"
	opl "ottopoint-purchase/hosts/opl/host"
	"ottopoint-purchase/models"
	"ottopoint-purchase/services"
	"ottopoint-purchase/utils"
	"time"

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

	//validate request
	header, resultValidate := ValidateRequest(ctx, true, req)
	if !resultValidate.Meta.Status {
		ctx.JSON(http.StatusOK, resultValidate)
		return
	}

	dataToken, _ := token.CheckToken(header)

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

	cekVoucher, errVoucher := opl.VoucherDetail(req.CampaignID)
	if errVoucher != nil || cekVoucher.Name == "" {
		sugarLogger.Info("[GetVoucherDetail]-[VoucherComulative-Controller]")
		sugarLogger.Info(fmt.Sprintf("Error : ", errVoucher))

		logs.Info("[GetVoucherDetail]-[VoucherComulative-Controller]")
		logs.Info(fmt.Sprintf("Error : ", errVoucher))

		res = utils.GetMessageResponse(res, 422, false, errors.New("Internal Server Error"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	var supplier string
	for _, val := range cekVoucher.Coupons {
		supplier = val
	}

	data := SwitchCheckData(supplier, req.Category)

	logs.Info("SupplierID : ", data.SupplierID)
	logs.Info("producrType : ", data.ProductType)

	sugarLogger.Info("SupplierID : ")
	sugarLogger.Info("producrType : ")

	param := models.Params{
		AccountNumber: dataToken.Data,
		MerchantID:    dataToken.MerchantID,
		InstitutionID: header.InstitutionID,
		SupplierID:    data.SupplierID,
		ProductType:   data.ProductType,
		NamaVoucher:   cekVoucher.Name,
		Category:      req.Category,
	}

	switch data.SupplierID {
	// case constants.UV:
	// 	res = voucherComulative.VoucherComulative(req, param)
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

func SwitchCheckData(supplier, product string) models.Params {
	res := models.Params{}

	supplierid := supplier[0:2]
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
	}

	return res
}
