package host

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	redismodels "ottopoint-purchase/hosts/redis_token/models"
	"ottopoint-purchase/models"
	"strings"
	"time"

	hcmodels "ottodigital.id/library/healthcheck/models"

	"github.com/astaxie/beego/logs"
	ODU "ottodigital.id/library/utils"
)

var (
	host          string
	name          string
	endpointToken string
	// HealthCheckKey string
)

func init() {
	host = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_REIDS_TOKEN", "http://13.228.25.85:8703")
	name = ODU.GetEnv("OTTOPOINT_PURCHASE_NAME_REDIS_TOKEN", "REDIS-TOKEN-OTTOPOINT")
	endpointToken = ODU.GetEnv("OTTOPOINT_PURCHASE_ENDPOINT_REDIS_TOKEN", "/ottopoint/v0.1.0/redis/service")
	// HealthCheckKey = ODU.GetEnv("OTTOPOINT_PURCHASE_KEY_HEALTHCHECK_REDIS_TOKEN", "OTTOPOINT-PURCHASE:REDIS_TOKEN_OTTOPOINT")
}

func CheckToken(header models.RequestHeader) (redismodels.TokenResp, error) {
	var resp redismodels.TokenResp

	urlSvr := host + endpointToken

	t := strings.ReplaceAll(header.Authorization, "Bearer ", "")
	token := header.InstitutionID + "-" + t
	logs.Info("Token : ", token)

	data, err := HTTPxFormWithHeader(urlSvr, token)
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

// get token from redis
func GetToken(Key string) (*redismodels.TokenResp, error) {
	var resp redismodels.TokenResp

	url := host + endpointToken
	data, err := HTTPPostWithHeader_GetRedis(url, Key)

	if err != nil {
		logs.Error("generate mpan ", err.Error())
		return &resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logs.Error("Failed to unmarshaling response  from Redis service Ottopoint ", err.Error())
		return &resp, err
	}

	return &resp, err

}

// GetServiceHealthCheck ..
func GetServiceHealthCheckRedisService() hcmodels.ServiceHealthCheck {
	res := hcmodels.ServiceHealthCheck{}
	var erorr interface{}
	// sugarLogger := service.General.OttoZapLog

	PublicAddress := host
	log.Print("url : ", PublicAddress)
	res.Name = name
	res.Address = PublicAddress
	res.UpdatedAt = time.Now().UTC()

	d, err := http.Get(PublicAddress)

	erorr = err
	if err != nil {
		log.Print("masuk error")
		res.Status = "Not OK"
		res.Description = fmt.Sprintf("%v", erorr)
		return res
	}
	if d.StatusCode != 200 {
		res.Status = "Not OK"
		res.Description = d.Status
		return res
	}

	res.Status = "OK"
	res.Description = ""

	return res
}
