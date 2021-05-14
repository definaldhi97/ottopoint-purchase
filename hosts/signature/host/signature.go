package host

import (
	"encoding/json"
	"net/http"
	"ottopoint-purchase/constants"
	https "ottopoint-purchase/hosts"
	"ottopoint-purchase/hosts/signature/models"
	headermodels "ottopoint-purchase/models"
	"ottopoint-purchase/utils"

	"github.com/astaxie/beego/logs"
	"github.com/sirupsen/logrus"
)

var (
	host              string
	name              string
	endpointSignature string

	HealthCheckKey string
)

func init() {
	host = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_SIGNATURE", "http://13.228.25.85:8666")
	name = utils.GetEnv("OTTOPOINT_PURCHASE_NAME_SIGNATURE", "SIGNATURE")

	endpointSignature = utils.GetEnv("OTTOPOINT_PURCHASE_ENDPOINT_VALIDATE_SIGNATURE", "/auth/v2/signature")

	HealthCheckKey = utils.GetEnv("OTTOPOINT_PURCHASE_KEY_HEALTHCHECK_SIGNATURE", "OTTOPOINT-PURCHASE:SIGNATURE")
}

// Signature
func Signature(signature interface{}, headers headermodels.RequestHeader) (*models.SignatureResp, error) {
	var resp models.SignatureResp

	logrus.Info("[Hit to API Signature]")

	urlSvr := host + endpointSignature

	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	header.Set("DeviceId", headers.DeviceID)
	header.Set("InstitutionId", headers.InstitutionID)
	header.Set("Geolocation", headers.Geolocation)
	header.Set("ChannelId", headers.ChannelID)
	header.Set("AppsId", headers.AppsID)
	header.Set("Timestamp", headers.Timestamp)
	header.Set("Signature", headers.Signature)

	if headers.ChannelID == constants.SDK_WEB {
		header.Set("Authorization", headers.Authorization)
	}

	data, err := https.HTTPxPOSTwithRequest(urlSvr, signature, header)
	// data, err := HTTPxFormPostWithHeader(urlSvr, HealthCheckKey, signature, header)
	if err != nil {
		logs.Error("Check error", err.Error())

		return &resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logs.Error("Failed to unmarshaling response from Signature API ", err.Error())

		return &resp, err
	}

	return &resp, nil
}
