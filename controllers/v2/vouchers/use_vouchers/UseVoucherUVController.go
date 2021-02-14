package use_vouchers

import (
	"errors"
	"fmt"

	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"

	"ottopoint-purchase/db"
	opl "ottopoint-purchase/hosts/opl/host"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"net/http"
)

func UseVouhcerUVController(ctx *gin.Context) {
	req := models.UseVoucherUVReq{}
	res := models.Response{}

	namectrl := "[PackageUserVoucherController]-[UseVouhcerUVController]"
	logReq := fmt.Sprintf("[AccountId : %v, VoucherCode : %v]", req.AccountId, req.VoucherCode)

	logrus.Info(namectrl)

	if err := ctx.ShouldBindJSON(&req); err != nil {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[ShouldBindJSON]-[Error : %v]", err))
		logrus.Println(logReq)

		res.Meta.Code = 03
		res.Meta.Message = "Gagal! Maaf transaksi Anda tidak dapat dilakukan saat ini. Silahkan dicoba lagi atau hubungi tim kami untuk informasi selengkapnya."
		ctx.JSON(http.StatusOK, res)

		return
	}

	getData, errData := db.GetUltraVoucher(req.VoucherCode, req.AccountId)
	if errData != nil || getData.CampaignID == "" {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[GetUltraVoucher]-[Error : %v]", errData))
		logrus.Println(logReq)

		res = utils.GetMessageResponse(res, 500, false, errors.New("User belum Eligible"))
		ctx.JSON(http.StatusOK, res)
		return
	}

	cekVoucher, errVoucher := opl.VoucherDetail(getData.CampaignID)
	if errVoucher != nil || cekVoucher.CampaignID == "" {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[VoucherDetail]-[Error : %v]", errVoucher))
		logrus.Println(logReq)

		res = utils.GetMessageResponse(res, 404, false, errors.New("Voucher Not Found"))
		ctx.JSON(http.StatusOK, res)
		return
	}

	// param := models.Params{
	// 	AccountNumber: getData.Phone,
	// 	// MerchantID:    dataToken.MerchantID,
	// 	// InstitutionID: header.InstitutionID,
	// 	SupplierID:  "UltraVoucher",
	// 	ProductType: cekVoucher.BrandName,
	// 	NamaVoucher: cekVoucher.Name,
	// 	CouponCode:  cekVoucher.Coupons[0],
	// 	Category:    cekVoucher.BrandName,
	// 	CouponID:    getData.CouponID,
	// 	Point:       cekVoucher.CostInPoints,
	// 	AccountId:   req.AccountId,
	// }

	// res = usevoucher.UseVoucherUV(req, param, getData.CampaignID)

	ctx.JSON(http.StatusOK, res)

}
