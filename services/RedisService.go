package services

import (
	redishost "ottopoint-purchase/integration/redis/host"
	redismodels "ottopoint-purchase/integration/redis/models"
	"ottopoint-purchase/models"
)

type RedisService struct {
	General models.GeneralModel
}

func (service *RedisService) GetData(key string) redismodels.ResponseRedis1 {

	//keyRedis := header.InstitutionID + "-" + header.Authorization
	dataRedis, _ := redishost.GetToken(key)

	return *dataRedis

}
