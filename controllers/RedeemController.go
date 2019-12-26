package controllers

import (
	"errors"
	"fmt"

	"ottopoint-purchase/utils"
	"time"

	ottomart "ottopoint-purchase/hosts/ottomart/host"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	zaplog "github.com/opentracing-contrib/go-zap/log"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	ottologer "ottodigital.id/library/logger"
	utilsgo "ottodigital.id/library/utils"

	"net/http"
	"ottopoint-purchase/constants"

	"ottopoint-purchase/models"
	"ottopoint-purchase/services"
)

func Redeem(ctx *gin.Context) {
	req := models.RedeemVoucherRequest{}
	res := models.Response{}

	sugarLogger := ottologer.GetLogger()
	namectrl := "[Redeem-Voucher]"

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
		DeviceID:      ctx.Request.Header.Get("Device-Id"),
		Authorization: ctx.Request.Header.Get("Authorization"),
	}

	dataToken, errToken := ottomart.CheckToken(header)
	if errToken != nil || dataToken.Data.AccountNumber == "" {
		sugarLogger.Info("[ValidateToken]-[controllers-RedeemController]")
		sugarLogger.Info(fmt.Sprintf("Error when validation request header"))

		logs.Info("[ValidateToken]-[controllers-RedeemController]")
		logs.Info(fmt.Sprintf("Error when validation request header"))

		res = utils.GetMessageResponse(res, 400, false, errors.New("Session aplikasi anda telah berakhir"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	spanid := utilsgo.GetSpanId(span)
	sugarLogger.Info("REQUEST:", zap.String("SPANID", spanid), zap.String("CTRL", namectrl),
		zap.Any("BODY", req),
		zap.Any("HEADER", ctx.Request.Header))

	redeemGame := services.RedeemGameService{
		General: models.GeneralModel{
			ParentSpan: span,
			OttoZaplog: sugarLogger,
			SpanId:     spanid,
			Context:    context,
		},
	}

	logs.Info("[Redeem-Voucher : %v]", res)

	switch req.Category {
	case constants.CategoryPulsa:
		res = redeemGame.RedeemGame(req, dataToken)
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
