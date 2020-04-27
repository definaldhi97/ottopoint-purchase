package host

import (
	"fmt"
	"log"
	"strconv"
	"time"

	//"fmt"
	"net/http"
	"ottopoint-purchase/utils"

	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"

	hcmodels "ottodigital.id/library/healthcheck/models"
	ODU "ottodigital.id/library/utils"
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

	host = ODU.GetEnv("OTTOPOINT_PURCHASE_OTTOAG_HOST", "http://13.228.25.85:8089/")
	endpointInquiry = ODU.GetEnv("OTTOPOINT_PURCHASE_OTTOAG_ENDPOINT_INQUIRY", "v1/inquiry")
	endpointPayment = ODU.GetEnv("OTTOPOINT_PURCHASE_OTTOAG_ENDPOINT_PAYMENT", "v1/payment")
	authorization = ODU.GetEnv("OTTOPOINT_PURCHASE_OTTOAG_AUTHORIZATION", "T1RQT0lOVA==")
	serverkey = ODU.GetEnv("OTTOPOINT_PURCHASE_OTTOAG_SESSIONKEY", "052CFD8A04F99AC48E4656BBDF19FE60")
	HealthCheckKey = ODU.GetEnv("OTTOPOINT_PURCHASE_HEALTHCHECK_OTTOAG", "OTTOPOINT_HEALTH_CHECK:OTTOAG")
	Name = ODU.GetEnv("OTTOPOINT_PURCHASE_NAME_OTTOAG", "OTTOAG")
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
	data, err := HTTPPostWithHeader(urlSvr, msgreq, header)

	logs.Info(fmt.Sprintf("[Response %s]", typetrans), fmt.Sprintf("[%s]", string(data)))
	return data, err
}

// GetServiceHealthCheck ..
func GetServiceHealthCheckOttoAG() hcmodels.ServiceHealthCheck {
	res := hcmodels.ServiceHealthCheck{}
	var erorr interface{}
	// sugarLogger := service.General.OttoZapLog

	PublicAddress := host + endpointInquiry
	log.Print("url : ", PublicAddress)
	res.Name = Name
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
