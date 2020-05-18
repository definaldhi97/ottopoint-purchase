package controllers

import (
	"fmt"
	"time"

	"ottopoint-purchase/db"
	token "ottopoint-purchase/hosts/redis_token/host"
	"ottopoint-purchase/services"
	"ottopoint-purchase/utils"

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

func ReversePointController(ctx *gin.Context) {
	req := models.ReversePointReq{}
	res := models.Response{}

	reqDeduct := models.PointReq{}

	sugarLogger := ottologer.GetLogger()
	namectrl := "[ReversePoint-Controller]"

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

	transferPoint := services.TransferPointServices{
		General: models.GeneralModel{
			ParentSpan: span,
			OttoZaplog: sugarLogger,
			SpanId:     spanid,
			Context:    context,
		},
	}
	dataDeduct, err := db.GetDataDeduct(req.TrxID)
	if err != nil {
		sugarLogger.Info("[CheckDeduction]-[controllers-ReversePoint]")
		sugarLogger.Info(fmt.Sprintf("Error when get data deduction"))

		logs.Info("[CheckDeduction]-[controllers-ReversePoint]")
		logs.Info(fmt.Sprintf("Error when get data deduction"))

		res = utils.GetMessageResponse(res, 400, false, err)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}
	// reqDeduct.AccountNumber = dataDeduct.CustomerID
	reqDeduct.Point = dataDeduct.Point + int(dataDeduct.Amount)
	reqDeduct.Text = "Reversal " + dataDeduct.ProductName
	reqDeduct.Type = "adding"

	res = transferPoint.NewTransferPointServices(reqDeduct, dataToken, header)
	_, err = db.UpdateDataDeduct(dataDeduct.TrxID)
	if err != nil {
		sugarLogger.Info("[UpdateStatusDeduction]-[controllers-ReversePoint]")
		sugarLogger.Info(fmt.Sprintf("Error when update status deduction"))

		logs.Info("[UpdateStatusDeduction]-[controllers-ReversePoint]")
		logs.Info(fmt.Sprintf("Error when update status deduction"))
		res = utils.GetMessageResponse(res, 400, false, err)
		ctx.JSON(http.StatusOK, res)
		return
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
