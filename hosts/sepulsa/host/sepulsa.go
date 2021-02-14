package host

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ottopoint-purchase/hosts/sepulsa/models"

	https "ottopoint-purchase/hosts"
	"ottopoint-purchase/utils"
)

var (
	host string
	name string

	ewalletInsertTransaction string
	ewalletDetailTransaction string
)

func init() {
	host = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_SEPULSA", "https://horven-api.sumpahpalapa.com/api")
	name = utils.GetEnv("OTTOPOINT_PURCHASE_NAME_SEPULSA", "SEPULSA")

	ewalletInsertTransaction = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_EWALLET_INSERT_TRANSACTION", "/transaction/ewallet.json")
	ewalletDetailTransaction = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_EWALLET_DETAIL_TRANSACTION", "/transaction/ewallet/")

}

func EwalletInsertTransaction(req models.EwalletInsertTrxReq) (*models.EwalletInsertTrxRes, error) {
	var resp models.EwalletInsertTrxRes

	header := make(http.Header)
	header.Set("User-Agent", "ottopoint")
	header.Set("Accept", "application/json")
	header.Set("Content-Type", "application/json")

	urlSvr := host + ewalletInsertTransaction
	data, err := https.HTTPxPOSTxSepulsa(urlSvr, req, header)
	// data, err := HTTPxFormPostSepulsa(urlSvr, req)
	if err != nil {

		fmt.Println(fmt.Sprintf("[Failed to HTTPxFormPostSepulsa]-[Error : %v]", err.Error()))
		fmt.Println("[PackageHostSepulsa]-[EwalletInsertTransaction]")

		return nil, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {

		fmt.Println(fmt.Sprintf("[Failed to unmarshalling response EwalletInsertTransaction]-[Error : %v]", err.Error()))
		fmt.Println("[PackageHostSepulsa]-[EwalletInsertTransaction]")

		return nil, err
	}
	return &resp, nil
}

func EwalletDetailTransaction(trxID string) (*models.CheckStatusSepulsaResp, error) {
	var resp models.CheckStatusSepulsaResp

	header := make(http.Header)
	header.Set("User-Agent", "ottopoint")
	header.Set("Accept", "application/json")
	header.Set("Content-Type", "application/json")

	urlSvr := host + ewalletDetailTransaction + trxID + ".json"

	data, err := https.HTTPxGETxSepulsa(urlSvr, header)
	// data, err := HTTPxFormGETSepulsa(urlSvr, nil)
	if err != nil {

		fmt.Println(fmt.Sprintf("[Failed to HTTPxFormPostSepulsa]-[Error : %v]", err.Error()))
		fmt.Println("[PackageHostSepulsa]-[EwalletDetailTransaction]")

		return nil, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {

		fmt.Println(fmt.Sprintf("[Failed to unmarshalling response EwalletDetailTransaction]-[Error : %v]", err.Error()))
		fmt.Println("[PackageHostSepulsa]-[EwalletDetailTransaction]")

		return nil, err
	}
	return &resp, nil
}
