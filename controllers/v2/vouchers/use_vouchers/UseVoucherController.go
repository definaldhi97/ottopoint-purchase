package use_vouchers

import (
	"errors"
	"fmt"
	"net/http"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/controllers"
	s "ottopoint-purchase/controllers/v2/vouchers"
	"ottopoint-purchase/db"
	token "ottopoint-purchase/hosts/redis_token/host"
	"ottopoint-purchase/models"
	vouchers "ottopoint-purchase/services/v2/vouchers/use_vouchers"
	"ottopoint-purchase/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// func UseVouhcerMigrateController(ctx *gin.Context) {
func UseVouchersControllers(ctx *gin.Context) {

	// fmt.Println("[ >>>>>>>>>>>>>>>>>> V2 Migrate Use Vouhcer Migrate Controller <<<<<<<<<<<<<<<< ]")

	req := models.UseVoucherReq{}
	res := models.Response{}

	namectrl := "[PackageUserVoucher]-[UseVouchersControllers]"
	logReq := fmt.Sprintf("[CouponID : %v]", req.CouponID)

	logrus.Info(namectrl)

	if err := ctx.ShouldBindJSON(&req); err != nil {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[ShouldBindJSON]-[Error : %v]", err))
		logrus.Println(logReq)

		res.Meta.Code = 03
		res.Meta.Message = "Transaksi gagal, silahkan dicoba kembali. Jika masih gagal silahkan hubungi customer support kami."
		ctx.JSON(http.StatusOK, res)
		return
	}

	//validate request
	header, resultValidate := controllers.ValidateRequest(ctx, true, req, true)
	if !resultValidate.Meta.Status {
		ctx.JSON(http.StatusOK, resultValidate)
		return
	}

	dataToken, _ := token.CheckToken(header)

	// Check voucher / get details voucher
	cekVoucher, errVoucher := db.GetVoucherDetails(req.CampaignID)
	if errVoucher != nil || cekVoucher.RewardID == "" {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[GetVoucherDetails]-[Error : %v]", errVoucher))
		logrus.Println(logReq)

		res = utils.GetMessageResponse(res, 404, false, errors.New("Voucher Not Found"))
		ctx.JSON(http.StatusOK, res)
		return
	}

	dataVouch := s.SwitchDataVoucher(cekVoucher)

	var custIdOPL, merchant string
	if dataVouch.SupplierID == constants.CODE_VENDOR_UV {
		fmt.Println("[ Voucher Ultra Voucher ]")
		getData, errData := db.CheckCouponUV(dataToken.Data, req.CampaignID, req.CouponID)
		if errData != nil || getData.AccountId == "" {

			logrus.Error(namectrl)
			logrus.Error(fmt.Sprintf("[CheckCouponUV]-[Error : %v]", errData))
			logrus.Println(logReq)

			res = utils.GetMessageResponse(res, 404, false, errors.New("Voucher Not Found"))
			ctx.JSON(http.StatusOK, res)
			return
		}
		custIdOPL = getData.AccountId
	} else {

		fmt.Println("[Voucher OttoAG]")
		dataUser, errUser := db.CheckUser(dataToken.Data)
		if errUser != nil || dataUser.CustID == "" {

			logrus.Error(namectrl)
			logrus.Error(fmt.Sprintf("[CheckUser]-[Error : %v]", errUser))
			logrus.Println(logReq)

			res = utils.GetMessageResponse(res, 500, false, errors.New("User belum Eligible"))
		}
		custIdOPL = dataUser.CustID
		merchant = dataUser.MerchantID
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

	logrus.Println("[Request]")
	logrus.Info("CampaignID : ", req.CampaignID, "CouponID : ", req.CouponID, "CustID : ", req.CustID, "CustID2 : ", req.CustID2)

	switch dataVouch.SupplierID {
	case constants.CODE_VENDOR_UV:
		res = vouchers.UseVoucherUVServices(req, param)
	case constants.CODE_VENDOR_AGREGATOR:
		res = vouchers.UseVoucherAggregatorService(req, param)
	}

	ctx.JSON(http.StatusOK, res)
	return

}
