package callbacks

import (
	"fmt"
	"net/http"
	"ottopoint-purchase/controllers"
	"ottopoint-purchase/models"
	callback "ottopoint-purchase/models/v21/callback"
	"time"

	service "ottopoint-purchase/services/v2.1/vouchers/callbacks"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CallBackVoucherAG_V21_Controller(ctx *gin.Context) {

	req := callback.CallbackVoucherAGReq{}
	res := models.Response{}

	namectrl := "[PackageCallbacks_V21]-[CallBackVoucherAG_V21_Controller]"

	logReq := fmt.Sprintf("[TransactionID : %v]", req.TransactionId)

	time.Sleep(time.Second * 5)

	fmt.Println(">>> Sleep 5 Detik <<<")

	logrus.Info(namectrl)

	clientIP := ctx.ClientIP()
	if clientIP != "" {
		fmt.Println("Caller IP : ", clientIP)
	}

	fmt.Println("From IP : ", clientIP)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.Meta.Code = 03
		res.Meta.Message = err.Error()

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[ShouldBindJSON]-[Error : %v]", err))
		logrus.Println(logReq)

		ctx.JSON(http.StatusOK, res)
		return
	}

	// validate request
	_, resultValidate := controllers.ValidateRequest(ctx, false, req, false)
	if !resultValidate.Meta.Status {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[ValidateRequest]-[Error : %v]", resultValidate))
		logrus.Println(logReq)

		ctx.JSON(http.StatusOK, resultValidate)
		return
	}

	logrus.Println("[Request]")
	logrus.Info("InstitutionId : ", req.InstitutionId, " NotificationType : ", req.NotificationType, " TransactionId : ", req.TransactionId, " VoucherType : ", req.VoucherType, " Data : ", req.Data)

	res = service.CallbackVoucherAG_V21_Service(req)

	ctx.JSON(http.StatusOK, res)

}
