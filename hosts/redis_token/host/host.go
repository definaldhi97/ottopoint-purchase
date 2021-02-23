package host

import (
	"encoding/json"
	"net/http"
	redismodels "ottopoint-purchase/hosts/redis_token/models"
	"ottopoint-purchase/models"
	"strings"

	https "ottopoint-purchase/hosts"
	"ottopoint-purchase/utils"

	"github.com/astaxie/beego/logs"
	"github.com/sirupsen/logrus"
)

var (
	host          string
	name          string
	endpointToken string
	// HealthCheckKey string
)

func init() {
	host = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_REIDS_TOKEN", "http://13.228.25.85:8703")
	name = utils.GetEnv("OTTOPOINT_PURCHASE_NAME_REDIS_TOKEN", "REDIS-TOKEN-OTTOPOINT")
	endpointToken = utils.GetEnv("OTTOPOINT_PURCHASE_ENDPOINT_REDIS_TOKEN", "/ottopoint/v0.1.0/redis/service")
	// HealthCheckKey = utils.GetEnv("OTTOPOINT_PURCHASE_KEY_HEALTHCHECK_REDIS_TOKEN", "OTTOPOINT-PURCHASE:REDIS_TOKEN_OTTOPOINT")
}

func CheckToken(headers models.RequestHeader) (redismodels.TokenResp, error) {
	var resp redismodels.TokenResp

	urlSvr := host + endpointToken

	t := strings.ReplaceAll(headers.Authorization, "Bearer ", "")
	token := headers.InstitutionID + "-" + t
	logrus.Info("Token : ", token)

	header := make(http.Header)
	header.Set("Key", token)
	header.Set("Action", "GET")

	data, err := https.HTTPxPOSTwithoutRequest(urlSvr, header)
	// data, err := HTTPxFormWithHeader(urlSvr, token)
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

	header := make(http.Header)
	// header.Set("Content-Type", "application/json")
	header.Set("Action", "GET")
	header.Set("Key", Key)

	data, err := https.HTTPxPOSTwithoutRequest(url, header)
	// data, err := HTTPPostWithHeader_GetRedis(url, Key)

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
