package host

import (
	"encoding/json"
	"fmt"
	"ottopoint-purchase/hosts/sepulsa/models"

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

	ewalletInsertTransaction = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_EWALLET_INSERT_TRANSACTION", "/transaction/ewallet.json")
	ewalletDetailTransaction = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_EWALLET_DETAIL_TRANSACTION", "/transaction/ewallet/")

}

func EwalletInsertTransaction(req models.EwalletInsertTrxReq) (*models.EwalletInsertTrxRes, error) {
	var resp models.EwalletInsertTrxRes

	urlSvr := host + ewalletInsertTransaction
	data, err := HTTPxFormPostSepulsa(urlSvr, req)
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

	urlSvr := host + ewalletDetailTransaction + trxID + ".json"

	data, err := HTTPxFormGETSepulsa(urlSvr, nil)
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
