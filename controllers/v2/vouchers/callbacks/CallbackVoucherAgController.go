package callbacks

import (
	"fmt"
	"net/http"
	"ottopoint-purchase/models"
	service "ottopoint-purchase/services/v2/vouchers/callbacks"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CallbackVoucherAggController(ctx *gin.Context) {
	// func (controller *CallbBackControllers) CallbackVoucherAggController(ctx *gin.Context) {

	// fmt.Println("[ >>>>>>>>>>>>>>>>>>>>> V2 Migrate Callbakc Voucher Agg Controller <<<<<<<<<<<<<<<<<< ]")

	var (
		req models.CallbackRequestVoucherAg
		res models.Response
	)

	namectrl := "[PackageCallBacksController]-[CallbackVoucherAggController]"

	logReq := fmt.Sprintf("[TransactionID : %v]", req.TransactionID)

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

	res = service.CallbackVoucherAgg(req)

	ctx.JSON(http.StatusOK, res)
}
