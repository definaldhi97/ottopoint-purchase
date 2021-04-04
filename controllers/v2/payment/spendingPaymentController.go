package payment

import (
	"errors"
	"fmt"
	"os"
	ctrl "ottopoint-purchase/controllers"
	"ottopoint-purchase/db"
	sp "ottopoint-purchase/models/v2/payment"
	"ottopoint-purchase/utils"

	"net/http"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/models"
	"ottopoint-purchase/services/v2/payment"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func SpendingPaymentController(ctx *gin.Context) {

	req := sp.SpendingPaymentReq{}
	res := models.Response{}
	param := models.Params{}

	namectrl := "[PackagePayment]-[SpendingPaymentController]"
	logReq := fmt.Sprintf("[AccountNumber : %v || ReferenceId : %v]", req.AccountNumber, req.ReferenceId)

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

	if req.TransType != constants.CodeRedeemtion && req.TransType != constants.CodeSplitBill {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[Invalid Mandatory]-[TransType : %v]", req.TransType))
		logrus.Println(logReq)

		res.Meta.Code = 196
		res.Meta.Message = "Mandatory request data"

		ctx.JSON(http.StatusOK, res)
		return
	}

	if req.Point == 0 || req.Cash == 0 || req.ReferenceId == "" {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[Invalid Mandatory]-[Point/Cash/ReferenceId]"))
		logrus.Println(logReq)

		res.Meta.Code = 196
		res.Meta.Message = "Mandatory request data"

		ctx.JSON(http.StatusOK, res)
		return
	}

	if req.Point+req.Cash != req.Amount {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[No Match Amount]-[Amount]"))
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

	// check user
	dataUser, errUser := db.UserWithInstitution(req.AccountNumber, header.InstitutionID)
	if errUser != nil || dataUser.CustID == "" {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[UserWithInstitution]-[Error : %v]", errUser))
		logrus.Println(logReq)

		res = utils.GetMessageResponse(res, 404, false, errors.New("User belum Eligible"))

		ctx.JSON(http.StatusOK, res)
		return
	}

	param.AccountNumber = req.AccountNumber
	param.AccountId = dataUser.CustID
	param.InstitutionID = header.InstitutionID
	param.Point = req.Point
	param.Amount = int64(req.Amount)
	param.RRN = req.ReferenceId

	logrus.Println("[Request]")
	logrus.Info("[AccountNumber : ", req.AccountNumber, "|| TransType : ", req.TransType, "|| ProductName : ", req.ProductName, "|| ReferenceId : ", req.ReferenceId,
		"|| TransactionTime : ", req.TransactionTime, "|| Point : ", req.Point, "|| Amount : ", req.Amount, "|| Comment : ", req.Comment, "|| PaymentMethod : ", req.PaymentMethod)

	res = payment.SpendingPaymentService(req, param, header)

	ctx.JSON(http.StatusOK, res)

}
