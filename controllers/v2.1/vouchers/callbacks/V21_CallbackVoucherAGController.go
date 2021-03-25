package callbacks

import (
	"fmt"
	"net/http"
	"ottopoint-purchase/controllers"
	redishost "ottopoint-purchase/hosts/redis_token/host"
	"ottopoint-purchase/models"
	callback "ottopoint-purchase/models/v21/callback"
	"ottopoint-purchase/utils"

	service "ottopoint-purchase/services/v2.1/vouchers/callbacks"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CallBackVoucherAG_V21_Controller(ctx *gin.Context) {

	req := callback.CallbackVoucherAGReq{}
	res := models.Response{}

	namectrl := "[PackageCallbacks_V21]-[CallBackVoucherAG_V21_Controller]"

	logReq := fmt.Sprintf("[TransactionID : %v]", req.TransactionId)

	logrus.Info(namectrl)

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
	header, resultValidate := controllers.ValidateRequest(ctx, true, req, false)
	if !resultValidate.Meta.Status {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[ValidateRequest]-[Error : %v]", resultValidate))
		logrus.Println(logReq)

		ctx.JSON(http.StatusOK, resultValidate)
		return
	}

	// get customer di redis
	_, errToken := redishost.CheckToken(header)
	if errToken != nil {
		logrus.Println("Failed Get Token .. ..")

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[CheckToken]-[Error : %v]", errToken))
		logrus.Println(logReq)

		res = utils.GetMessageFailedErrorNew(res, 500, "Internal Server Error")
		ctx.JSON(http.StatusOK, res)
		return
	}

	logrus.Println("[Request]")
	logrus.Info("InstitutionId : ", req.InstitutionId, "NotificationType : ", req.NotificationType, "TransactionId : ", req.TransactionId, "VoucherType : ", req.VoucherType, "Data : ", req.Data)

	res = service.CallbackVoucherAG_V21_Service(req)

	ctx.JSON(http.StatusOK, res)

}
