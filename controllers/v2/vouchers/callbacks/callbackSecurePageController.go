package callbacks

import (
	"fmt"

	"ottopoint-purchase/constants"
	validate "ottopoint-purchase/controllers"
	c "ottopoint-purchase/controllers/v2/vouchers"
	"ottopoint-purchase/db"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"

	// service "ottopoint-purchase/services/v2/vouchers/callbacks"
	redeemtion "ottopoint-purchase/services/v2.1/vouchers/redeemtion"

	"net/http"

	"errors"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CallBackSecurePageController(ctx *gin.Context) {

	req := models.CallBackSGReq{}
	res := models.Response{}

	namectrl := "[PackageCallBacks_V2_Controller]-[CallBackSecurePageController]"

	logReq := fmt.Sprintf("[OttoRefNo : %v, IssuerRefNo : %v]", req.OttoRefNo, req.IssuerRefNo)

	logrus.Info(namectrl)

	if err := ctx.ShouldBindJSON(&req); err != nil {

		logrus.Error(namectrl)
		logrus.Error("[ShouldBindJSON]-[Error : %v]", err)
		logrus.Println(logReq)

		res.Meta.Code = 03
		res.Meta.Message = "Gagal! Maaf transaksi Anda tidak dapat dilakukan saat ini. Silahkan dicoba lagi atau hubungi tim kami untuk informasi selengkapnya."

		ctx.JSON(http.StatusOK, res)
		return
	}

	// validate request
	header, resultValidate := validate.ValidateRequest(ctx, true, req, true)
	if !resultValidate.Meta.Status {
		ctx.JSON(http.StatusOK, resultValidate)
		return
	}

	getData, errData := db.GetDataSplitBillbyTrxID(req.TrxRef)
	if errData != nil || len(getData) < 2 {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[GetDataSplitBillbyTrxID]-[Error : %v]", errData))

		res = utils.GetMessageResponse(res, 500, false, errors.New("Internal Server Error"))
		ctx.JSON(http.StatusOK, res)
		return
	}

	var balancePoint int64
	var accountID, accountNumber, custID, campaignID string

	cekVoucher := models.VoucherDetailsManagement{}

	for _, val := range getData {

		balancePoint = 0
		accountID = val.AccountId
		accountNumber = val.AccountNumber
		campaignID = val.CampaignId
		custID = val.CustID

		if val.ExternalReffId == req.TrxRef {
			balancePoint = val.Value

			categoryId := []string{val.ProductCategoryId}

			cekVoucher = models.VoucherDetailsManagement{
				RewardID:    val.MRewardID,
				VoucherName: val.Voucher,
				CostPoints:  float64(val.Point),
				// UsageLimit         : ,
				BrandName: val.ProductType,
				// ActivityActiveFrom : ,
				// ActivityActiveTo   : ,
				CategoriesID: categoryId,
				CodeSuplier:  val.SupplierID,
				// RewardCodes        : ,
				ExternalProductCode: val.ProductCode,
				// InternalProductCode: ,
				ProductID: val.MProductID,
			}

		}
	}

	if balancePoint == 0 {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[GetDataSplitBillbyTrxID]-[TrxRef not Found]"))

		res = utils.GetMessageResponse(res, 500, false, errors.New("Internal Server Error"))
		ctx.JSON(http.StatusOK, res)
		return
	}

	param := c.ParamRedeemtion(accountID, custID, cekVoucher)

	if param.ResponseCode != 200 {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[ParamRedeemtion]-[Response : %v]", param))

		res = utils.GetMessageResponse(res, 404, false, errors.New("Invalid BrandName / Prefix"))

		ctx.JSON(http.StatusOK, res)
		return

	}

	reqOP := models.VoucherComultaiveReq{
		Jumlah:     1,
		CampaignID: campaignID,
		CustID:     utils.Before(custID, "|| "),
		CustID2:    utils.After(custID, " ||"),
	}

	if param.Category == constants.CategoryVidio {
		reqOP.CustID = "0"
	}

	param.CampaignID = campaignID
	param.AccountId = accountID
	param.AccountNumber = accountNumber

	logrus.Println("[Request]")
	logrus.Info("Amount : ", req.Amount, "Issuer : ", req.Issuer, "IssuerRefNo : ", req.IssuerRefNo, "OttoRefNo : ", req.OttoRefNo,
		"ResponseCode : ", req.ResponseCode, "ResponseDescription : ", req.ResponseDescription, "TrxRef : ", req.TrxRef, "TransactionType : ", req.TransactionType, "UserId : ", req.UserId)

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

	ctx.JSON(http.StatusOK, res)

	return

}
