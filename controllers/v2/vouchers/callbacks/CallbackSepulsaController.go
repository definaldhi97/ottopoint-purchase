package callbacks

import (
	"fmt"
	"net/http"
	"ottopoint-purchase/models"

	sepulsaModels "ottopoint-purchase/hosts/sepulsa/models"
	services "ottopoint-purchase/services/v2/vouchers/callbacks"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CallbackSepulsaController(ctx *gin.Context) {

	req := sepulsaModels.CallbackTrxReq{}
	res := models.Response{}

	namectrl := "[PackageCallBacksController]-[CallbackSepulsaController]"

	logReq := fmt.Sprintf("[TransactionID : %v, CustomerNumber : %v]", req.TransactionID, req.CustomerNumber)

	logrus.Info(namectrl)

	if err := ctx.ShouldBindJSON(&req); err != nil {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[ShouldBindJSON]-[Error : %v]", err))
		logrus.Println(logReq)

		res.Meta.Code = 03
		res.Meta.Message = err.Error()
		ctx.JSON(http.StatusOK, res)
		return
	}

	logrus.Println("[Request]")
	logrus.Info("TransactionID : ", req.TransactionID, "Type : ", req.Type, "Created : ", req.Created, "Changed : ", req.Changed, "CustomerNumber : ", req.CustomerNumber,
		"OrderID : ", req.OrderID, "Price : ", req.Price, "Status : ", req.Status, "ResponseCode : ", req.ResponseCode, "SerialNumber : ", req.SerialNumber,
		"Amount : ", req.Amount, "ProductID : ", req.ProductID, "Token : ", req.Token, "Data : ", req.Data)

	res = services.CallbackVoucherSepulsaService(req)

	ctx.JSON(http.StatusOK, res)

}
