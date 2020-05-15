package controllers

import (
	token "ottopoint-purchase/hosts/redis_token/host"
	"ottopoint-purchase/services"
	"ottopoint-purchase/utils"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	zaplog "github.com/opentracing-contrib/go-zap/log"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	ottologer "ottodigital.id/library/logger"
	utilsgo "ottodigital.id/library/utils"

	"net/http"

	"ottopoint-purchase/models"
)

func EarningRuleController(ctx *gin.Context) {
	req := models.EarningRuleReq{}
	res := models.Response{}

	sugarLogger := ottologer.GetLogger()
	namectrl := "[EarningRuleController]"

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

	earningRule := services.EarningRuleServices{
		General: models.GeneralModel{
			ParentSpan: span,
			OttoZaplog: sugarLogger,
			SpanId:     spanid,
			Context:    context,
		},
	}

	logs.Info("Request : ", req)
	res = earningRule.NewEarningRuleServices(req, dataToken, header)

	sugarLogger.Info("RESPONSE:", zap.String("SPANID", spanid), zap.String("CTRL", namectrl),
		zap.Any("BODY", res))

	datalog := utils.LogSpanMax(res)
	zaplog.InfoWithSpan(span, namectrl,
		zap.Any("RESP", datalog),
		zap.Duration("backoff", time.Second))

	defer span.Finish()
	ctx.JSON(http.StatusOK, res)

}
