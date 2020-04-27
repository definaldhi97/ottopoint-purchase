package host

import (
	"encoding/json"
	"fmt"
	"log"
	"ottopoint-purchase/integration/redis"
	redismodels "ottopoint-purchase/integration/redis/models"

	"github.com/astaxie/beego/logs"
	"github.com/kelseyhightower/envconfig"
)

type HostRedisServiceOttopoint struct {
	Host                 string `envconfig:"HOSTADDRES_REDIS_SERVICE" default:"http://13.228.25.85:8703"`
	EndpointServiceRedis string `envconfig:"ENDPOINT_SERVICEREDIS" default:"/ottopoint/v0.1.0/redis/service"`
}

var (
	hostRedisServiceOttopoint HostRedisServiceOttopoint
)

func init() {
	err := envconfig.Process("REDIS_SERVICE_OTTOPOINT", &hostRedisServiceOttopoint)
	if err != nil {
		fmt.Println("Failed to get REDIS_SERVICE_OTTOPOINT env:", err)
	}
}

// save token to redis
func SaveToken(Key, Expire, value string) (*redismodels.ResponseApi, error) {
	var resp redismodels.ResponseApi

	// reqHeader := http.Header{}
	// reqHeader.Set("Action", "SET")
	// reqHeader.Set("Expire", Expire)
	// reqHeader.Set("Key", Key)

	jsonData := map[string]interface{}{
		"value": value,
	}

	url := hostRedisServiceOttopoint.Host + hostRedisServiceOttopoint.EndpointServiceRedis
	log.Print("ini url ", hostRedisServiceOttopoint.EndpointServiceRedis)

	data, err := redis.HTTPPostWithHeader_SaveRedis(url, jsonData, Key, Expire)

	if err != nil {
		logs.Error("generate mid ", err.Error())
		return &resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logs.Error("Failed to unmarshaling response GenMID from StringBuilder ", err.Error())
		return &resp, err
	}

	return &resp, err
}

// get token from redis
func GetToken(Key string) (*redismodels.ResponseRedis1, error) {
	var resp redismodels.ResponseRedis1
	// reqHeader := http.Header{}
	// reqHeader.Set("Action", "GET")
	// reqHeader.Set("Key", Key)

	url := hostRedisServiceOttopoint.Host + hostRedisServiceOttopoint.EndpointServiceRedis
	data, err := redis.HTTPPostWithHeader_GetRedis(url, Key)

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
