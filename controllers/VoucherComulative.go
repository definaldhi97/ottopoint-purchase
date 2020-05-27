package controllers

import (
	"errors"
	"fmt"
	"strings"

	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	opl "ottopoint-purchase/hosts/opl/host"
	modelsopl "ottopoint-purchase/hosts/opl/models"
	token "ottopoint-purchase/hosts/redis_token/host"
	"ottopoint-purchase/models"
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
)

func VoucherComulativeController(ctx *gin.Context) {
	req := models.VoucherComultaiveReq{}
	res := models.Response{}

	sugarLogger := ottologer.GetLogger()
	namectrl := "[VoucherComulativeController]"

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
	sugarLogger.Info("REQUEST : ", zap.String("SPANID", spanid), zap.String("CTRL", namectrl),
		zap.Any("BODY : ", req),
		zap.Any("HEADER : ", ctx.Request.Header))

	voucherComulative := services.VoucherComulativeService{
		General: models.GeneralModel{
			ParentSpan: span,
			OttoZaplog: sugarLogger,
			SpanId:     spanid,
			Context:    context,
		},
	}

	ultraVoucher := services.UseVoucherUltraVoucher{
		General: models.GeneralModel{
			ParentSpan: span,
			OttoZaplog: sugarLogger,
			SpanId:     spanid,
			Context:    context,
		},
	}

	cekVoucher, errVoucher := opl.VoucherDetail(req.CampaignID)
	if errVoucher != nil || cekVoucher.CampaignID == "" {
		sugarLogger.Info("[VoucherComulativeController]-[VoucherDetail]")
		sugarLogger.Info(fmt.Sprintf("Error : ", errVoucher))

		logs.Info("[VoucherComulativeController]-[VoucherDetail]")
		logs.Info(fmt.Sprintf("Error : ", errVoucher))

		res = utils.GetMessageResponse(res, 404, false, errors.New("Voucher Not Found"))
		ctx.JSON(http.StatusOK, res)
		return
	}

	dataUser, errUser := db.CheckUser(dataToken.Data)
	if errUser != nil || dataUser.CustID == "" {
		logs.Info("Internal Server Error : ", errUser)
		logs.Info("[VoucherComulativeController]-[CheckUser]")
		logs.Info("[Failed from DB]-[Get Data User]")

		// sugarLogger.Info("Internal Server Error : ", errredeem)
		sugarLogger.Info("[VoucherComulativeController]-[CheckUser]")
		sugarLogger.Info("[Failed from DB]-[Get Data User]")

		res = utils.GetMessageResponse(res, 500, false, errors.New("User belum Eligible"))
		ctx.JSON(http.StatusOK, res)
		return
	}

	data := SwitchCheckData(cekVoucher)

	logs.Info("SupplierID : ", data.SupplierID)
	logs.Info("producrType : ", data.ProductType)

	// sugarLogger.Info("SupplierID : ", data.SupplierID)
	// sugarLogger.Info("producrType : ", data.ProductType)

	if data.SupplierID == "OttoAG" {
		switch data.Category {
		case constants.CategoryPulsa:
			fmt.Println("Category Pulsa")
		case constants.CategoryFreeFire, constants.CategoryMobileLegend:
			fmt.Println("Category Game")
		case constants.CategoryPLN:
			fmt.Println("Category PLN")
		default:
			fmt.Println("Invalid Category")
			res = utils.GetMessageResponse(res, 500, false, errors.New("Invalid BrandName"))
			ctx.JSON(http.StatusOK, res)
			return
		}
	}

	param := models.Params{
		AccountNumber: dataToken.Data,
		MerchantID:    dataToken.MerchantID,
		InstitutionID: header.InstitutionID,
		AccountId:     dataUser.CustID,
		CampaignID:    req.CampaignID,
		SupplierID:    data.SupplierID,
		ProductType:   data.ProductType,
		ProductCode:   data.ProductCode,
		CouponCode:    data.CouponCode,
		NamaVoucher:   data.NamaVoucher,
		Point:         data.Point,
		Category:      data.Category,
	}

	switch data.SupplierID {
	case constants.UV:
		res = ultraVoucher.UltraVoucherServices(req, param)
	case constants.OttoAG:
		res = voucherComulative.VoucherComulative(req, param)
	}

	sugarLogger.Info("RESPONSE : ", zap.String("SPANID", spanid), zap.String("CTRL", namectrl),
		zap.Any("BODY : ", res))

	datalog := utils.LogSpanMax(res)
	zaplog.InfoWithSpan(span, namectrl,
		zap.Any("RESP : ", datalog),
		zap.Duration("backoff : ", time.Second))

	defer span.Finish()
	ctx.JSON(http.StatusOK, res)

}

func SwitchCheckData(data modelsopl.VoucherDetailResp) models.Params {
	res := models.Params{}

	coupon := data.Coupons[0]

	supplierid := coupon[:2]
	var supplierID string
	if supplierid == "UV" {
		supplierID = "Ultra Voucher"
		coupon = coupon[3:]
	} else {
		supplierID = "OttoAG"
	}

	var producrType string
	t := strings.ToLower(data.BrandName)
	switch t {
	case constants.CategoryPulsa:
		producrType = "Pulsa"
	case constants.CategoryFreeFire, constants.CategoryMobileLegend:
		producrType = "Game"
	case constants.CategoryPLN:
		producrType = "PLN"
	default:
		producrType = data.BrandName
	}

	res = models.Params{
		ProductType: producrType,
		ProductCode: coupon,
		CouponCode:  coupon,
		SupplierID:  supplierID,
		NamaVoucher: data.Name,
		Point:       data.CostInPoints,
		Category:    t,
		ExpDate:     data.CampaignActivity.ActiveTo,
	}

	return res
}
