package earnings

import (
	"fmt"
	ctrl "ottopoint-purchase/controllers"
	earning "ottopoint-purchase/services/v2/earnings"

	"net/http"

	"github.com/gin-gonic/gin"

	"ottopoint-purchase/models"

	"github.com/sirupsen/logrus"
)

func CheckStatusEarningController(ctx *gin.Context) {

	req := models.CheckStatusEarningReq{}
	res := models.Response{}

	namectrl := "[PackageEarnings]-[CheckStatusEarningController]"

	// Debug, Println (Putih), Warn (Kuning), Error (Merah), Info (Biru)

	logrus.Info(namectrl)

	if err := ctx.ShouldBindJSON(&req); err != nil {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[ShouldBindJSON]-[Error : %v]", err))
		logrus.Println("ReferenceId : ", req.ReferenceId)

		res.Meta.Code = 03
		res.Meta.Message = "Gagal! Maaf transaksi Anda tidak dapat dilakukan saat ini. Silahkan dicoba lagi atau hubungi tim kami untuk informasi selengkapnya."
		ctx.JSON(http.StatusOK, res)
		return
	}

	// validate request
	header, resultValidate := ctrl.ValidateRequestWithoutAuth(ctx, req)
	if !resultValidate.Meta.Status {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[ValidateRequestWithoutAuth]-[Error : %v]", resultValidate))
		logrus.Println("ReferenceId : ", req.ReferenceId)

		ctx.JSON(http.StatusOK, resultValidate)
		return
	}

	logrus.Println("[Request]")
	logrus.Info("ReferenceId : ", req.ReferenceId)

	checkStatusEarning := new(earning.EarningsServices)

	res = checkStatusEarning.CheckStatusEarningServices(req.ReferenceId, header.InstitutionID)
	ctx.JSON(http.StatusOK, res)

}
