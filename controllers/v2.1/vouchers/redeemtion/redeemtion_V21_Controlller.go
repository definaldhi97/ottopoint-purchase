package redeemtion

import (
	"errors"
	"fmt"
	"net/http"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/controllers"
	c "ottopoint-purchase/controllers/v2/vouchers"
	"ottopoint-purchase/db"
	redishost "ottopoint-purchase/hosts/redis_token/host"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"

	redeemtion "ottopoint-purchase/services/v2.1/vouchers/redeemtion"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// func (controller *V21_RedeemtionVoucherController) V21_RedeemtionVoucherController(ctx *gin.Context) {
func RedeemtionControllerV21(ctx *gin.Context) {
	// logrus.Println("[ >>>>>>>>>>>>>>>>>> V2 Migrate Redeemtion Voucher Controller <<<<<<<<<<<<<<<< ]")

	req := models.VoucherComultaiveReq{}
	res := models.Response{}

	namectrl := "[PackageRedeemtion]-[RedeemtionController_V2.1]"

	logrus.Info(namectrl)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		message := "Gagal! Maaf transaksi Anda tidak dapat dilakukan saat ini. Silahkan dicoba lagi atau hubungi tim kami untuk informasi selengkapnya."
		res = utils.GetMessageFailedErrorNew(res, 03, message)

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[ShouldBindJSON]-[Error : %v]", err))

		ctx.JSON(http.StatusOK, res)
		return
	}

	// validate request
	header, resultValidate := controllers.ValidateRequest(ctx, true, req, false)
	if !resultValidate.Meta.Status {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[ValidateRequest]-[Error : %v]", resultValidate))

		ctx.JSON(http.StatusOK, resultValidate)
		return
	}

	// get customer di redis
	dataToken, errToken := redishost.CheckToken(header)
	if errToken != nil {
		logrus.Println("Failed Get Token .. ..")

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[CheckToken]-[Error : %v]", errToken))

		res = utils.GetMessageFailedErrorNew(res, 500, "Internal Server Error")
		ctx.JSON(http.StatusOK, res)
		return
	}

	// Check voucher / get details voucher
	cekVoucher, errVoucher := db.GetVoucherDetails(req.CampaignID)
	if errVoucher != nil || cekVoucher.RewardID == "" {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[GetVoucherDetails]-[Error : %v]", errVoucher))

		res = utils.GetMessageResponse(res, 404, false, errors.New("Voucher Not Found"))
		ctx.JSON(http.StatusOK, res)
		return
	}

	// check user
	dataUser, errUser := db.UserWithInstitution(dataToken.Data, header.InstitutionID)
	if errUser != nil || dataUser.CustID == "" {
		logrus.Info("Internal Server Error : ", errUser)

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[CheckToken]-[Error : %v]", errUser))

		res = utils.GetMessageResponse(res, 404, false, errors.New("User belum Eligible"))
		ctx.JSON(http.StatusOK, res)
		return
	}

	param := c.ParamRedeemtion(req.CustID, cekVoucher)
	if param.ResponseCode != 200 {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[ParamRedeemtion]-[Response : %v]", param))

		res = utils.GetMessageResponse(res, 404, false, errors.New("Invalid BrandName / Prefix"))

		ctx.JSON(http.StatusOK, res)
		return

	}

	if param.Category == constants.CategoryVidio {
		req.CustID = "0"
	}

	param.AccountId = dataUser.CustID
	param.InstitutionID = header.InstitutionID
	param.CampaignID = req.CampaignID
	param.AccountNumber = dataToken.Data
	param.MerchantID = dataUser.MerchantID

	logrus.Println("[Request]")
	logrus.Info("CampaignId : ", req.CampaignID, "CustID : ", req.CustID, "CustID2 : ", req.CustID2, "Jumlah : ", req.Jumlah)

	switch param.SupplierID {
	case constants.CODE_VENDOR_DUMY:
		logrus.Println(" [ Product Dummy ]")
		res = redeemtion.RedeemtionDummyService(req, param, header)
	case constants.CODE_VENDOR_OTTOAG:
		logrus.Println(" [ Product OTTOAG ]")
		res = redeemtion.RedeemtionOttoAG_V21_Service(req, param, header)
	case constants.CODE_VENDOR_UV:
		logrus.Println(" [ Product Ultra Voucher ]")
		res = redeemtion.RedeemtionUV_V21_Service(req, param, header)
	case constants.CODE_VENDOR_SEPULSA:
		logrus.Println(" [ Product Sepulsa ]")
		res = redeemtion.RedeemtionSepulsa_V21_Service(req, param, header)
	case constants.CODE_VENDOR_AGREGATOR:
		logrus.Println(" [ Product Agregator ]")
		header.DeviceID = "H2H"
		res = redeemtion.RedeemtionAG_V21_Services(req, param, header)
	}

	ctx.JSON(http.StatusOK, res)
	return

}