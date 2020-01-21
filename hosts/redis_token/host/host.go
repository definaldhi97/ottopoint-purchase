package host

import (
	"encoding/json"
	redismodels "ottopoint-purchase/hosts/redis_token/models"
	"ottopoint-purchase/models"
	"ottopoint-purchase/redis"

	hcmodels "ottodigital.id/library/healthcheck/models"
	hcutils "ottodigital.id/library/healthcheck/utils"

	"github.com/astaxie/beego/logs"
	ODU "ottodigital.id/library/utils"
)

var (
	host           string
	name           string
	endpointToken  string
	HealthCheckKey string
)

func init() {
	host = ODU.GetEnv("host.openloyalty", "http://13.228.25.85:8703")
	name = ODU.GetEnv("name.redis.token", "REDIS-TOKEN-OTTOPOINT")
	endpointToken = ODU.GetEnv("host.voucher_redeem", "/ottopoint/v0.1.0/redis/service")

	HealthCheckKey = ODU.GetEnv("key.healthcheck.redis.token", "REDIS-TOKEN-OTTOPOINT:REDIS_TOKEN_OTTOPOINT")
}

func CheckToken(header models.RequestHeader) (redismodels.TokenResp, error) {
	var resp redismodels.TokenResp

	urlSvr := host + endpointToken

	token := header.InstitutionID + "-" + header.Authorization

	data, err := HTTPxFormWithHeader(urlSvr, token, HealthCheckKey)
	if err != nil {
		logs.Error("Check error", err.Error())

		return resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logs.Error("Failed to unmarshaling response from open-loyalty ", err.Error())

		return resp, err
	}

	return resp, nil
}

// GetServiceHealthCheck ..
func GetServiceHealthCheck() hcmodels.ServiceHealthCheck {
	redisClient := redis.GetRedisConnection()
	return hcutils.GetServiceHealthCheck(&redisClient, &hcmodels.ServiceEnv{
		Name:           name,
		Address:        host,
		HealthCheckKey: HealthCheckKey,
	})
}
