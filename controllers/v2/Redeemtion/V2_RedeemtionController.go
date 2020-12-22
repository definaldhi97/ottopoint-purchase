package Redeemtion

import (
	"errors"
	"fmt"
	"net/http"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/controllers"
	"ottopoint-purchase/db"
	redishost "ottopoint-purchase/hosts/redis_token/host"
	"ottopoint-purchase/models"
	"ottopoint-purchase/services/v2/Redeemtion"
	"ottopoint-purchase/utils"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	ottologer "ottodigital.id/library/logger"
	utilsgo "ottodigital.id/library/utils"
)

type V2_RedeemtionVoucherController struct{}

func (controller *V2_RedeemtionVoucherController) V2_RedeemtionVoucherController(ctx *gin.Context) {
	fmt.Println("[ >>>>>>>>>>>>>>>>>> V2 Migrate Redeemtion Voucher Controller <<<<<<<<<<<<<<<< ]")

	req := models.VoucherComultaiveReq{}
	res := models.Response{}

	sugarLogger := ottologer.GetLogger()
	namectrl := "[RedeemtionVoucherController]"

	if err := ctx.ShouldBindJSON(&req); err != nil {
		message := "Gagal! Maaf transaksi Anda tidak dapat dilakukan saat ini. Silahkan dicoba lagi atau hubungi tim kami untuk informasi selengkapnya."
		res = utils.GetMessageFailedErrorNew(res, 03, message)
		go sugarLogger.Error("Error, body Request", zap.Error(err))
		ctx.JSON(http.StatusOK, res)
		return
	}

	span := controllers.TracingFirstControllerCtx(ctx, req, namectrl)
	c := ctx.Request.Context()
	context := opentracing.ContextWithSpan(c, span)

	// validate request
	header, resultValidate := controllers.ValidateRequest(ctx, true, req, true)
	if !resultValidate.Meta.Status {
		ctx.JSON(http.StatusOK, resultValidate)
		return
	}

	spanid := utilsgo.GetSpanId(span)
	sugarLogger.Info("REQUEST : ", zap.String("SPANID", spanid), zap.String("CTRL", namectrl),
		zap.Any("BODY : ", req),
		zap.Any("HEADER : ", ctx.Request.Header))

	// get customer di redis
	dataToken, errToken := redishost.CheckToken(header)
	if errToken != nil {
		fmt.Println("Failed Get Token .. ..")
		res = utils.GetMessageFailedErrorNew(res, 500, "Internal Server Error")
		go sugarLogger.Error("Internal Server Error Get token Customer to redis", zap.Error(errToken))
		ctx.JSON(http.StatusOK, res)
		return
	}

	// check user
	dataUser, errUser := db.UserWithInstitution(dataToken.Data, header.InstitutionID)
	if errUser != nil || dataUser.CustID == "" {
		logs.Info("Internal Server Error : ", errUser)
		sugarLogger.Info("Customer not found")
		res = utils.GetMessageResponse(res, 404, false, errors.New("User belum Eligible"))
		ctx.JSON(http.StatusOK, res)
		return
	}

	// Check voucher / get details voucher
	cekVoucher, errVoucher := db.GetVoucherDetails(req.CampaignID)
	if errVoucher != nil || cekVoucher.RewardID == "" {
		sugarLogger.Info(fmt.Sprintf("Failed Get Voucher/Reward Details : ", errVoucher))
		res = utils.GetMessageResponse(res, 404, false, errors.New("Voucher Not Found"))
		ctx.JSON(http.StatusOK, res)
		return
	}

	dataVouch := SwitchDataVoucher(cekVoucher)

	if dataVouch.Category == constants.CategoryVidio {
		req.CustID = "0"
	}

	if dataVouch.SupplierID == constants.CODE_VENDOR_OTTOAG {
		validateVerfix := controllers.ValidatePerfix(req.CustID, dataVouch.ProductCode, dataVouch.Category)
		if validateVerfix == false {
			fmt.Println("Invalid verfix")
			res = utils.GetMessageResponse(res, 500, false, errors.New("Nomor akun ini tidak terdafatr"))
			ctx.JSON(http.StatusOK, res)
			return
		}
		if dataVouch.Category == "" {
			fmt.Println("Invalid Category")
			res = utils.GetMessageResponse(res, 500, false, errors.New("Invalid BrandName"))
			ctx.JSON(http.StatusOK, res)
			return
		}
	}

	VoucherOttoAgMigrateService := Redeemtion.V2_VoucherOttoAgMigrateService{
		General: models.GeneralModel{
			ParentSpan: span,
			OttoZaplog: sugarLogger,
			SpanId:     spanid,
			Context:    context,
		},
	}

	VoucherUVMigrateService := Redeemtion.V2_VoucherUVService{
		General: models.GeneralModel{
			ParentSpan: span,
			OttoZaplog: sugarLogger,
			SpanId:     spanid,
			Context:    context,
		},
	}

	VoucherSepulsaMigrateService := Redeemtion.V2_VoucherSepulsaService{
		General: models.GeneralModel{
			ParentSpan: span,
			OttoZaplog: sugarLogger,
			SpanId:     spanid,
			Context:    context,
		},
	}

	VoucherAgMigrateServices := Redeemtion.V2_VoucherAgServices{
		General: models.GeneralModel{
			ParentSpan: span,
			OttoZaplog: sugarLogger,
			SpanId:     spanid,
			Context:    context,
		},
	}

	param := models.Params{
		AccountNumber:       dataToken.Data,
		MerchantID:          dataUser.MerchantID,
		InstitutionID:       header.InstitutionID,
		AccountId:           dataUser.CustID,
		CampaignID:          req.CampaignID,
		SupplierID:          dataVouch.SupplierID,
		ProductType:         dataVouch.ProductType,
		ProductCode:         dataVouch.ProductCode,
		CouponCode:          dataVouch.CouponCode,
		NamaVoucher:         dataVouch.NamaVoucher,
		Point:               dataVouch.Point,
		Category:            dataVouch.Category,
		UsageLimitVoucher:   dataVouch.UsageLimitVoucher,
		ProductCodeInternal: dataVouch.ProductCodeInternal,
		ProductID:           cekVoucher.ProductID,
		CategoryID:          dataVouch.CategoryID,
		RewardID:            dataVouch.RewardID,
		ExpDate:             dataVouch.ExpDate,
	}

	switch dataVouch.SupplierID {
	case constants.CODE_VENDOR_OTTOAG:
		fmt.Println(" [ Product OTTOAG ]")
		res = VoucherOttoAgMigrateService.VoucherOttoAg(req, param)
	case constants.CODE_VENDOR_UV:
		fmt.Println(" [ Product Ultra Voucher ]")
		res = VoucherUVMigrateService.VoucherUV(req, param)
	case constants.CODE_VENDOR_SEPULSA:
		fmt.Println(" [ Product Sepulsa ]")
		res = VoucherSepulsaMigrateService.VoucherSepulsa(req, param)
	case constants.CODE_VENDOR_AGREGATOR:
		fmt.Println(" [ Product Agregator ]")
		header.DeviceID = "H2H"
		res = VoucherAgMigrateServices.VoucherAg(req, param, header)
	}

	// sugarLogger.Info("RESPONSE : ", zap.String("SPANID", spanid), zap.String("CTRL", namectrl),
	// 	zap.Any("BODY : ", res))

	// datalog := utils.LogSpanMax(res)
	// zaplog.InfoWithSpan(span, namectrl,
	// 	zap.Any("RESP : ", datalog),
	// 	zap.Duration("backoff : ", time.Second))

	defer span.Finish()
	ctx.JSON(http.StatusOK, res)

}

func SwitchDataVoucher(data models.VoucherDetailsManagement) models.Params {
	var result models.Params

	var producrType string
	t := strings.ToLower(data.BrandName)
	switch t {
	case constants.CategoryFreeFire, constants.CategoryMobileLegend:
		producrType = "Game"
	default:
		producrType = data.BrandName
	}

	var categoriesID *string
	if len(data.CategoriesID) == 0 {
		fmt.Println("Not CategoriesID")
		categoriesID = nil
	} else {
		categoriesID = &data.CategoriesID[0]
	}
	result = models.Params{
		ProductType:         producrType,
		ProductCode:         data.ExternalProductCode,
		CouponCode:          data.ExternalProductCode,
		SupplierID:          data.CodeSuplier,
		NamaVoucher:         data.VoucherName,
		Point:               int(data.CostPoints),
		Category:            strings.ToLower(producrType),
		ExpDate:             data.ActivityActiveTo,
		CategoryID:          categoriesID,
		UsageLimitVoucher:   data.UsageLimit,
		ProductCodeInternal: data.InternalProductCode,
		RewardID:            data.RewardID,
		ProductID:           data.ProductID,
	}

	return result
}