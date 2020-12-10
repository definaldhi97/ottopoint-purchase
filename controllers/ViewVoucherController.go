package controllers

import (
	"fmt"
	"net/http"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/models"
	"ottopoint-purchase/services"
	"ottopoint-purchase/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	ottologer "ottodigital.id/library/logger"
	utilsgo "ottodigital.id/library/utils"
)

func ViewVoucherController(ctx *gin.Context) {
	fmt.Println("[ >>>>>>>>>>>>>>>>>>>>> View Voucher Controller <<<<<<<<<<<<<<<<<< ]")

	namectrl := "[ ViewVoucherController ]"
	sugarLogger := ottologer.GetLogger()

	var resp models.Response

	couponId := ctx.Request.URL.Query().Get("couponId")

	// header
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

	//check header request
	if header.AppsID == "" || header.ChannelID == "" || header.InstitutionID == "" || header.DeviceID == "" || header.Geolocation == "" {
		resp = utils.GetMessageFailedErrorNew(resp, constants.RC_ERROR_HEADER_MANDATORY, constants.RD_ERROR_HEADER_MANDATORY)
		ctx.JSON(http.StatusOK, resp)
		return
	}

	span := TracingFirstControllerCtx(ctx, couponId, namectrl)
	c := ctx.Request.Context()
	context := opentracing.ContextWithSpan(c, span)

	// Validate Token user
	var RedisService = new(services.RedisService)
	auth := strings.ReplaceAll(header.Authorization, "Bearer ", "")
	keyRedis := header.InstitutionID + "-" + auth
	dataRedis := RedisService.GetData(keyRedis)

	if dataRedis.ResponseCode != "00" {
		resp = utils.GetMessageFailedErrorNew(resp, constants.RC_ERROR_INVALID_TOKEN, constants.RD_ERROR_INVALID_TOKEN)
		ctx.JSON(http.StatusOK, resp)
		return
	}

	spanid := utilsgo.GetSpanId(span)
	sugarLogger.Info("REQUEST:", zap.String("SPANID", spanid), zap.String("CTRL", namectrl),
		zap.Any("BODY", couponId),
		zap.Any("HEADER", ctx.Request.Header))

	viewVoucherSerivice := services.ViewVoucherService{
		General: models.GeneralModel{
			ParentSpan: span,
			OttoZaplog: sugarLogger,
			SpanId:     spanid,
			Context:    context,
		},
	}

	// service
	resp = viewVoucherSerivice.ViewVoucher(dataRedis.Value, couponId)
	sugarLogger.Info("RESPONSE:", zap.String("SPANID", spanid), zap.String("CTRL", namectrl),
		zap.Any("BODY", resp))

	defer span.Finish()
	ctx.JSON(http.StatusOK, resp)

}
