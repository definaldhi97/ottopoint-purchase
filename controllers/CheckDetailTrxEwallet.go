package controllers

import (
	"net/http"
	"ottopoint-purchase/models"
	"ottopoint-purchase/services"
	"ottopoint-purchase/utils"
	"time"

	"github.com/gin-gonic/gin"
	zaplog "github.com/opentracing-contrib/go-zap/log"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	ottologer "ottodigital.id/library/logger"
	utilsgo "ottodigital.id/library/utils"
)

func CheckStatusTrxController(ctx *gin.Context) {

	res := models.Response{}

	transactionID := ctx.Param("transaction_id")

	sugarLogger := ottologer.GetLogger()
	namectrl := "[CheckStatusTrxController]"

	span := TracingFirstControllerCtx(ctx, map[string]interface{}{
		"transaction_id": transactionID,
	}, namectrl)
	c := ctx.Request.Context()
	context := opentracing.ContextWithSpan(c, span)

	spanId := utilsgo.GetSpanId(span)
	sugarLogger.Info("REQUEST:", zap.String("SPANID", spanId), zap.String("CTRL", namectrl),
		zap.Any("HEADER", ctx.Request.Header))

	sepulsaSvc := services.UseSepulsaService{
		General: models.GeneralModel{
			ParentSpan: span,
			OttoZaplog: sugarLogger,
			SpanId:     spanId,
			Context:    context,
		},
	}

	res = sepulsaSvc.CheckStatusTrx(transactionID)

	sugarLogger.Info("RESPONSE : ", zap.String("SPANID", spanId), zap.String("CTRL", namectrl),
		zap.Any("BODY : ", res))

	datalog := utils.LogSpanMax(res)
	zaplog.InfoWithSpan(span, namectrl,
		zap.Any("RESP : ", datalog),
		zap.Duration("backoff : ", time.Second))

	defer span.Finish()

	ctx.JSON(http.StatusOK, res)

}
