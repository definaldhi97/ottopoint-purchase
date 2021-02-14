package check_status

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	services "ottopoint-purchase/services/v2/vouchers/check_status"

	"net/http"
)

func SchedulerCheckStatusController(ctx *gin.Context) {
	// res := models.Response{}

	namectrl := "[PackageCheckStatus]-[SchedulerCheckStatusController]"

	logrus.Info(namectrl)

	res := services.SchedulerCheckStatusService()

	ctx.JSON(http.StatusOK, res)
	return
}
