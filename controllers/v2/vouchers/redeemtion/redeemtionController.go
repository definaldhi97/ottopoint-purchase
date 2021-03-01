package redeemtion

import (
	"errors"
	"fmt"
	"net/http"
	"ottopoint-purchase/constants"
	validate "ottopoint-purchase/controllers"
	c "ottopoint-purchase/controllers/v2/vouchers"
	"ottopoint-purchase/db"
	redishost "ottopoint-purchase/hosts/redis_token/host"
	"ottopoint-purchase/models"
	redeemtion "ottopoint-purchase/services/v2/vouchers/redeemtion"
	"ottopoint-purchase/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// func V2_RedeemtionVoucherController(ctx *gin.Context) {
func RedeemtionController(ctx *gin.Context) {

	req := models.VoucherComultaiveReq{}
	res := models.Response{}

	namectrl := "[PackageRedeemtion]-[RedeemtionController]"

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
	header, resultValidate := validate.ValidateRequest(ctx, true, req, true)
	if !resultValidate.Meta.Status {
		ctx.JSON(http.StatusOK, resultValidate)
		return
	}

	// get customer di redis
	dataToken, errToken := redishost.CheckToken(header)
	if errToken != nil {
		fmt.Println("Failed Get Token .. ..")
		res = utils.GetMessageFailedErrorNew(res, 500, "Internal Server Error")

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[CheckToken]-[Error : %v]", errToken))

		ctx.JSON(http.StatusOK, res)
		return
	}

	// check user
	dataUser, errUser := db.UserWithInstitution(dataToken.Data, header.InstitutionID)
	if errUser != nil || dataUser.CustID == "" {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[UserWithInstitution]-[Error : %v]", errUser))

		res = utils.GetMessageResponse(res, 404, false, errors.New("User belum Eligible"))
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

	param := c.ParamRedeemtion(dataUser.CustID, cekVoucher)
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

	param.InstitutionID = header.InstitutionID
	param.CampaignID = req.CampaignID
	param.AccountNumber = dataToken.Data
	param.MerchantID = dataUser.MerchantID
	param.Fields = cekVoucher.Fields

	logrus.Println("[Request]")
	logrus.Info("CampaignId : ", req.CampaignID, "CustID : ", req.CustID, "CustID2 : ", req.CustID2, "Jumlah : ", req.Jumlah)

	switch param.SupplierID {
	case constants.CODE_VENDOR_OTTOAG:
		fmt.Println(" [ Product OTTOAG ]")
		res = redeemtion.RedeemtionOttoAGServices(req, param)
	case constants.CODE_VENDOR_UV:
		fmt.Println(" [ Product Ultra Voucher ]")
		res = redeemtion.RedeemtionUVServices(req, param)
	case constants.CODE_VENDOR_SEPULSA:
		fmt.Println(" [ Product Sepulsa ]")
		res = redeemtion.RedeemtionSepulsaServices(req, param)
	case constants.CODE_VENDOR_AGREGATOR:
		fmt.Println(" [ Product Agregator ]")
		header.DeviceID = "H2H"
		res = redeemtion.RedeemtionAggServices(req, param, header)
	}

	ctx.JSON(http.StatusOK, res)

	return

}
