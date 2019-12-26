package utils

import (
	"encoding/json"
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/models"
	"ottopoint-purchase/redis"
	"time"

	"github.com/leekchan/accounting"

	ODU "ottodigital.id/library/utils"
)

var (
	rrnkey            string
	DefaultStatusCode int
	DefaultStatusMsg  string
	RedisKeyAuth      string
	LimitTRXPoint     string
)

func init() {
	DefaultStatusCode = 200
	DefaultStatusMsg = "OK"
	rrnkey = ODU.GetEnv("REDISKEY.OTTOFIN.RRN", "OTTOFIN:KEYRRN")
	RedisKeyAuth = ODU.GetEnv("redis.key.auth", "Ottopoint-Token-Admin :")
	LimitTRXPoint = ODU.GetEnv("limit.trx.point", "999999999999999")
}

func GetMessageResponse(res models.Response, code int, status bool, err error) models.Response {

	res = models.Response{}

	res.Meta.Code = code
	res.Meta.Status = status
	res.Meta.Message = err.Error()

	return res
}

func LogSpanMax(request interface{}) interface{} {
	data, _ := json.Marshal(request)
	if len(data) > constants.MAXUDP {
		request = fmt.Sprint("%s", data[:constants.MAXUDP])
	}
	return request
}

func GetFormatUangInt(amount int64) string {
	// amt, err := strconv.Atoi(amount)
	// if err != nil {
	// 	return "0"
	// }
	ac := accounting.Accounting{Symbol: "Rp ", Precision: 0, Thousand: ",", Decimal: "."}
	return ac.FormatMoney(amount)
}

func GetTimeFormatResponse() string {
	t := time.Now().Local()
	ts := fmt.Sprintf("%s %02d:%02d", t.Format("02-Jan-2006"), t.Hour(), t.Minute())
	return ts
}

func GetRrn() string {
	//res, err := redis.GetRedisKey(rrnkey)
	counter, err := redis.SaveRedisCounter(rrnkey)
	if err != nil {
		counter = 1
	}
	t11 := time.Now().Local()
	return fmt.Sprintf("%02d%02d%02d%02d%04d", t11.Day(), t11.Hour(), t11.Minute(), t11.Second(), counter)
}

func GetTimeFormatYYMMDDHHMMSS() string {
	t11 := time.Now().Local()
	strthn := fmt.Sprintf("%v", t11.Year())
	return fmt.Sprintf("%s%02d%02d%02d%02d%02d", strthn[2:4], t11.Month(), t11.Day(), t11.Hour(), t11.Minute(), t11.Second())
}

func Operator(code int) string {
	var operator string
	switch code {
	case 0:
		operator = "Telkomsel"
		break
	case 1:
		operator = "Indosat"
		break
	case 2:
		operator = "XL"
		break
	case 3:
		operator = "Three"
		break
	case 4:
		operator = "Smartfren"
		break
	}

	return operator
}

func ProductPulsa(code string) string {
	var productCode string
	switch code {
	case "1080", "1250", "1251":
		productCode = "Telkomsel"
		break
	case "1081", "1255", "1254":
		productCode = "XL"
		break
	case "1082", "1253", "1252":
		productCode = "Indosat"
		break
		// case "1083": //108300
		// 	productCode = "smartfren"
		// 	break
		// case "1084": //108400
		// 	productCode = "three"
		// 	break
	}

	return productCode
}

func TypeProduct(code string) string {
	var productCode string
	switch code {
	case "108":
		productCode = "pulsa"
		break
	case "112":
		productCode = "paket data"
		break
	}

	return productCode
}
