package controllers

import (
	"fmt"
	"net/http"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	"time"

	auth "ottopoint-purchase/hosts/auth/host"
	opl "ottopoint-purchase/hosts/opl/host"
	ottoag "ottopoint-purchase/hosts/ottoag/host"
	ottomart "ottopoint-purchase/hosts/ottomart/host"
	publisher "ottopoint-purchase/hosts/publisher/host"
	redisToken "ottopoint-purchase/hosts/redis_token/host"
	signature "ottopoint-purchase/hosts/signature/host"
	uv "ottopoint-purchase/hosts/ultra_voucher/host"
	"ottopoint-purchase/models"
	"ottopoint-purchase/redis"
	"ottopoint-purchase/utils"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	zaplog "github.com/opentracing-contrib/go-zap/log"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	hcmodels "ottodigital.id/library/healthcheck/models"
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

func getHealthCheckStatus() hcmodels.HealthCheckResponse {
	// redis
	redisHc := make([]hcmodels.RedisHealthCheck, 0)
	redisHc = append(redisHc, redis.GetRedisClusterHealthCheck())
	// TODO more redis health check

	// database
	databaseHc := make([]hcmodels.DatabaseHealthCheck, 0)
	databaseHc = append(databaseHc, db.GetHealthCheck())
	// TODO more database health check

	// service
	serviceHc := make([]hcmodels.ServiceHealthCheck, 0)
	serviceHc = append(serviceHc, opl.GetServiceHealthCheckOPL())
	serviceHc = append(serviceHc, ottoag.GetServiceHealthCheckOttoAG())
	serviceHc = append(serviceHc, redisToken.GetServiceHealthCheckRedisService())
	serviceHc = append(serviceHc, signature.GetServiceHealthCheckSignature())
	serviceHc = append(serviceHc, auth.GetServiceHealthCheckAuth())
	serviceHc = append(serviceHc, ottomart.GetServiceHealthCheckOttomart())
	serviceHc = append(serviceHc, uv.GetServiceHealthCheckUV())
	serviceHc = append(serviceHc, publisher.GetServiceHealthCheckPublisher())
	// TODO more service health check

	return hcmodels.HealthCheckResponse{
		Redis:    redisHc,
		Database: databaseHc,
		Service:  serviceHc,
	}
}
