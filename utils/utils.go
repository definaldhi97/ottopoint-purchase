package utils

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/models"
	"ottopoint-purchase/redis"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/leekchan/accounting"
	"github.com/vjeantet/jodaTime"

	ODU "ottodigital.id/library/utils"
)

var (
	rrnkey            string
	DefaultStatusCode int
	DefaultStatusMsg  string
	RedisKeyAuth      string
	LimitTRXPoint     string
	MemberID          string
	PathCSV           string

	ListErrorCode []models.MappingErrorCodes
)

func init() {
	DefaultStatusCode = 200
	DefaultStatusMsg = "OK"
	rrnkey = ODU.GetEnv("REDISKEY.OTTOFIN.RRN", "OTTOFIN:KEYRRN")
	RedisKeyAuth = ODU.GetEnv("redis.key.auth", "Ottopoint-Token-Admin :")
	LimitTRXPoint = ODU.GetEnv("limit.trx.point", "999999999999999")
	MemberID = ODU.GetEnv("OTTOPOINT_PURCHASE_OTTOAG_MEMBERID", "OTPOINT")
	// PathCSV = ODU.GetEnv("PATH_CSV", "//Users/abdulrohmat/Documents/Golang/src/ottopoint-purchase/utils/")
	PathCSV = ODU.GetEnv("PATH_CSV", "/opt/ottopoint-purchase/csv/")

}

func GetMessageResponse(res models.Response, code int, status bool, err error) models.Response {

	res = models.Response{}

	res.Meta.Code = code
	res.Meta.Status = status
	res.Meta.Message = err.Error()

	return res
}

func GetMessageFailedError(res models.Response, code int, err error) models.Response {

	res = models.Response{}

	res.Meta.Code = code
	res.Meta.Status = false
	res.Meta.Message = err.Error()

	return res
}

func GetMessageResponseData(res models.Response, resData models.ResponseData, code int, status bool, err error) models.Response {

	res = models.Response{}

	res.Data = resData
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

func ResponseMetaOK() models.MetaData {
	return models.MetaData{
		Status:  true,
		Code:    200,
		Message: "SUCCESS",
	}
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
	case "1081", "1255", "1254", "":
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

func ProductPaketData(code string) string {
	var productCode string
	switch code {
	case "12560":
		productCode = "Telkomsel"
		break
	case "12565":
		productCode = "XL"
		break
	case "12570":
		productCode = "Indosat"
		break
	case "12575": //108300
		productCode = "Three"
		break
	case "12580": //108400
		productCode = "Smartfren"
		break
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

// GetMetaResponse ..
func GetMetaResponse(key string) models.MetaData {
	fmt.Println("Get response by key:", key)

	var meta models.MetaData

	if key == constants.KeyResponseSucceed {
		meta.Code = 200
		meta.Message = "Successful"
		meta.Status = true
		return meta
	}

	for _, element := range ListErrorCode {
		if element.Key == key {
			meta.Status = false
			meta.Code = element.Content.Code
			meta.Message = element.Content.Message
			return meta
		}
	}

	meta.Code = 400
	meta.Message = "Terjadi kesalahan pada server"
	meta.Status = false

	return meta
}

// GetFormattedToken ...
func GetFormattedToken(token string) string {

	tokenLen := len(token)

	if tokenLen == 0 {
		return token
	}

	formattedToken := ""

	for index := 0; index < tokenLen; index++ {
		if (index%4) == 0 && index != 0 && index != tokenLen-1 {
			formattedToken = fmt.Sprintf("%s %s", formattedToken, token[index:index+1])

			continue
		}

		formattedToken = fmt.Sprintf("%s%s", formattedToken, token[index:index+1])
	}

	return formattedToken
}

// // generate token using UUID and base64
// func GenerateTokenUUID() string {
// 	out, err := exec.Command("uuidgen").Output()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("%s", out)
// 	tokenString := string(out)

// 	tokenString = strings.ReplaceAll(tokenString, "\n", "")

// 	encode64Token := base64.StdEncoding.EncodeToString([]byte(tokenString))
// 	log.Print(encode64Token)
// 	return encode64Token
// }

func GenerateTokenUUID() string {
	value := uuid.Must(uuid.NewRandom())
	fmt.Println("ini ID : ", value)
	out := value.String()
	fmt.Printf("%s", out)
	tokenString := string(out)

	tokenString = strings.ReplaceAll(tokenString, "\n", "")
	tokenString = strings.ToLower(tokenString)
	return tokenString
}

// ReffNumb
func GenTransactionId() string {

	currentTime := fmt.Sprintf(time.Now().Format("060102"))
	currentMilitmp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	currentMili := currentMilitmp[8:len(currentMilitmp)]
	randomvalue := strconv.Itoa(Random(11111, 99999))
	transactionID := currentTime + currentMili + randomvalue

	return transactionID
}

func Random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func GetMessageFailedErrorNew(res models.Response, resCode int, resDesc string) models.Response {
	res = models.Response{}
	res.Meta.Status = false
	res.Meta.Code = resCode
	res.Meta.Message = resDesc

	return res
}

func ValidateTimeActive(status, allTime bool, startAt, endAt time.Time) bool {

	if status == false {

		fmt.Println("=== Earning InActive ===")
		return false
	}

	if allTime == false {
		now := jodaTime.Format("dd-MM-YYYY", time.Now())
		start := jodaTime.Format("dd-MM-YYYY", startAt)
		end := jodaTime.Format("dd-MM-YYYY", endAt)

		// validate masa active earning
		if now == end || now == start {

			fmt.Println("=== Earning Kadaluarsa ===")
			return false
		}
	}

	return true
}

func FormatTimeString(timestamp time.Time, year, month, day int) string {

	t := timestamp.AddDate(year, month, day)

	// 2020-05-01T17:40:24+0700
	res := jodaTime.Format("YYYY-MM-dd", t)

	return res
}

func CreateCSVFile(data interface{}, name string) {
	fmt.Println(">>> createCSVFile <<<")

	time := jodaTime.Format("dd-MM-YYYY HH:mm:ss", time.Now())
	path := PathCSV + name
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		if file, err = os.Create(path); err != nil {
			return
		}
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	returnError := writer.Write([]string{""})
	returnError = writer.Write([]string{time})
	if returnError != nil {
		fmt.Println(returnError)
	}

	writer.Flush()

	datas, _ := json.Marshal(&data)

	dataString := strings.Fields(string(datas))

	for _, value := range dataString {
		_, err := file.WriteString(strings.TrimSpace(value))
		if err != nil { //exception handler
			fmt.Println(err)
			break
		}
	}

	writer.Flush()

}

func Before(value string, a string) string {
	// Get substring before a string.
	pos := strings.Index(value, a)
	if pos == -1 {
		return ""
	}
	return value[0:pos]
}

func After(value string, a string) string {
	// Get substring after a string.
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:len(value)]
}
