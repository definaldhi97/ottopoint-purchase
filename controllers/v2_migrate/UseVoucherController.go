package v2_migrate

import (
	"errors"
	"fmt"
	"net/http"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/controllers"
	"ottopoint-purchase/db"
	token "ottopoint-purchase/hosts/redis_token/host"
	"ottopoint-purchase/models"
	"ottopoint-purchase/services"
	"ottopoint-purchase/services/v2_migrate"
	"ottopoint-purchase/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	zaplog "github.com/opentracing-contrib/go-zap/log"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	ottologer "ottodigital.id/library/logger"
	utilsgo "ottodigital.id/library/utils"
)

type UseVouhcerMigrateController struct {
}

func (controller UseVouhcerMigrateController) UseVouhcerMigrateController(ctx *gin.Context) {

	fmt.Println("[ >>>>>>>>>>>>>>>>>> Use Vouhcer Migrate Controller <<<<<<<<<<<<<<<< ]")

	req := models.UseVoucherReq{}
	res := models.Response{}

	sugarLogger := ottologer.GetLogger()
	namectrl := "[UseVoucherController]"

	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.Meta.Code = 03
		res.Meta.Message = "Transaksi gagal, silahkan dicoba kembali. Jika masih gagal silahkan hubungi customer support kami."
		ctx.JSON(http.StatusOK, res)
		go sugarLogger.Error("Error, body Request", zap.Error(err))
		return
	}

	span := controllers.TracingFirstControllerCtx(ctx, req, namectrl)
	c := ctx.Request.Context()
	context := opentracing.ContextWithSpan(c, span)

	//validate request
	header, resultValidate := controllers.ValidateRequest(ctx, true, req, true)
	if !resultValidate.Meta.Status {
		ctx.JSON(http.StatusOK, resultValidate)
		return
	}

	dataToken, _ := token.CheckToken(header)

	spanid := utilsgo.GetSpanId(span)
	sugarLogger.Info("REQUEST:", zap.String("SPANID", spanid), zap.String("CTRL", namectrl),
		zap.Any("BODY", req),
		zap.Any("HEADER", ctx.Request.Header))

	// Check voucher / get details voucher
	cekVoucher, errVoucher := db.GetVoucherDetails(req.CampaignID)
	if errVoucher != nil || cekVoucher.RewardID == "" {
		sugarLogger.Info(fmt.Sprintf("Failed Get Voucher/Reward Details : ", errVoucher))
		res = utils.GetMessageResponse(res, 404, false, errors.New("Voucher Not Found"))
		ctx.JSON(http.StatusOK, res)
		return
	}

	dataVouch := SwitchDataVoucher(cekVoucher)

	var custIdOPL, merchant string
	if dataVouch.SupplierID == constants.CODE_VENDOR_UV {
		fmt.Println("[ Voucher Ultra Voucher ]")
		getData, errData := db.CheckCouponUV(dataToken.Data, req.CampaignID, req.CouponID)
		if errData != nil || getData.AccountId == "" {
			fmt.Println(fmt.Sprintf("Internal Server Error : %v\n", errData))
			res = utils.GetMessageResponse(res, 404, false, errors.New("Voucher Not Found"))
			ctx.JSON(http.StatusOK, res)
			return
		}
		custIdOPL = getData.AccountId
	} else {
		fmt.Println("[Voucher OttoAG]")
		dataUser, errUser := db.CheckUser(dataToken.Data)
		if errUser != nil || dataUser.CustID == "" {
			fmt.Println(fmt.Sprintf("Internal Server Error : %v\n", errUser))

			res = utils.GetMessageResponse(res, 500, false, errors.New("User belum Eligible"))
		}
		custIdOPL = dataUser.CustID
		merchant = dataUser.MerchantID
	}

	usevoucher := v2_migrate.UseVoucherMigrateServices{
		General: models.GeneralModel{
			ParentSpan: span,
			OttoZaplog: sugarLogger,
			SpanId:     spanid,
			Context:    context,
		},
	}

	param := models.Params{
		AccountNumber: dataToken.Data,
		MerchantID:    merchant,
		InstitutionID: header.InstitutionID,
		SupplierID:    dataVouch.SupplierID,
		AccountId:     custIdOPL,
		CampaignID:    req.CampaignID,
		ProductType:   dataVouch.ProductType,
		ProductCode:   dataVouch.ProductCode,
		NamaVoucher:   dataVouch.NamaVoucher,
		Category:      dataVouch.Category,
		CouponID:      dataVouch.CouponID,
		Point:         dataVouch.Point,
		ExpDate:       dataVouch.ExpDate,
	}

	switch dataVouch.SupplierID {
	case constants.CODE_VENDOR_UV:
		res = usevoucher.UseVoucherUV(req, param)
	case constants.CODE_VENDOR_AGREGATOR:
		res = usevoucher.UseVoucherAggregator(req, param)
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

func (controller UseVouhcerMigrateController) UseVoucherVidioController(ctx *gin.Context) {
	fmt.Println("[ >>>>>>>>>>>>>>>>>> Use Vouhcer Vidio Controller <<<<<<<<<<<<<<<< ]")

	var resp models.Response
	namectrl := "[ UseVoucherVidioController ]"
	sugarLogger := ottologer.GetLogger()

	couponId := ctx.Request.URL.Query().Get("couponId")

	// header
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

	//check header request
	if header.AppsID == "" || header.ChannelID == "" || header.InstitutionID == "" || header.DeviceID == "" || header.Geolocation == "" {
		resp = utils.GetMessageFailedErrorNew(resp, constants.RC_ERROR_HEADER_MANDATORY, constants.RD_ERROR_HEADER_MANDATORY)
		ctx.JSON(http.StatusOK, resp)
		return
	}

	span := controllers.TracingFirstControllerCtx(ctx, couponId, namectrl)
	c := ctx.Request.Context()
	context := opentracing.ContextWithSpan(c, span)

	// Validate Token user
	var RedisService = new(services.RedisService)
	auth := strings.ReplaceAll(header.Authorization, "Bearer ", "")
	keyRedis := header.InstitutionID + "-" + auth
	dataRedis := RedisService.GetData(keyRedis)

	if dataRedis.ResponseCode != "00" {
		resp = utils.GetMessageFailedErrorNew(resp, constants.RC_ERROR_INVALID_TOKEN, constants.RD_ERROR_INVALID_TOKEN)
		ctx.JSON(http.StatusOK, resp)
		return
	}

	spanid := utilsgo.GetSpanId(span)
	sugarLogger.Info("REQUEST:", zap.String("SPANID", spanid), zap.String("CTRL", namectrl),
		zap.Any("BODY", couponId),
		zap.Any("HEADER", ctx.Request.Header))

	UseVoucher := v2_migrate.UseVoucherMigrateServices{
		General: models.GeneralModel{
			ParentSpan: span,
			OttoZaplog: sugarLogger,
			SpanId:     spanid,
			Context:    context,
		},
	}

	resp = UseVoucher.UseVoucherVidio(couponId)

	sugarLogger.Info("RESPONSE:", zap.String("SPANID", spanid), zap.String("CTRL", namectrl),
		zap.Any("BODY", resp))

	datalog := utils.LogSpanMax(resp)
	zaplog.InfoWithSpan(span, namectrl,
		zap.Any("RESP", datalog),
		zap.Duration("backoff", time.Second))

	defer span.Finish()
	ctx.JSON(http.StatusOK, resp)

}
