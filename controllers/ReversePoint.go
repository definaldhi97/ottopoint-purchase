package controllers

import (
	"errors"
	"fmt"
	"time"

	"ottopoint-purchase/db"
	token "ottopoint-purchase/hosts/redis_token/host"
	signature "ottopoint-purchase/hosts/signature/host"
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

func ReversePoint(ctx *gin.Context) {
	req := models.ReversePointReq{}
	res := models.Response{}
	reqDeduct := models.DeductPointReq{}

	sugarLogger := ottologer.GetLogger()
	namectrl := "[ReversePoint-Controller]"

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
	if errSignature != nil || ValidateSignature.ResponseCode == "" {
		sugarLogger.Info("[ValidateSignature]-[controllers-DeductPoint]")
		sugarLogger.Info(fmt.Sprintf("Error when validation request header"))

		logs.Info("[ValidateSignature]-[controllers-DeductPoint]")
		logs.Info(fmt.Sprintf("Error when validation request header"))

		res = utils.GetMessageResponse(res, 400, false, errors.New("Silahkan login kembali"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	dataToken, errToken := token.CheckToken(header)
	if errToken != nil || dataToken.ResponseCode != "00" {
		sugarLogger.Info("[ValidateToken]-[controllers-DeductPoint]")
		sugarLogger.Info(fmt.Sprintf("Error when validation request header"))

		logs.Info("[ValidateToken]-[controllers-DeductPoint]")
		logs.Info(fmt.Sprintf("Error when validation request header"))

		res = utils.GetMessageResponse(res, 400, false, errors.New("Silahkan login kembali"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
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

		res = utils.GetMessageResponse(res, 400, false, errors.New("Internal Server Error !"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	reqDeduct.AccountNumber = dataDeduct.CustomerID
	reqDeduct.Point = dataDeduct.Point
	reqDeduct.Text = dataDeduct.ProductName
	reqDeduct.Type = "adding"

	res = transferPoint.NewTransferPointServices(reqDeduct, dataToken, header)
	_, err = db.UpdateDataDeduct(dataDeduct.TrxID)
	if err != nil {
		sugarLogger.Info("[UpdateStatusDeduction]-[controllers-ReversePoint]")
		sugarLogger.Info(fmt.Sprintf("Error when update status deduction"))

		logs.Info("[UpdateStatusDeduction]-[controllers-ReversePoint]")
		logs.Info(fmt.Sprintf("Error when update status deduction"))
		res = utils.GetMessageResponse(res, 400, false, errors.New("Internal Server Error !"))
		ctx.JSON(http.StatusBadRequest, res)
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
