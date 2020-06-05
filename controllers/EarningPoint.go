package controllers

import (
	"fmt"
	"ottopoint-purchase/constants"
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

func EarningsPointController(ctx *gin.Context) {
	req := models.EarningReq{}
	res := models.Response{}

	sugarLogger := ottologer.GetLogger()
	namectrl := "[EarningsPointController]"

	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.Meta.Code = 03
		res.Meta.Message = "Gagal! Maaf transaksi Anda tidak dapat dilakukan saat ini. Silahkan dicoba lagi atau hubungi tim kami untuk informasi selengkapnya."
		ctx.JSON(http.StatusOK, res)
		go sugarLogger.Error("Error, body Request", zap.Error(err))
		return
	}

	span := TracingFirstControllerCtx(ctx, req, namectrl)
	c := ctx.Request.Context()
	context := opentracing.ContextWithSpan(c, span)

	// validate request
	header, resultValidate := ValidateRequest(ctx, true, req)
	if !resultValidate.Meta.Status {
		ctx.JSON(http.StatusOK, resultValidate)
		return
	}

	// dataToken, _ := token.CheckToken(header)

	spanid := utilsgo.GetSpanId(span)
	sugarLogger.Info("REQUEST:", zap.String("SPANID", spanid), zap.String("CTRL", namectrl),
		zap.Any("BODY", req),
		zap.Any("HEADER", ctx.Request.Header))

	earningPoint := services.EarningPointServices{
		General: models.GeneralModel{
			ParentSpan: span,
			OttoZaplog: sugarLogger,
			SpanId:     spanid,
			Context:    context,
		},
	}

	fmt.Println("Request : ", req)
	fmt.Println("Code : ", req.Earning)

	code := req.Earning[:2]
	switch code {
	case constants.GeneralSpending:
		res = earningPoint.GeneralSpendingService(req, header.InstitutionID)
	// case constants.Multiply        :
	// 	res = earningPoint.GeneralSpendingService(req, header.InstitutionID)
	case constants.InstantReward:
		res = earningPoint.InstantRewardService(req, header.InstitutionID)
	case constants.EventRule:
		res = earningPoint.EventRuleService(req, header.InstitutionID)
	case constants.CustomerReferral:
		res = earningPoint.CustomerReferralService(req, header.InstitutionID)
	case constants.CustomeEventRule:
		res = earningPoint.CustomeEventRuleService(req, header.InstitutionID)
	default:
		// belum ada response

	}
	// res = earningRule.EarningsPointServuc(req, dataToken, header)

	sugarLogger.Info("RESPONSE:", zap.String("SPANID", spanid), zap.String("CTRL", namectrl),
		zap.Any("BODY", res))

	datalog := utils.LogSpanMax(res)
	zaplog.InfoWithSpan(span, namectrl,
		zap.Any("RESP", datalog),
		zap.Duration("backoff", time.Second))

	defer span.Finish()
	ctx.JSON(http.StatusOK, res)

}
