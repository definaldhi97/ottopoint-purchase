package Controller

import (
	"net/http"
	"ottopoint-purchase/constants"
	controller_v1 "ottopoint-purchase/controllers"
	token "ottopoint-purchase/hosts/redis_token/host"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"

	"github.com/astaxie/beego/logs"

	service_v2_1 "ottopoint-purchase/services/v2.1/voucher"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	ottologer "ottodigital.id/library/logger"
	utilsgo "ottodigital.id/library/utils"
)

func ViewVoucherController(ctx *gin.Context) {
	logs.Info("[ View Voucher Controller ]")

	namectrl := "[ ViewVoucherController ]"
	sugarLogger := ottologer.GetLogger()

	var resp models.Response

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

	// get param
	keys, ok := ctx.Request.URL.Query()["couponId"]
	if !ok || len(keys[0]) < 1 {
		logs.Info("Param 'key' is missing")
		resp = utils.GetMessageFailedErrorNew(resp, constants.RC_PARAMETER_INVALID, constants.RD_PARAMETER_INVALID)
		ctx.JSON(http.StatusOK, resp)
		// go sugarLogger.Error("Error, body Request", zap.Error(ok))
		return
	}
	couponId := keys[0]

	span := controller_v1.TracingFirstControllerCtx(ctx, couponId, namectrl)
	c := ctx.Request.Context()
	context := opentracing.ContextWithSpan(c, span)

	// //validate request
	// header, resultValidate := controller_v1.ValidateRequest(ctx, false, couponId)
	// if !resultValidate.Meta.Status {
	// 	ctx.JSON(http.StatusOK, resultValidate)
	// 	// go sugarLogger.Error("Error, body Request", zap.Error(err))
	// 	return
	// }

	authorization, errAuth := token.CheckToken(header)
	if errAuth != nil {
		logs.Info("Internal server error")
		logs.Info("Chek Token : ", errAuth)
		resp = utils.GetMessageFailedErrorNew(resp, 500, "Internal Server Error")
		ctx.JSON(http.StatusOK, resp)
		return
	}
	if authorization.ResponseCode != "00" {
		resp = utils.GetMessageFailedErrorNew(resp, constants.RC_ERROR_INVALID_TOKEN, constants.RD_ERROR_INVALID_TOKEN)
		ctx.JSON(http.StatusOK, resp)
		return
	}

	spanid := utilsgo.GetSpanId(span)
	sugarLogger.Info("REQUEST:", zap.String("SPANID", spanid), zap.String("CTRL", namectrl),
		zap.Any("BODY", couponId),
		zap.Any("HEADER", ctx.Request.Header))

	viewVoucherSerivice := service_v2_1.ViewVoucherService{
		General: models.GeneralModel{
			ParentSpan: span,
			OttoZaplog: sugarLogger,
			SpanId:     spanid,
			Context:    context,
		},
	}

	// service
	resp = viewVoucherSerivice.ViewVoucher(authorization.Data, couponId)
	sugarLogger.Info("RESPONSE:", zap.String("SPANID", spanid), zap.String("CTRL", namectrl),
		zap.Any("BODY", resp))

	defer span.Finish()
	ctx.JSON(http.StatusOK, resp)
}
