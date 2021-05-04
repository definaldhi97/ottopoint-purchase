package utils

import (
	"crypto/aes"
	"crypto/cipher"

	cryptRand "crypto/rand"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/models"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/leekchan/accounting"
	"github.com/vjeantet/jodaTime"
)

var (
	DefaultStatusCode int
	DefaultStatusMsg  string
	RedisKeyAuth      string
	LimitTRXPoint     string
	MemberID          string
	PathCSV           string

	// Topics
	TopicsEarning string
	TopicsNotif   string
	TopicNotifSMS string

	// URL
	UrlImage      string
	ListErrorCode []models.MappingErrorCodes
)

func init() {
	DefaultStatusCode = 200
	DefaultStatusMsg = "OK"
	RedisKeyAuth = GetEnv("redis.key.auth", "Ottopoint-Token-Admin :")
	LimitTRXPoint = GetEnv("limit.trx.point", "999999999999999")
	MemberID = GetEnv("OTTOPOINT_PURCHASE_OTTOAG_MEMBERID", "OTPOINT")
	PathCSV = GetEnv("PATH_CSV", "/opt/ottopoint-purchase/csv/")
	TopicsNotif = GetEnv("TOPICS_NOTIF", "ottopoint-notification-topics")
	TopicsEarning = GetEnv("TOPICS_EARNING", "ottopoint-earning-topics")
	TopicNotifSMS = GetEnv("TOPIC_NOTIF_SMS", "ottopoint-sms-notification-topics")
	UrlImage = GetEnv("URL_IMAGES", "https://apidev.ottopoint.id/product/v2.1/image/")

}

func GetEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
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

func GenerateUUID() string {
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
	randomString, _ := GenerateRandomString(12)
	transactionID := currentTime + randomString

	return transactionID
}

func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789"
	bytes, err := GenerateRandomBytes(n)
	if err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := cryptRand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
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

func EncryptAES(plaintext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(cryptRand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func DecryptAES(ciphertext []byte, key []byte) ([]byte, error) {

	fmt.Println("Prosess decrytp ")
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	fmt.Println("value nonceSize Decrypt : ", nonceSize)
	fmt.Println("value len ciphertext Decrypt : ", len(ciphertext))
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

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

func ExpireDateVoucherAGt(value int) time.Time {
	now := time.Now().Local()
	layOut := "2006-01-02"
	expierDate := now.AddDate(0, 0, value)
	date := expierDate.Format("2006-01-02")
	dateStamp, _ := time.Parse(layOut, date)
	return dateStamp
}

func DefaultNulTime(date time.Time) *time.Time {
	if !date.IsZero() {
		return &date
	}
	return nil
}

func GetTimeFormatMillisecond() string {
	now := time.Now().Local()
	unixNano := now.UnixNano()
	umillisec := unixNano / 1000000
	// fmt.Println("(correct)Millisecond : ", umillisec)
	convString := strconv.FormatInt(umillisec, 10)
	return convString

}

func EncryptVoucherCode(data, key string) string {

	var codeVoucher string
	if data == "" {
		return codeVoucher
	}

	a := []rune(key)
	key32 := string(a[0:32])
	screetKey := []byte(key32)
	codeByte := []byte(data)
	chiperText, _ := EncryptAES(codeByte, screetKey)
	codeVoucher = string(chiperText)
	return codeVoucher
}
