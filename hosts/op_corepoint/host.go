package op_corepoint

import (
	"encoding/json"
	"fmt"
	"net/http"
	https "ottopoint-purchase/hosts"
	"ottopoint-purchase/models"

	"ottopoint-purchase/utils"

	"github.com/sirupsen/logrus"
)

var (
	host               string
	endpointAdding     string
	endpointSpending   string
	endpointGetBalance string
)

func init() {
	host = utils.GetEnv("HOST_OTTOPOINT_COREPOINT", "http://13.228.25.85:8402")
	endpointAdding = utils.GetEnv("ENDPOINT_ADDING_OTTOPOINT_COREPOINT", "/v1/points/transfer/add")
	endpointSpending = utils.GetEnv("ENDPOINT_SEPENDING_OTTOPOINT_COREPOINT", "/v1/points/transfer/spend")
	endpointGetBalance = utils.GetEnv("ENDPOINT_SEPENDING_OTTOPOINT_COREPOINT", "/v1/account/wallet/balance")
}

func AddingPoint(req AddingPointReq, headerReq models.RequestHeader) (*TrxPointRes, error) {
	fmt.Println("[ >>>>>>>>>>>>>> package Trx AddingPoint to ottopoint-corepoint <<<<<<<<<<<<<<<< ]")
	var result TrxPointRes
	urlSvr := host + endpointAdding

	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	header.Set("InstitutionId", headerReq.InstitutionID)
	header.Set("Authorization", headerReq.Authorization)
	header.Set("Deviceid", headerReq.DeviceID)
	header.Set("Geolocation", headerReq.Geolocation)
	header.Set("ChannelId", headerReq.ChannelID)
	header.Set("AppsId", headerReq.AppsID)

	data, err := https.HTTPxPOSTwithRequest(urlSvr, req, header)
	logrus.Info("Response Trx Adding Point to wallet ottopoint-corepoint : ", data)

	if err != nil {
		logrus.Error("Failed Trx Adding Point : ", err.Error())
		return &result, err
	}

	err = json.Unmarshal(data, &result)
	if err != nil {

		logrus.Error("Failed to unmarshaling response adding pointt : ", err.Error())
		return &result, err
	}

	return &result, err

}

func SependingPoint(req SpendingPointReq, headerReq models.RequestHeader) (*TrxPointRes, error) {
	fmt.Println("[ >>>>>>>>>>>>>> package Trx SpendingPoint to ottopoint-corepoint <<<<<<<<<<<<<<<< ]")

	var result TrxPointRes
	urlSvr := host + endpointSpending

	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	header.Set("InstitutionId", headerReq.InstitutionID)
	header.Set("Authorization", headerReq.Authorization)
	header.Set("Deviceid", headerReq.DeviceID)
	header.Set("Geolocation", headerReq.Geolocation)
	header.Set("ChannelId", headerReq.ChannelID)
	header.Set("AppsId", headerReq.AppsID)

	data, err := https.HTTPxPOSTwithRequest(urlSvr, req, header)
	logrus.Info("Response Trx Spending Point to wallet ottopoint-corepoint : ", data)

	if err != nil {
		logrus.Error("Failed Trx Spending Point : ", err.Error())
		return &result, err
	}

	err = json.Unmarshal(data, &result)
	if err != nil {

		logrus.Error("Failed to unmarshaling response Spending pointt : ", err.Error())
		return &result, err
	}

	return &result, err
}

func GetBalancePoint(accountID string) (*GetBalanceResponse, error) {

	var result GetBalanceResponse

	logrus.Info("CustID : ", accountID)

	accountId := "?accountId=" + accountID
	urlSvr := host + endpointGetBalance + accountId

	header := make(http.Header)

	data, err := https.HTTPxGET(urlSvr, header)
	logrus.Error("Response GetBalancePoint : ", data)

	if err != nil {
		logrus.Error("Failed GetBalancePoint : ", err.Error())
		return &result, err
	}

	err = json.Unmarshal(data, &result)
	if err != nil {

		logrus.Error("Failed to unmarshaling response GetBalancePoint : ", err.Error())
		return &result, err
	}

	return &result, err
}
