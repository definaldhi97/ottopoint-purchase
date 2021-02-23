package redeemtion

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/controllers"
	v "ottopoint-purchase/controllers/v2/vouchers"
	"ottopoint-purchase/db"
	core "ottopoint-purchase/hosts/op_corepoint"
	redishost "ottopoint-purchase/hosts/redis_token/host"
	"ottopoint-purchase/models"
	services "ottopoint-purchase/services/v2/vouchers/redeemtion"
	"ottopoint-purchase/utils"

	redeemtion "ottopoint-purchase/services/v2.1/vouchers/redeemtion"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Fields struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func PaymentSplitBillController(ctx *gin.Context) {

	req := models.PaymentSplitBillReq{}
	res := models.Response{}

	namectrl := "[PackageRedeemtionController]-[PaymentSplitBillController]"

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
	header, resultValidate := controllers.ValidateRequest(ctx, true, req, true)
	if !resultValidate.Meta.Status {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[ValidateRequest]-[Response : %v]", resultValidate))

		ctx.JSON(http.StatusOK, resultValidate)
		return
	}

	// get customer di redis
	dataToken, errToken := redishost.CheckToken(header)
	if errToken != nil {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[CheckToken]-[Error : %v]", errToken))

		res = utils.GetMessageFailedErrorNew(res, 500, "Internal Server Error")

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
	cekVoucher, errVoucher := db.GetVoucherDetails(req.CampaignId)
	if errVoucher != nil || cekVoucher.RewardID == "" {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[GetVoucherDetails]-[Error : %v]", errVoucher))

		res = utils.GetMessageResponse(res, 404, false, errors.New("Voucher Not Found"))
		ctx.JSON(http.StatusOK, res)
		return
	}

	// Balance
	balance, errBalance := core.GetBalancePoint(dataUser.CustID)
	if errBalance != nil || balance.ResponseCode != "00" {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[GetBalancePoint]-[Error : %v]", errBalance))

		res = utils.GetMessageFailedErrorNew(res, 500, "Internal Server Error")
		ctx.JSON(http.StatusOK, res)
		return
	}

	param := v.ParamRedeemtion(dataUser.CustID, cekVoucher)

	param.InstitutionID = header.InstitutionID
	param.CampaignID = req.CampaignId
	param.Email = dataUser.Email
	param.FirstName = dataUser.Nama
	param.LastName = dataUser.LastName
	param.AccountNumber = dataToken.Data
	param.MerchantID = dataUser.MerchantID

	balanceAmount := int64(cekVoucher.CostPoints) - int64(balance.Balance)
	balancePoint := balance.Balance

	var fields []Fields
	var custId, custId2 string

	f, _ := json.Marshal(req.FieldValue)
	errFields := json.Unmarshal(f, &fields)
	logrus.Error("Error Unmarshal errFields : ", errFields)

	for i := 0; i < len(fields); i++ {
		param.CustID = fields[i].Value
		custId = fields[i].Value

		if len(fields) > 1 {
			param.CustID = fields[0].Value + " || " + fields[1].Value

			custId = fields[0].Value
			custId2 = fields[1].Value
		}

	}

	reqOP := models.VoucherComultaiveReq{
		Jumlah:     1,
		CampaignID: param.CampaignID,
		CustID:     custId,
		CustID2:    custId2,
	}

	logrus.Println("[Request]")
	logrus.Info("CampaignId : ", req.CampaignId, "FieldValue : ", req.FieldValue, "PaymentMethod : ", req.PaymentMethod)

	switch req.PaymentMethod {
	case constants.SplitBillMethod:
		logrus.Println(" [ SplitBillMethod ]")
		res = services.PaymentSplitBillServices(req, param, int64(balancePoint), balanceAmount)
	case constants.FullPointMethod:
		logrus.Println(" [ FullPointMethod ]")
		switch param.SupplierID {

		case constants.CODE_VENDOR_OTTOAG:
			logrus.Println(" [ Product OTTOAG ]")
			res = redeemtion.RedeemtionOttoAG_V21_Service(reqOP, param, header)
		case constants.CODE_VENDOR_UV:
			logrus.Println(" [ Product Ultra Voucher ]")
			res = redeemtion.RedeemtionUV_V21_Service(reqOP, param, header)
		case constants.CODE_VENDOR_SEPULSA:
			logrus.Println(" [ Product Sepulsa ]")
			res = redeemtion.RedeemtionSepulsa_V21_Service(reqOP, param, header)
		case constants.CODE_VENDOR_AGREGATOR:
			logrus.Println(" [ Product Agregator ]")
			header.DeviceID = "H2H"
			res = redeemtion.RedeemtionAG_V21_Services(reqOP, param, header)
		}

	default:
		logrus.Error(namectrl)
		logrus.Error("Invalid Payment Method : ", req.PaymentMethod)
		res = utils.GetMessageFailedErrorNew(res, 500, "Internal Server Error")
	}

	ctx.JSON(http.StatusOK, res)

}
