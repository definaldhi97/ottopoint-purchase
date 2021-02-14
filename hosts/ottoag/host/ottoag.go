package host

import (
	"fmt"
	"strconv"
	"time"

	//"fmt"
	"net/http"
	https "ottopoint-purchase/hosts"
	"ottopoint-purchase/utils"

	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"
)

var (
	host            string
	endpointInquiry string
	endpointPayment string
	authorization   string

	HealthCheckKey string
	Name           string

	serverkey string
	memberID  string
)

// HeaderHTTP ..
type HeaderHTTP struct {
	Signature string
	Timestamp string
}

func init() {
	// http://13.228.25.85:8089/ottoaggo/biller/v1.0.0/inquiry

	host = utils.GetEnv("OTTOPOINT_PURCHASE_OTTOAG_HOST", "http://13.228.25.85:8089/")
	endpointInquiry = utils.GetEnv("OTTOPOINT_PURCHASE_OTTOAG_ENDPOINT_INQUIRY", "v1/inquiry")
	endpointPayment = utils.GetEnv("OTTOPOINT_PURCHASE_OTTOAG_ENDPOINT_PAYMENT", "v1/payment")
	authorization = utils.GetEnv("OTTOPOINT_PURCHASE_OTTOAG_AUTHORIZATION", "T1RQT0lOVA==")              // dev
	serverkey = utils.GetEnv("OTTOPOINT_PURCHASE_OTTOAG_SESSIONKEY", "052CFD8A04F99AC48E4656BBDF19FE60") // dev
	HealthCheckKey = utils.GetEnv("OTTOPOINT_PURCHASE_HEALTHCHECK_OTTOAG", "OTTOPOINT_HEALTH_CHECK:OTTOAG")
	Name = utils.GetEnv("OTTOPOINT_PURCHASE_NAME_OTTOAG", "OTTOAG")
}

// PackMessageHeader ..
func PackMessageHeader(req interface{}) HeaderHTTP {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature := utils.OttoAGCreateSignature(timestamp, req, serverkey)

	headhttp := HeaderHTTP{
		Timestamp: timestamp,
		Signature: signature,
		//Auth: auth ,
	}

	return headhttp
}

// Send ..
func Send(msgreq interface{}, head HeaderHTTP, typetrans string) ([]byte, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	// authorization := "Basic " + base64.StdEncoding.EncodeToString([]byte(instuition))

	header := make(http.Header)
	header.Add("Accept", "*/*")
	header.Add("Content-Type", "application/json")
	header.Add("Authorization", fmt.Sprintf("Basic %v", authorization))
	// header.Add("Authorization", "Basic Q0xBUEFQRQ==")
	header.Add("Signature", head.Signature)
	header.Add("Timestamp", head.Timestamp)
	urlSvr := ""
	fmt.Println("header 1", header)
	switch typetrans {
	case "INQUIRY":
		fmt.Println("Inquiry")
		urlSvr = host + endpointInquiry
		break
	case "PAYMENT":
		fmt.Println("Payment")
		urlSvr = host + endpointPayment
		break
	}

	datareq, _ := json.Marshal(msgreq)
	logs.Info(fmt.Sprintf("[Request %s]", typetrans), fmt.Sprintf("[%s]", string(datareq)))
	data, err := https.HTTPxPOSTwithRequest(urlSvr, msgreq, header)

	logs.Info(fmt.Sprintf("[Response %s]", typetrans), fmt.Sprintf("[%s]", string(data)))
	return data, err
}
