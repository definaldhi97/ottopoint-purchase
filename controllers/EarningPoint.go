package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"ottopoint-purchase/constants"
	kafka "ottopoint-purchase/hosts/publisher/host"
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

	header := models.RequestHeader{}
	header.InstitutionID = "PSM0001"
	// validate request
	// header, resultValidate := ValidateRequest(ctx, true, req)
	// if !resultValidate.Meta.Status {
	// 	ctx.JSON(http.StatusOK, resultValidate)
	// 	return
	// }

	// dataToken, _ := token.CheckToken(header)
	// fmt.Println(dataToken)

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

	fmt.Println(earningPoint)
	fmt.Println(fmt.Sprintf("[Request : %v]", req))
	fmt.Println(fmt.Sprintf("[Code : %v]", req.Earning))

	res = utils.GetMessageResponse(res, 200, true, errors.New("Transaksi sedang di proses"))

	code := req.Earning[:3]
	switch code {
	case constants.GeneralSpending:
		fmt.Println("===== GeneralSpending =====")
		go publishEarning(req, header)
		// res = earningPoint.GeneralSpendingService(req, header.InstitutionID)
	// case constants.Multiply        :
	// 	res = earningPoint.GeneralSpendingService(req, header.InstitutionID)
	case constants.InstantReward:
		fmt.Println("===== InstantReward =====")
		go publishEarning(req, header)
		// res = earningPoint.InstantRewardService(req, header.InstitutionID)
	case constants.EventRule:
		fmt.Println("===== EventRule =====")
		go publishEarning(req, header)
		// res = earningPoint.EventRuleService(req, header.InstitutionID)
	case constants.CustomerReferral:
		fmt.Println("===== CustomerReferral =====")
		go publishEarning(req, header)
		// res = earningPoint.CustomerReferralService(req, header.InstitutionID)
	case constants.CustomeEventRule:
		fmt.Println("===== CustomeEventRule =====")
		go publishEarning(req, header)
		// res = earningPoint.CustomeEventRuleService(req, header.InstitutionID)
	default:
		fmt.Println("===== Invalid Code =====")
		res = utils.GetMessageResponse(res, 178, false, errors.New("Earning Rule not found"))
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

func publishEarning(req models.EarningReq, header models.RequestHeader) {
	fmt.Println(">>>>> Publisher Earning <<<<<")

	pubReq := models.PublishEarningReq{
		Header:         header,
		Earning:        req.Earning,
		ReferenceId:    req.ReferenceId,
		ProductCode:    req.ProductCode,
		ProductName:    req.ProductName,
		AccountNumber1: req.AccountNumber1,
		AccountNumber2: req.AccountNumber2,
		Amount:         req.Amount,
		Remark:         req.Remark,
	}

	bytePub, _ := json.Marshal(pubReq)

	kafkaReq := kafka.PublishReq{
		Topic: "ottopoint-earning-topics",
		Value: bytePub,
	}

	kafkaRes, err := kafka.SendPublishKafka(kafkaReq)
	if err != nil {
		fmt.Println("Gagal Send Publisher")
		fmt.Println("Error : ", err)
	}

	fmt.Println("Response Publisher : ", kafkaRes)
}
