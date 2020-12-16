package controllers

import (
	"fmt"
	"net/http"
	"ottopoint-purchase/constants"
	databaseHc "ottopoint-purchase/db"
	redishc "ottopoint-purchase/redis"
	"time"

	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	zaplog "github.com/opentracing-contrib/go-zap/log"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	ottologer "ottodigital.id/library/logger"
	utilsgo "ottodigital.id/library/utils"
)

func HealthCheckService(ctx *gin.Context) {
	fmt.Println(">>> Health Check - Service <<<")

	response := models.Response{
		Meta: utils.GetMetaResponse("default"),
	}

	sugarLogger := ottologer.GetLogger()
	namectrl := "[HealthCheckService]"
	span := TracingFirstControllerCtx(ctx, "", namectrl)
	c := ctx.Request.Context()
	context := opentracing.ContextWithSpan(c, span)
	spanid := utilsgo.GetSpanId(span)
	logs.Info("context :", context)

	data := getHealthCheckStatus()

	response = models.Response{
		Meta: utils.GetMetaResponse(constants.KeyResponseSucceed),
		Data: data,
	}

	sugarLogger.Info("RESPONSE:", zap.String("SPANID", spanid), zap.String("CTRL", namectrl),
		zap.Any("BODY", response))

	datalog := utils.LogSpanMax(response)
	zaplog.InfoWithSpan(span, namectrl,
		zap.Any("RESP", datalog),
		zap.Duration("backoff", time.Second))

	defer span.Finish()
	ctx.JSON(http.StatusOK, response)
}

func getHealthCheckStatus() *models.HealthcheckResponse {

	// //service
	// serviceHc := make([]models.ServicesHealthcheckResponse, 0)
	// serviceHc = append(serviceHc, hostopl.GetHealthCheckOPL())

	return &models.HealthcheckResponse{
		Redis:    redishc.GetRedisHealthCheck(),
		Database: databaseHc.GetDatabaseHealthCheck(),
		// Services: serviceHc,
	}
}
