package controllers

import (
	"fmt"
	"net/http"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"

	"ottopoint-purchase/models"
	"ottopoint-purchase/redis"
	"ottopoint-purchase/utils"

	"github.com/gin-gonic/gin"
)

func HealthCheckService(ctx *gin.Context) {
	fmt.Println(">>> Health Check - Service <<<")

	response := models.Response{
		Meta: utils.GetMetaResponse("default"),
	}

	response = models.Response{
		Meta: utils.GetMetaResponse(constants.KeyResponseSucceed),
		Data: models.HealthcheckResponse{
			Redis:    redis.GetRedisClusterHealthCheck(),
			Database: db.GetHealthCheck(),
			// Service:  serviceHc,
		},
	}

	ctx.JSON(http.StatusOK, response)
	return
}
