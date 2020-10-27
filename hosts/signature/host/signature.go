package host

import (
	"encoding/json"
	"ottopoint-purchase/hosts/signature/models"
	headermodels "ottopoint-purchase/models"
	"time"

	"github.com/astaxie/beego/logs"
	hcmodels "ottodigital.id/library/healthcheck/models"
	ODU "ottodigital.id/library/utils"
)

var (
	host              string
	name              string
	endpointSignature string

	HealthCheckKey string
)

func init() {
	host = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_SIGNATURE", "http://13.228.25.85:8666")
	name = ODU.GetEnv("OTTOPOINT_PURCHASE_NAME_SIGNATURE", "SIGNATURE")

	endpointSignature = ODU.GetEnv("OTTOPOINT_PURCHASE_ENDPOINT_VALIDATE_SIGNATURE", "/auth/v2/signature")

	HealthCheckKey = ODU.GetEnv("OTTOPOINT_PURCHASE_KEY_HEALTHCHECK_SIGNATURE", "OTTOPOINT-PURCHASE:SIGNATURE")
}

// Signature
func Signature(signature interface{}, header headermodels.RequestHeader) (*models.SignatureResp, error) {
	var resp models.SignatureResp

	logs.Info("[Hit to API Signature]")

	urlSvr := host + endpointSignature

	data, err := HTTPxFormPostWithHeader(urlSvr, HealthCheckKey, signature, header)
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

// GetServiceHealthCheck ..
func GetServiceHealthCheckSignature() hcmodels.ServiceHealthCheck {
	return hcmodels.ServiceHealthCheck{
		Name:    name,
		Address: host,
		Status:  "OK",
		// Description: ,
		UpdatedAt: time.Now(),
	}
}
