package redis

import (
	"errors"
	"fmt"
	"time"

	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"

	"gopkg.in/redis.v5"
)

var (
	ClientRed       *redis.ClusterClient
	addres          string
	addresscluster1 string
	addresscluster2 string
	addresscluster3 string
	port            string
	dbtype          int
	queuename       string
)

func init() {

	addresscluster1 = utils.GetEnv("REDIS_HOST_CLUSTER1", "13.228.23.160:8079")
	addresscluster2 = utils.GetEnv("REDIS_HOST_CLUSTER2", "13.228.23.160:8078")
	addresscluster3 = utils.GetEnv("REDIS_HOST_CLUSTER3", "13.228.23.160:8077")
	ClientRed = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{addresscluster1, addresscluster2, addresscluster3},
	})
	pong, err := ClientRed.Ping().Result()

	//logger.Info("Redis Status ", zap.String("Ping",pong), zap.String("Error",err.Error()))

	fmt.Println("Redis Ping ", pong)
	fmt.Println("Redis Ping ", err)

	addres = utils.GetEnv("REDIS_HOST_OTTOPOINT_PURCHASE", "13.228.23.160")
	port = utils.GetEnv("REDIS_PORT_OTTOPOINT_PURCHASE", "6077")
	// dbtype = 0
	queuename = utils.GetEnv("REDIS_QUEUE", "ottopoint-purchase")

	// ClientRed = redis.NewClient(&redis.Options{
	// 	Addr:     addres + ":" + port,
	// 	Password: "",     // no password set
	// 	DB:       dbtype, // use default DB
	// })

	// pong, err := ClientRed.Ping().Result()
	// fmt.Println("Redis Ping ", pong)
	// fmt.Println("Redis Ping ", err)
	// //RunSubscriber()
}

// GetRedisConnection ...
func GetRedisConnection() *redis.ClusterClient {
	return ClientRed
}

// GetRedisUri ...
func GetRedisUri() string {
	return "redis://" + addres + ":" + port + "/"
}

// GetQueueName ...
func GetQueueName() string {
	return queuename
}

/*
 Redis Standard Set
*/
func SaveRedis(key string, val interface{}) error {
	var err error
	for i := 0; i < 3; i++ {
		err = ClientRed.Set(key, val, 0).Err()
		if err == nil {
			break
		}
	}
	return err
}

/*
 Redis Standard Get
*/
func GetRedisKey(Key string) (string, error) {
	val2, err := ClientRed.Get(Key).Result()
	if err == redis.Nil {
		err = errors.New("Key Does Not Exists")
		fmt.Println("keystruct does not exists")
	} else if err != nil {
		fmt.Println("Error : ", err.Error())
	} //else {
	//fmt.Println("keystruct", val2)
	//}
	return val2, err
}

func DelRedisKey(key string) error {
	return ClientRed.Del(key).Err()
}

/*
delayto * max = total timeout
*/
func GetDataRedis(key string, delayto, max int) (bool, string) {
	for i := 0; i < max; i++ {
		data, err := GetRedisKey(key)
		fmt.Println(" Err : ", err)
		fmt.Println(" data : ", data)
		if err == nil {
			return true, data
		}
		time.Sleep(time.Duration(delayto) * time.Second)
	}
	return false, ""
}

/*
 Redis Standard Set Expired
*/
func SaveRedisExp(key string, menit string, val interface{}) error {
	var err error
	for i := 0; i < 3; i++ {
		duration, _ := time.ParseDuration(menit)
		err = ClientRed.Set(key, val, duration).Err()
		if err == nil {
			break
		}
		fmt.Println("Error : ", err)
	}
	return err
}

func SaveRedisCounter(key string) (int64, error) {
	incr := ClientRed.Incr(key)
	return incr.Val(), incr.Err()
}

func SaveRedisCounterAuto(key string, autonom int64) (int64, error) {
	incr := ClientRed.IncrBy(key, autonom)
	return incr.Val(), incr.Err()
}

func GetRedisCounter(key string) (int64, error) {
	decr := ClientRed.Decr(key)
	return decr.Val(), decr.Err()

}

func GetRedisCounterIncr(key string) (int64, error) {
	decr := ClientRed.Incr(key)
	return decr.Val(), decr.Err()

}

func GetRrn(rrnkey string) string {
	//res, err := redis.GetRedisKey(rrnkey)
	counter, err := SaveRedisCounter(rrnkey)
	if err != nil {
		counter = 1
	}
	t11 := time.Now().Local()
	return fmt.Sprintf("%02d%02d%02d%02d%04d", t11.Day(), t11.Hour(), t11.Minute(), t11.Second(), counter)
}

// GetRedisClusterHealthCheck ..
func GetRedisClusterHealthCheck() models.RedisHealthcheckResponse {
	return models.RedisHealthcheckResponse{
		Type:    "Cluster",
		Address: addres,
		Port:    port,
	}
}
