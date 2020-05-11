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

	"github.com/astaxie/beego/logs"

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
	namectrl := "[UseVouhcerController]"

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

	dataUser, errUser := db.CheckUser(dataToken.Data)
	if errUser != nil || dataUser.CustID == "" {
		logs.Info("Internal Server Error : ", errUser)
		logs.Info("[UltraVoucherServices]-[CheckUser]")
		logs.Info("[Failed Redeem Voucher]-[Get Data User]")

		// sugarLogger.Info("Internal Server Error : ", errredeem)
		sugarLogger.Info("[UltraVoucherServices]-[CheckUser]")
		sugarLogger.Info("[Failed Redeem Voucher]-[Get Data User]")

		res = utils.GetMessageResponse(res, 500, false, errors.New("User belum Eligible"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

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
		sugarLogger.Info("[HistoryVoucherCustomer]-[VoucherComulative-Controller]")
		sugarLogger.Info(fmt.Sprintf("Error : ", errVoucher))

		logs.Info("[HistoryVoucherCustomer]-[VoucherComulative-Controller]")
		logs.Info(fmt.Sprintf("Error : ", errVoucher))

		res = utils.GetMessageResponse(res, 422, false, errors.New("Internal Server Error"))
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	data := SwitchCheckData(cekVoucher)

	logs.Info("SupplierID : ", data.SupplierID)
	logs.Info("producrType : ", data.ProductType)

	sugarLogger.Info("=== SupplierID ===")
	sugarLogger.Info(data.SupplierID)

	sugarLogger.Info("=== producrType ===")
	sugarLogger.Info(data.ProductType)

	param := models.Params{
		AccountNumber: dataToken.Data,
		MerchantID:    dataToken.MerchantID,
		InstitutionID: header.InstitutionID,
		SupplierID:    data.SupplierID,
		CustID:        dataUser.CustID,
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
