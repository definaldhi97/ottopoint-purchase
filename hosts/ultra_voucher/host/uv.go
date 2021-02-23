package host

import (
	"encoding/json"
	"net/http"
	"ottopoint-purchase/hosts/ultra_voucher/models"

	https "ottopoint-purchase/hosts"
	"ottopoint-purchase/utils"

	"github.com/astaxie/beego/logs"
	"github.com/sirupsen/logrus"
)

var (
	host string
	name string

	endpointOrderVoucher     string
	endpointUseVoucher       string
	endpointCheckStatusOrder string

	// healthCheckKey string
)

func init() {
	host = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_UV", "http://13.228.25.85:8704/uv-service/v0.1.0")

	name = utils.GetEnv("OTTOPOINT_PURCHASE_NAME_UV", "ULTRA VOUCHER")

	endpointOrderVoucher = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_UV_ORDER", "/purchase/order")
	endpointUseVoucher = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_UV_USE", "/voucher/use")
	endpointCheckStatusOrder = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_UV_CHECK_ORDER", "/check/status-order-voucher")

	// healthCheckKey = utils.GetEnv("OTTOPOINT_PURCHASE_KEY_HEALTHCHECK_UV", "OTTOPOINT-PURCHASE:OTTOPOINT-UV")
}

// OrderVoucher
func OrderVoucher(req models.OrderVoucherReq, institutionID string) (*models.OrderVoucherResp, error) {
	var resp models.OrderVoucherResp

	logrus.Info("[PackageHostUV]-[OrderVoucher]")

	urlSvr := host + endpointOrderVoucher

	header := make(http.Header)

	if institutionID == "" {
		header.Set("Content-Type", "application/json")
	} else {
		header.Set("Content-Type", "application/json")
		header.Set("InstitutionId", institutionID)
	}

	data, err := https.HTTPxPOSTwithRequest(urlSvr, req, header)
	// data, err := HTTPxFormPostUV(urlSvr, institutionID, req)
	if err != nil {
		logs.Error("Check error : ", err.Error())

		return &resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logs.Error("Failed to unmarshaling response OrderVoucher from Ultra Voucher ", err.Error())

		return &resp, err
	}

	return &resp, nil
}

// UseVoucher
func UseVoucherUV(req models.UseVoucherUVReq) (*models.UseVoucherUVResp, error) {
	var resp models.UseVoucherUVResp

	logrus.Info("[Package Host UV]-[UseVoucher]")

	urlSvr := host + endpointUseVoucher

	header := make(http.Header)
	header.Set("Content-Type", "application/json")

	data, err := https.HTTPxPOSTwithRequest(urlSvr, req, header)
	// data, err := HTTPxFormPostUV(urlSvr, "", req)
	if err != nil {
		logs.Error("Check error : ", err.Error())

		return &resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logs.Error("Failed to unmarshaling response OrderVoucher from Ultra Voucher ", err.Error())

		return &resp, err
	}

	return &resp, nil
}

// CheckStatusOrder
func CheckStatusOrder(InstitutionReff, InstitutionId string) (models.OrderVoucherResp, error) {
	var resp models.OrderVoucherResp

	logrus.Info("[Package Host UV]-[CheckStatusOrder]")

	urlSvr := host + endpointCheckStatusOrder

	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	header.Set("InstitutionId", InstitutionId)
	header.Set("InstitutionRefno", InstitutionReff)

	data, err := https.HTTPxGET(urlSvr, header)
	// data, err := HTTPxFormGETUV(urlSvr, InstitutionReff, InstitutionId)
	if err != nil {
		logs.Error("Check error : ", err.Error())

		return resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logs.Error("Failed to unmarshaling response CheckStatusOrder from Ultra Voucher ", err.Error())

		return resp, err
	}

	return resp, nil
}
