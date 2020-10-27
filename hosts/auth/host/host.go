package host

import (
	"encoding/json"
	authModel "ottopoint-purchase/hosts/auth/models"
	"time"

	hcmodels "ottodigital.id/library/healthcheck/models"

	"github.com/astaxie/beego/logs"
	ODU "ottodigital.id/library/utils"
)

var (
	host                         string
	name                         string
	endpointClearCacheGetBalance string
	// HealthCheckKey string
)

func init() {
	host = ODU.GetEnv("HOST_ADDRESS_OTTOPOINT_AUTH", "http://13.228.25.85:8666")
	name = ODU.GetEnv("NAME_OTTOPOINT_AUTH", "OTTOPOINT-AUTH")
	endpointClearCacheGetBalance = ODU.GetEnv("ENDPOINT_CLEAR_CACHE_BALANCE", "/auth/v2/clear-cache-balance-point")
	// HealthCheckKey = ODU.GetEnv("OTTOPOINT_PURCHASE_KEY_HEALTHCHECK_REDIS_TOKEN", "OTTOPOINT-PURCHASE:REDIS_TOKEN_OTTOPOINT")
}

func ClearCacheBalance(phone string) (authModel.RespClearCacheBalance, error) {
	var result authModel.RespClearCacheBalance

	urlSvr := host + endpointClearCacheGetBalance + "?phone=" + phone

	// header := make(http.Header)
	// header.Set("DeviceId", headerReq.DeviceID)
	// header.Set("InstitutionId", headerReq.InstitutionID)
	// header.Set("Geolocation", headerReq.Geolocation)
	// header.Set("ChannelId", headerReq.ChannelID)
	// header.Set("AppsId", headerReq.AppsID)
	// header.Set("Timestamp", headerReq.Timestamp)
	// header.Set("Authorization", headerReq.Authorization)

	data, err := HTTPGet(urlSvr, nil)
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

func GetServiceHealthCheckAuth() hcmodels.ServiceHealthCheck {
	return hcmodels.ServiceHealthCheck{
		Name:    name,
		Address: host,
		Status:  "OK",
		// Description: ,
		UpdatedAt: time.Now(),
	}
}
