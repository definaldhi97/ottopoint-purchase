package v2_migrate

import (
	"errors"
	"fmt"

	"ottopoint-purchase/models"
	"ottopoint-purchase/services/v2_migrate"
	"ottopoint-purchase/utils"

	"ottopoint-purchase/controllers"
	"ottopoint-purchase/db"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	ottologer "ottodigital.id/library/logger"
	utilsgo "ottodigital.id/library/utils"
)

func CallBackUVController(ctx *gin.Context) {

	logrus.Info("[ >>>>>>>>>>>>>>>>>>>>> Callbakc Voucher UV COntroller <<<<<<<<<<<<<<<<<< ]")

	req := models.UseVoucherUVReq{}
	res := models.Response{}

	sugarLogger := ottologer.GetLogger()
	namectrl := "[UseVouhcerUVController]"

	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.Meta.Code = 03
		res.Meta.Message = "Gagal! Maaf transaksi Anda tidak dapat dilakukan saat ini. Silahkan dicoba lagi atau hubungi tim kami untuk informasi selengkapnya."
		ctx.JSON(http.StatusOK, res)
		go sugarLogger.Error("Error, body Request", zap.Error(err))
		return
	}

	span := controllers.TracingFirstControllerCtx(ctx, req, namectrl)
	c := ctx.Request.Context()
	context := opentracing.ContextWithSpan(c, span)

	// header := models.RequestHeader{
	// 	// DeviceID:      ctx.Request.Header.Get("DeviceId"),
	// 	InstitutionID: ctx.Request.Header.Get("InstitutionId"),
	// 	// Geolocation:   ctx.Request.Header.Get("Geolocation"),
	// 	// ChannelID:     ctx.Request.Header.Get("ChannelId"),
	// 	// AppsID:        ctx.Request.Header.Get("AppsId"),
	// 	Timestamp: ctx.Request.Header.Get("Timestamp"),
	// 	// Authorization: ctx.Request.Header.Get("Authorization"),
	// 	Signature: ctx.Request.Header.Get("Signature"),
	// }

	// jsonSignature, _ := json.Marshal(req)

	// ValidateSignature, errSignature := signature.Signature(req, header)
	// if errSignature != nil || ValidateSignature.ResponseCode == "" {
	// 	sugarLogger.Info("[ValidateSignature]-[UseVouhcerUVController]")
	// 	sugarLogger.Info(fmt.Sprintf("Error when validation request header"))

	// 	logs.Info("[ValidateSignature]-[UseVouhcerUVController]")
	// 	logs.Info(fmt.Sprintf("Error when validation request header"))

	// 	res = utils.GetMessageResponse(res, 400, false, errors.New("Silahkan login kembali"))
	// 	ctx.JSON(http.StatusOK, res)
	// 	return
	// }

	getData, errData := db.GetUltraVoucher(req.VoucherCode, req.AccountId)
	if errData != nil || getData.CampaignID == "" {
		logrus.Info("[UseVouhcerUVController]-[GetUltraVoucher]")
		logrus.Info("[Failed from DB]-[Get Data Voucher-UV]")

		res = utils.GetMessageResponse(res, 404, false, errors.New("Voucher Not Found"))
		ctx.JSON(http.StatusOK, res)
		return
	}

	spanid := utilsgo.GetSpanId(span)
	sugarLogger.Info("REQUEST:", zap.String("SPANID", spanid), zap.String("CTRL", namectrl),
		zap.Any("BODY", req),
		zap.Any("HEADER", ctx.Request.Header))

	CallabckUVServices := v2_migrate.CallabckUVServices{
		General: models.GeneralModel{
			ParentSpan: span,
			OttoZaplog: sugarLogger,
			SpanId:     spanid,
			Context:    context,
		},
	}

	// Check voucher / get details voucher
	dtaVoucher, errVouc := db.GetVoucherDetails(getData.CampaignID)
	if errVouc != nil || dtaVoucher.RewardID == "" {
		sugarLogger.Info(fmt.Sprintf("Failed Get Voucher/Reward Details : ", errVouc))
		res = utils.GetMessageResponse(res, 404, false, errors.New("Voucher Not Found"))
		ctx.JSON(http.StatusOK, res)
		return
	}

	param := models.Params{
		AccountNumber: getData.Phone,
		// MerchantID:    dataToken.MerchantID,
		// InstitutionID: header.InstitutionID,
		SupplierID:  "UltraVoucher",
		ProductType: dtaVoucher.BrandName,
		NamaVoucher: dtaVoucher.VoucherName,
		CouponCode:  dtaVoucher.ExternalProductCode,
		Category:    dtaVoucher.BrandName,
		CouponID:    getData.CouponID,
		Point:       int(dtaVoucher.CostPoints),
		AccountId:   req.AccountId,
	}

	res = CallabckUVServices.CallbackVoucherUV(req, param, getData.CampaignID)

	// sugarLogger.Info("RESPONSE:", zap.String("SPANID", spanid), zap.String("CTRL", namectrl),
	// 	zap.Any("BODY", res))

	// datalog := utils.LogSpanMax(res)
	// zaplog.InfoWithSpan(span, namectrl,
	// 	zap.Any("RESP", datalog),
	// 	zap.Duration("backoff", time.Second))

	defer span.Finish()
	ctx.JSON(http.StatusOK, res)

}
