package host

import (
	"encoding/json"
	"net/http"
	https "ottopoint-purchase/hosts"
	authModel "ottopoint-purchase/hosts/auth/models"

	"ottopoint-purchase/utils"

	"github.com/astaxie/beego/logs"
)

var (
	host                         string
	endpointClearCacheGetBalance string
	// HealthCheckKey string
)

func init() {
	host = utils.GetEnv("HOST_ADDRESS_OTTOPOINT_AUTH", "http://13.228.25.85:8666")
	endpointClearCacheGetBalance = utils.GetEnv("ENDPOINT_CLEAR_CACHE_BALANCE", "/auth/v2/clear-cache-balance-point")
	// HealthCheckKey = utils.GetEnv("OTTOPOINT_PURCHASE_KEY_HEALTHCHECK_REDIS_TOKEN", "OTTOPOINT-PURCHASE:REDIS_TOKEN_OTTOPOINT")
}

func ClearCacheBalance(phone string) (authModel.RespClearCacheBalance, error) {
	var result authModel.RespClearCacheBalance

	urlSvr := host + endpointClearCacheGetBalance + "?phone=" + phone

	header := make(http.Header)
	header.Set("Content-Type", "application/json")

	data, err := https.HTTPxGET(urlSvr, header)
	if err != nil {
		logs.Error("Check error", err.Error())

		return result, err
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		logs.Error("Failed to unmarshaling response from open-loyalty ", err.Error())

		return result, err
	}

	return result, nil
}
