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
	token "ottopoint-purchase/hosts/redis_token/host"

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
	if errVoucher != nil || cekVoucher.Campaigns[0].CampaignID == "" {
		sugarLogger.Info("[HistoryVoucherCustomer]-[UseVouhcerController]")
		sugarLogger.Info(fmt.Sprintf("Error : ", errVoucher))

		logs.Info("[HistoryVoucherCustomer]-[UseVouhcerController]")
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
