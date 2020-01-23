package controllers

import (
	"errors"
	"fmt"

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

	"ottopoint-purchase/models"
	"ottopoint-purchase/services"
)

func VoucherRedeem(ctx *gin.Context) {
	req := models.RedeemReq{}
	res := models.Response{}

	sugarLogger := ottologer.GetLogger()
	namectrl := "[Voucher-Voucher]"

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
	if errSignature != nil || ValidateSignature.ResponseCode != "00" {
		sugarLogger.Info("[ValidateSignature]-[controllers-UseVouhcer]")
		sugarLogger.Info(fmt.Sprintf("Error when validation request header"))

		logs.Info("[ValidateSignature]-[controllers-UseVouhcer]")
		logs.Info(fmt.Sprintf("Error when validation request header"))

		res = utils.GetMessageResponse(res, 400, false, errors.New("Signature salah"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	dataToken, errToken := token.CheckToken(header)
	if errToken != nil || dataToken.ResponseCode != "00" {
		sugarLogger.Info("[ValidateToken]-[controllers-VoucherController]")
		sugarLogger.Info(fmt.Sprintf("Error when validation request header"))

		logs.Info("[ValidateToken]-[controllers-VoucherController]")
		logs.Info(fmt.Sprintf("Error when validation request header"))

		res = utils.GetMessageResponse(res, 400, false, errors.New("Silahkan login kembali"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	logs.Info("Response Token :", dataToken)

	spanid := utilsgo.GetSpanId(span)
	sugarLogger.Info("REQUEST:", zap.String("SPANID", spanid), zap.String("CTRL", namectrl),
		zap.Any("BODY", req),
		zap.Any("HEADER", ctx.Request.Header))

	voucherRedeem := services.VoucherRedeemServices{
		General: models.GeneralModel{
			ParentSpan: span,
			OttoZaplog: sugarLogger,
			SpanId:     spanid,
			Context:    context,
		},
	}

	res = voucherRedeem.VoucherRedeem(req, dataToken.Data)

	sugarLogger.Info("RESPONSE:", zap.String("SPANID", spanid), zap.String("CTRL", namectrl),
		zap.Any("BODY", res))

	datalog := utils.LogSpanMax(res)
	zaplog.InfoWithSpan(span, namectrl,
		zap.Any("RESP", datalog),
		zap.Duration("backoff", time.Second))

	defer span.Finish()
	ctx.JSON(http.StatusOK, res)

}
