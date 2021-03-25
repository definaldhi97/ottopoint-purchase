package callbacks

import (
	"fmt"
	"net/http"
	"ottopoint-purchase/models"
	callback "ottopoint-purchase/models/v21/callback"

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

	// if req.VoucherType == constants.VoucherTypePPOB {
	// 	req.Data = callback.DataVoucherTypePPOB{}
	// } else if req.VoucherType == constants.VoucherTypeVoucherCode {
	// 	req.Data = callback.DataVoucherTypeVoucherCode{}
	// } else {

	// 	res.Meta.Code = 03
	// 	res.Meta.Message = "Invalid VoucherType"

	// 	logrus.Error(namectrl)
	// 	logrus.Error(fmt.Sprintf("[Invalid VoucherType]-[Request : %v]", req))
	// 	logrus.Println(logReq)

	// 	ctx.JSON(http.StatusOK, res)
	// 	return
	// }

	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.Meta.Code = 03
		res.Meta.Message = err.Error()

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[ShouldBindJSON]-[Error : %v]", err))
		logrus.Println(logReq)

		ctx.JSON(http.StatusOK, res)
		return
	}

	logrus.Println("[Request]")
	logrus.Info("InstitutionId : ", req.InstitutionId, "NotificationType : ", req.NotificationType, "NotificationTo : ", req.NotificationTo, "TransactionId : ", req.TransactionId, "VoucherType : ", req.VoucherType, "Data : ", req.Data)

	res = service.CallbackVoucherAG_V21_Service(req)

	ctx.JSON(http.StatusOK, res)

}
