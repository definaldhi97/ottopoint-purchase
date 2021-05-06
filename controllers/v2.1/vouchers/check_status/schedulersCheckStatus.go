package check_status

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	services "ottopoint-purchase/services/v2.1/vouchers/check_status"

	"net/http"
)

func SchedulerCheckStatusV21Controller(ctx *gin.Context) {
	// res := models.Response{}

	namectrl := "[PackageCheckStatus]-[SchedulerCheckStatusV21Controller]"

	logrus.Info(namectrl)

	res := services.SchedulerCheckStatusServiceV21()

	ctx.JSON(http.StatusOK, res)
	return
}
