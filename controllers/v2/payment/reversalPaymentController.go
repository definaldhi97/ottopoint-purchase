package payment

import (
	"fmt"
	"os"
	ctrl "ottopoint-purchase/controllers"
	sp "ottopoint-purchase/models/v2/payment"

	"net/http"
	"ottopoint-purchase/models"
	"ottopoint-purchase/services/v2/payment"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ReversalPaymentController(ctx *gin.Context) {

	req := sp.ReversalPaymentReq{}
	res := models.Response{}
	param := models.Params{}

	namectrl := "[PackagePayment]-[ReversalPaymentController]"
	logReq := fmt.Sprintf("[ReferenceId : %v]", req.ReferenceId)

	logrus.Info(namectrl)

	// logrus
	var log = logrus.New()
	log.Formatter = new(logrus.JSONFormatter)
	log.Formatter = new(logrus.TextFormatter)                     //default
	log.Formatter.(*logrus.TextFormatter).DisableColors = true    // remove colors
	log.Formatter.(*logrus.TextFormatter).DisableTimestamp = true // remove timestamp from test output
	log.Level = logrus.TraceLevel
	log.Out = os.Stdout

	if err := ctx.ShouldBindJSON(&req); err != nil {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[ShouldBindJSON]-[Error : %v]", err))
		logrus.Println(logReq)

		res.Meta.Code = 03
		res.Meta.Message = "Gagal! Maaf transaksi Anda tidak dapat dilakukan saat ini. Silahkan dicoba lagi atau hubungi tim kami untuk informasi selengkapnya."
		ctx.JSON(http.StatusOK, res)
		return
	}

	if req.ReferenceId == "" {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[Invalid Mandatory]-[Point/Cash/ReferenceId]"))
		logrus.Println(logReq)

		res.Meta.Code = 196
		res.Meta.Message = "Mandatory request data"

		ctx.JSON(http.StatusOK, res)
		return
	}

	// validate request
	header, resultValidate := ctrl.ValidateRequestWithoutAuth(ctx, req)
	if !resultValidate.Meta.Status {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[ValidateRequestWithoutAuth]-[Error : %v]", resultValidate))
		logrus.Println(logReq)

		ctx.JSON(http.StatusOK, resultValidate)
		return
	}

	param.InstitutionID = header.InstitutionID
	param.RRN = req.ReferenceId

	logrus.Println("[Request]")
	logrus.Info("[ReferenceId : %v]", req.ReferenceId)

	res = payment.ReversalPaymentService(req, param, header)

	ctx.JSON(http.StatusOK, res)

}
