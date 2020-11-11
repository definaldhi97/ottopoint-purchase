package controllers

import (
	"errors"
	"fmt"

	"ottopoint-purchase/constants"
	"ottopoint-purchase/models"
	"ottopoint-purchase/services"
	"ottopoint-purchase/utils"
	"time"

	"ottopoint-purchase/db"
	opl "ottopoint-purchase/hosts/opl/host"
	token "ottopoint-purchase/hosts/redis_token/host"

	"github.com/gin-gonic/gin"
	zaplog "github.com/opentracing-contrib/go-zap/log"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	ottologer "ottodigital.id/library/logger"
	utilsgo "ottodigital.id/library/utils"

	"net/http"
)

func UseVouhcerController(ctx *gin.Context) {
	req := models.UseVoucherReq{}
	res := models.Response{}

	sugarLogger := ottologer.GetLogger()
	namectrl := "[UseVoucherController]"

	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.Meta.Code = 03
		res.Meta.Message = "Transaksi gagal, silahkan dicoba kembali. Jika masih gagal silahkan hubungi customer support kami."
		// res.Meta.Message = "Gagal! Maaf transaksi Anda tidak dapat dilakukan saat ini. Silahkan dicoba lagi atau hubungi tim kami untuk informasi selengkapnya."
		ctx.JSON(http.StatusOK, res)
		go sugarLogger.Error("Error, body Request", zap.Error(err))
		return
	}

	span := TracingFirstControllerCtx(ctx, req, namectrl)
	c := ctx.Request.Context()
	context := opentracing.ContextWithSpan(c, span)

	//validate request
	header, resultValidate := ValidateRequest(ctx, true, req, true)
	if !resultValidate.Meta.Status {
		ctx.JSON(http.StatusOK, resultValidate)
		return
	}

	dataToken, _ := token.CheckToken(header)

	spanid := utilsgo.GetSpanId(span)
	sugarLogger.Info("REQUEST:", zap.String("SPANID", spanid), zap.String("CTRL", namectrl),
		zap.Any("BODY", req),
		zap.Any("HEADER", ctx.Request.Header))

	usevoucher := services.UseVoucherServices{
		General: models.GeneralModel{
			ParentSpan: span,
			OttoZaplog: sugarLogger,
			SpanId:     spanid,
			Context:    context,
		},
	}

	cekVoucher, errVoucher := opl.VoucherDetail(req.CampaignID)
	if errVoucher != nil || cekVoucher.CampaignID == "" {
		sugarLogger.Info("[UseVoucherController]-[VoucherDetail]")
		sugarLogger.Info(fmt.Sprintf("Error : ", errVoucher))

		fmt.Println("[UseVoucherController]-[VoucherDetail]")
		fmt.Println(fmt.Sprintf("Error : ", errVoucher))

		res = utils.GetMessageResponse(res, 404, false, errors.New("Voucher Not Found"))
		ctx.JSON(http.StatusOK, res)
		return
	}

	data := SwitchCheckData(cekVoucher)

	var custIdOPL, merchant string
	if data.SupplierID == "Ultra Voucher" {
		fmt.Println("[Voucher Ultra Voucher]")
		getData, errData := db.CheckCouponUV(dataToken.Data, req.CampaignID, req.CouponID)
		if errData != nil || getData.AccountId == "" {
			fmt.Println(fmt.Sprintf("Internal Server Error : %v\n", errData))
			sugarLogger.Info("[UseVoucherController]-[CheckCouponUV]")
			sugarLogger.Info("[Failed Failed from DB]-[Get Data Voucher-UV]")

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
			sugarLogger.Info("[UseVoucherController]-[CheckUser]")
			sugarLogger.Info("[Failed from DB]-[Get Data User]")

			res = utils.GetMessageResponse(res, 500, false, errors.New("User belum Eligible"))
		}
		custIdOPL = dataUser.CustID
		merchant = dataUser.MerchantID
	}

	fmt.Println("SupplierID : ", data.SupplierID)
	fmt.Println("producrType : ", data.ProductType)

	sugarLogger.Info("=== SupplierID ===")
	sugarLogger.Info(data.SupplierID)

	sugarLogger.Info("=== producrType ===")
	sugarLogger.Info(data.ProductType)

	param := models.Params{
		AccountNumber: dataToken.Data,
		MerchantID:    merchant,
		InstitutionID: header.InstitutionID,
		SupplierID:    data.SupplierID,
		AccountId:     custIdOPL,
		CampaignID:    req.CampaignID,
		ProductType:   data.ProductType,
		ProductCode:   data.ProductCode,
		NamaVoucher:   data.NamaVoucher,
		Category:      data.Category,
		CouponID:      req.CouponID,
		Point:         data.Point,
		ExpDate:       data.ExpDate,
	}

	switch data.SupplierID {
	case constants.UV:
		res = usevoucher.GetVoucherUV(req, param)
	case constants.OttoAG:
		res = usevoucher.UseVoucherOttoAG(req, param)
	case constants.VoucherAg:
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
