package Redeemtion

import (
	"errors"
	"fmt"
	"net/http"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/controllers"
	v2_redeemtion "ottopoint-purchase/controllers/v2/Redeemtion"
	"ottopoint-purchase/db"
	redishost "ottopoint-purchase/hosts/redis_token/host"
	"ottopoint-purchase/models"
	v21_redeemtion "ottopoint-purchase/services/v2.1/Redeemtion"
	"ottopoint-purchase/utils"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	ottologer "ottodigital.id/library/logger"
	utilsgo "ottodigital.id/library/utils"
)

type V21_RedeemtionVoucherController struct{}

func (controller *V21_RedeemtionVoucherController) V21_RedeemtionVoucherController(ctx *gin.Context) {
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
	header, resultValidate := controllers.ValidateRequest(ctx, true, req, false)
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

	dataVouch := v2_redeemtion.SwitchDataVoucher(cekVoucher)

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

	VoucherOttoAgMigrateService := v21_redeemtion.V21_VoucherOttoAgService{
		General: models.GeneralModel{
			ParentSpan: span,
			OttoZaplog: sugarLogger,
			SpanId:     spanid,
			Context:    context,
		},
	}

	VoucherUVMigrateService := v21_redeemtion.V21_VoucherUVService{
		General: models.GeneralModel{
			ParentSpan: span,
			OttoZaplog: sugarLogger,
			SpanId:     spanid,
			Context:    context,
		},
	}

	VoucherSepulsaMigrateService := v21_redeemtion.V21_VoucherSepulsaService{
		General: models.GeneralModel{
			ParentSpan: span,
			OttoZaplog: sugarLogger,
			SpanId:     spanid,
			Context:    context,
		},
	}

	VoucherAgMigrateServices := v21_redeemtion.V21_VoucherAgServices{
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
		res = VoucherOttoAgMigrateService.V21_VoucherOttoAg(req, param, header)
	case constants.CODE_VENDOR_UV:
		fmt.Println(" [ Product Ultra Voucher ]")
		res = VoucherUVMigrateService.V21_VoucherUV(req, param, header)
	case constants.CODE_VENDOR_SEPULSA:
		fmt.Println(" [ Product Sepulsa ]")
		res = VoucherSepulsaMigrateService.V21_VoucherSepulsa(req, param, header)
	case constants.CODE_VENDOR_AGREGATOR:
		fmt.Println(" [ Product Agregator ]")
		header.DeviceID = "H2H"
		res = VoucherAgMigrateServices.V21_VoucherAg(req, param, header)
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
