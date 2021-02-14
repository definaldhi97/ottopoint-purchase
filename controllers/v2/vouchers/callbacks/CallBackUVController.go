package callbacks

import (
	"errors"
	"fmt"

	"ottopoint-purchase/models"
	service "ottopoint-purchase/services/v2/vouchers/callbacks"
	"ottopoint-purchase/utils"

	"ottopoint-purchase/db"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CallBackUVController(ctx *gin.Context) {

	req := models.UseVoucherUVReq{}
	res := models.Response{}

	namectrl := "[PackageCallBacksController]-[CallBackUVController]"

	logReq := fmt.Sprintf("[AccountId : %v, VoucherCode : %v]", req.AccountId, req.VoucherCode)

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

	getData, errData := db.GetUltraVoucher(req.VoucherCode, req.AccountId)
	if errData != nil || getData.CampaignID == "" {

		logrus.Error(namectrl)
		logrus.Error("[GetUltraVoucher]-[Error : %v]", errData)
		logrus.Println(logReq)

		res = utils.GetMessageResponse(res, 404, false, errors.New("Voucher Not Found"))
		ctx.JSON(http.StatusOK, res)
		return
	}

	// Check voucher / get details voucher
	dtaVoucher, errVouc := db.GetVoucherDetails(getData.CampaignID)
	if errVouc != nil || dtaVoucher.RewardID == "" {

		logrus.Error(namectrl)
		logrus.Error("[GetVoucherDetails]-[Error : %v]", errVouc)
		logrus.Println(logReq)

		res = utils.GetMessageResponse(res, 404, false, errors.New("Voucher Not Found"))
		ctx.JSON(http.StatusOK, res)
		return
	}

	param := models.Params{
		AccountNumber: getData.Phone,
		// MerchantID:    dataToken.MerchantID,
		// InstitutionID: header.InstitutionID,
		SupplierID:  "UltraVoucher",
		ProductType: dtaVoucher.BrandName,
		NamaVoucher: dtaVoucher.VoucherName,
		CouponCode:  dtaVoucher.ExternalProductCode,
		Category:    dtaVoucher.BrandName,
		CouponID:    getData.CouponID,
		Point:       int(dtaVoucher.CostPoints),
		AccountId:   req.AccountId,
	}

	logrus.Println("[Request]")
	logrus.Info("AccountId : ", req.AccountId, "VoucherCode : ", req.VoucherCode)

	res = service.CallbackVoucherUV(req, param, getData.CampaignID)

	ctx.JSON(http.StatusOK, res)

	return

}
