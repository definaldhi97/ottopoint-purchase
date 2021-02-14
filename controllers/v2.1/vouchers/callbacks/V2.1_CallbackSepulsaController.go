package callbacks

import (
	"fmt"
	"net/http"
	"ottopoint-purchase/models"

	service "ottopoint-purchase/services/v2.1/vouchers/callbacks"

	sepulsaModels "ottopoint-purchase/hosts/sepulsa/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// func (controller *V21_CallbackSepulsaController) V21_VoucherCallbackSepulsaController(ctx *gin.Context) {
func CallBackSepulsa_V21_Controller(ctx *gin.Context) {

	// logrus.Info("[ >>>>>>>>>>>>>>>>>>>>> Callbakc Sepulsa COntroller <<<<<<<<<<<<<<<<<< ]")

	req := sepulsaModels.CallbackTrxReq{}
	res := models.Response{}

	namectrl := "[PackageCallbacks_V21]-[CallBackSepulsa_V21_Controller]"

	logReq := fmt.Sprintf("[TransactionID : %v, CustomerNumber : %v]", req.TransactionID, req.CustomerNumber)

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

	logrus.Println("[Request]")
	logrus.Info("TransactionID : ", req.TransactionID, "Type : ", req.Type, "Created : ", req.Created, "Changed : ", req.Changed, "CustomerNumber : ", req.CustomerNumber,
		"OrderID : ", req.OrderID, "Price : ", req.Price, "Status : ", req.Status, "ResponseCode : ", req.ResponseCode, "SerialNumber : ", req.SerialNumber,
		"Amount : ", req.Amount, "ProductID : ", req.ProductID, "Token : ", req.Token, "Data : ", req.Data)

	res = service.CallbackVoucherSepulsa_V21_Service(req)

	ctx.JSON(http.StatusOK, res)

}
