package controllers

import (
	"fmt"
	"ottopoint-purchase/constants"
	opl "ottopoint-purchase/hosts/opl/host"
	ottoag "ottopoint-purchase/hosts/ottoag/host"
	redisToken "ottopoint-purchase/hosts/redis_token/host"
	signature "ottopoint-purchase/hosts/signature/host"

	"ottopoint-purchase/models"
	"ottopoint-purchase/redis"
	"ottopoint-purchase/utils"

	"github.com/opentracing/opentracing-go"
	hcmodels "ottodigital.id/library/healthcheck/models"
)

type HealthCheckService struct {
	General models.GeneralModel
}

func InitializeHealthCheckService(general models.GeneralModel) *HealthCheckService {
	return &HealthCheckService{
		General: general,
	}
}

func (service *HealthCheckService) HealthCheck() models.Response {
	fmt.Println(">>> Health Check - Service <<<")

	response := models.Response{
		Meta: utils.GetMetaResponse("default"),
	}

	sugarLogger := service.General.OttoZaplog
	sugarLogger.Info("Service:HealthCheck ")
	span, _ := opentracing.StartSpanFromContext(service.General.Context, "Service:HealthCheck")
	defer span.Finish()

	data := getHealthCheckStatus()

	response = models.Response{
		Meta: utils.GetMetaResponse(constants.KeyResponseSucceed),
		Data: data,
	}

	return response
}

func getHealthCheckStatus() hcmodels.HealthCheckResponse {
	// redis
	redisHc := make([]hcmodels.RedisHealthCheck, 0)
	redisHc = append(redisHc, redis.GetRedisClusterHealthCheck())
	// TODO more redis health check

	// database
	databaseHc := make([]hcmodels.DatabaseHealthCheck, 0)
	// databaseHc = append(databaseHc, db.GetDatabaseHealthCheck())
	// TODO more database health check

	// service
	serviceHc := make([]hcmodels.ServiceHealthCheck, 0)
	serviceHc = append(serviceHc, opl.GetServiceHealthCheck())
	serviceHc = append(serviceHc, ottoag.GetServiceHealthCheck())
	serviceHc = append(serviceHc, redisToken.GetServiceHealthCheck())
	serviceHc = append(serviceHc, signature.GetServiceHealthCheck())
	// TODO more service health check

	return hcmodels.HealthCheckResponse{
		Redis:    redisHc,
		Database: databaseHc,
		Service:  serviceHc,
	}
}
