package host

import (
	"encoding/json"
	"ottopoint-purchase/models"

	ottomart "ottopoint-purchase/hosts/ottomart/models"

	"github.com/astaxie/beego/logs"

	ODU "ottodigital.id/library/utils"
)

var (
	host          string
	endpointToken string
	endpointNotif string
)

func init() {

	host = ODU.GetEnv("OTTOMART_HOST", "http://13.228.25.85:8999")
	endpointToken = ODU.GetEnv("OTTOMART_ENDPOINT_TOKEN", "/ottopay/v0.1.0/auth_token")
	endpointNotif = ODU.GetEnv("OTTOMART_ENDPOINT_NOTIF", "/ottopay/v0.1.0/ottopoint/notif")

}

func CheckToken(header models.RequestHeader) (ottomart.ResponseToken, error) {
	var resp ottomart.ResponseToken

	urlSvr := host + endpointToken
	data, err := HTTPxFormTokenOttomart(urlSvr, header)
	if err != nil {
		logs.Error("MerchantCheckPhone ", err.Error())
		//fmt.Printf("Check error %v", err.Error())
		return resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logs.Error("Failed to unmarshaling response from open-loyalty ", err.Error())
		//fmt.Printf("Failed to unmarshaling response from open-loyalty %v", err)

		return resp, err
	}

	logs.Info("========== Phone Ottomart", resp.Data.AccountNumber)

	return resp, nil
}

type Notif struct {
	Phone string `json:"phone"`
	Point int    `json:"point"`
	Rc    string `json:"rc"`
}

func NotifInboxOttomart(notif ottomart.NotifReq, header models.RequestHeader) (*ottomart.NotifResp, error) {
	var resp ottomart.NotifResp

	urlSvr := host + endpointNotif
	jsonData := notif
	logs.Info("===== Response Ottomart ======")
	data, err := HTTPxFormOttomartNotif(urlSvr, jsonData, header)
	if err != nil {
		logs.Error("MerchantCheckPhone ", err.Error())
		//fmt.Printf("Check error %v", err.Error())
		return &resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logs.Error("Failed to unmarshaling response from open-loyalty ", err.Error())
		//fmt.Printf("Failed to unmarshaling response from open-loyalty %v", err)

		return &resp, err
	}

	return &resp, nil
}
