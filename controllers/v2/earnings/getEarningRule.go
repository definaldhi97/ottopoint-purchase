package earnings

import (
	"errors"
	earning "ottopoint-purchase/services/v2/earnings"
	"ottopoint-purchase/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"net/http"

	"ottopoint-purchase/models"
)

func GetEarningRuleController(ctx *gin.Context) {
	res := models.Response{}

	namectrl := "[PackageEarnings]-[GetEarningRuleController]"

	// Debug, Println (Putih), Warn (Kuning), Error (Merah), Info (Biru)

	productCode := ctx.Request.URL.Query().Get("code")

	if productCode == "" {

		logrus.Error(namectrl)
		logrus.Error("[ValidateRequestMandatory]-[Invalid Mandatory]")
		logrus.Println("Request : ", productCode)

		res = utils.GetMessageResponse(res, 61, false, errors.New("Invalid Mandatory"))

		ctx.JSON(http.StatusOK, res)

		return

	}

	logrus.Println("[Request]")
	logrus.Info("Code Earning : ", productCode)

	getEarningRule := new(earning.EarningsServices)
	res = getEarningRule.NewGetEarningRuleService(productCode)

	ctx.JSON(http.StatusOK, res)

}
