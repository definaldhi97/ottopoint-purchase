package controllers

import (
	services "ottopoint-purchase/services/earnings"
	"ottopoint-purchase/utils"
	"time"

	"github.com/gin-gonic/gin"
	zaplog "github.com/opentracing-contrib/go-zap/log"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	ottologer "ottodigital.id/library/logger"
	utilsgo "ottodigital.id/library/utils"

	"net/http"

	"ottopoint-purchase/models"
)

func GetEarningRuleController(ctx *gin.Context) {
	res := models.Response{}

	sugarLogger := ottologer.GetLogger()
	namectrl := "[GetEarningRuleController]"

	productCode := ctx.Request.URL.Query().Get("productCode")

	span := TracingFirstControllerCtx(ctx, "", namectrl)
	c := ctx.Request.Context()
	context := opentracing.ContextWithSpan(c, span)

	spanid := utilsgo.GetSpanId(span)
	sugarLogger.Info("REQUEST:", zap.String("SPANID", spanid), zap.String("CTRL", namectrl))

	getEarningRule := services.GetEarningRuleService{
		General: models.GeneralModel{
			ParentSpan: span,
			OttoZaplog: sugarLogger,
			SpanId:     spanid,
			Context:    context,
		},
	}

	res = getEarningRule.NewGetEarningRuleService(productCode)

	sugarLogger.Info("RESPONSE:", zap.String("SPANID", spanid), zap.String("CTRL", namectrl),
		zap.Any("BODY", res))

	datalog := utils.LogSpanMax(res)
	zaplog.InfoWithSpan(span, namectrl,
		zap.Any("RESP", datalog),
		zap.Duration("backoff", time.Second))

	defer span.Finish()
	ctx.JSON(http.StatusOK, res)

}
