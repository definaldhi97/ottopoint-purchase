package host

import (
	"encoding/json"
	"ottopoint-purchase/hosts/sepulsa/models"

	"github.com/astaxie/beego/logs"
	ODU "ottodigital.id/library/utils"
)

var (
	host string
	name string

	ewalletInsertTransaction string
	ewalletDetailTransaction string
)

func init() {
	host = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_SEPULSA", "https://horven-api.sumpahpalapa.com/api")
	name = ODU.GetEnv("OTTOPOINT_PURCHASE_NAME_SEPULSA", "SEPULSA")

	ewalletInsertTransaction = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_EWALLET_INSERT_TRANSACTION", "/transaction/ewallet")
	ewalletDetailTransaction = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_EWALLET_DETAIL_TRANSACTION", "/transaction/ewallet/")

}

func EwalletInsertTransaction(req models.EwalletInsertTrxReq) (*models.EwalletInsertTrxRes, error) {
	var resp models.EwalletInsertTrxRes

	urlSvr := host + ewalletInsertTransaction

	data, err := HTTPxFormPostSepulsa(urlSvr, req)
	if err != nil {
		logs.Error("Check error : ", err.Error())
		return nil, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logs.Error("Failed to unmarshalling response EwalletInsertTransaction from Sepulsa", err.Error())
		return nil, err
	}
	return &resp, nil
}

func EwalletDetailTransaction(trxID string) (map[string]interface{}, error) {
	var resp map[string]interface{}

	urlSvr := host + ewalletDetailTransaction + trxID

	data, err := HTTPxFormGETSepulsa(urlSvr, nil)
	if err != nil {
		logs.Error("Check error : ", err.Error())
		return nil, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logs.Error("Failed to unmarshalling response EwalletDetailTransaction from Sepulsa", err.Error())
		return nil, err
	}
	return resp, nil
}
