package controllers

import (
	"net/http"
	sepulsaModels "ottopoint-purchase/hosts/sepulsa/models"
	"ottopoint-purchase/models"
	"ottopoint-purchase/services"
	"ottopoint-purchase/utils"
	"time"

	utilsgo "ottodigital.id/library/utils"

	"github.com/gin-gonic/gin"
	zaplog "github.com/opentracing-contrib/go-zap/log"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	ottologer "ottodigital.id/library/logger"
)

func HandleCallbackSepulsa(ctx *gin.Context) {

	req := sepulsaModels.CallbackTrxReq{}
	res := models.Response{}

	sugarLogger := ottologer.GetLogger()
	namectrl := "[HandleCallbackSepulsa]"

	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.Meta.Code = 03
		res.Meta.Message = err.Error()
		ctx.JSON(http.StatusOK, res)
		go sugarLogger.Error("Error, body request", zap.Error(err))
		return
	}

	span := TracingFirstControllerCtx(ctx, req, namectrl)
	c := ctx.Request.Context()
	context := opentracing.ContextWithSpan(c, span)
	spanId := utilsgo.GetSpanId(span)

	sepulsaSvc := services.UseSepulsaService{
		General: models.GeneralModel{
			ParentSpan: span,
			OttoZaplog: sugarLogger,
			SpanId:     spanId,
			Context:    context,
		},
	}

	res = sepulsaSvc.HandleCallbackRequest(req)

	sugarLogger.Info("RESPONSE : ", zap.String("SPANID", spanId), zap.String("CTRL", namectrl),
		zap.Any("BODY : ", res))

	datalog := utils.LogSpanMax(res)
	zaplog.InfoWithSpan(span, namectrl,
		zap.Any("RESP : ", datalog),
		zap.Duration("backoff : ", time.Second))

	defer span.Finish()

	ctx.JSON(http.StatusOK, res)
}
