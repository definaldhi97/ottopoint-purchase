package host

import (
	"encoding/json"
	"ottopoint-purchase/hosts/ultra_voucher/models"
	"ottopoint-purchase/redis"

	"github.com/astaxie/beego/logs"
	hcmodels "ottodigital.id/library/healthcheck/models"
	hcutils "ottodigital.id/library/healthcheck/utils"
	ODU "ottodigital.id/library/utils"
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
	host = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_UV", "http://13.228.25.85:8704/uv-service/v0.1.0")

	name = ODU.GetEnv("OTTOPOINT_PURCHASE_NAME_UV", "ULTRA VOUCHER")

	endpointOrderVoucher = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_UV_ORDER", "/purchase/order")
	endpointUseVoucher = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_UV_USE", "/voucher/use")
	endpointCheckStatusOrder = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_UV_CHECK_ORDER", "/check/status-order-voucher")

	// healthCheckKey = ODU.GetEnv("OTTOPOINT_PURCHASE_KEY_HEALTHCHECK_UV", "OTTOPOINT-PURCHASE:OTTOPOINT-UV")
}

// OrderVoucher
func OrderVoucher(req models.OrderVoucherReq, institutionID string) (*models.OrderVoucherResp, error) {
	var resp models.OrderVoucherResp

	logs.Info("[PackageHostUV]-[OrderVoucher]")

	urlSvr := host + endpointOrderVoucher

	data, err := HTTPxFormPostUV(urlSvr, institutionID, req)
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

	logs.Info("[Package Host UV]-[UseVoucher]")

	// req := models.UseVoucherUVReq{
	// 	Account:     accountNumber,
	// 	VoucherCode: code,
	// }

	urlSvr := host + endpointUseVoucher

	data, err := HTTPxFormPostUV(urlSvr, "", req)
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

	logs.Info("[Package Host UV]-[CheckStatusOrder]")

	// req := models.UseVoucherUVReq{
	// 	Account:     accountNumber,
	// 	VoucherCode: code,
	// }

	urlSvr := host + endpointCheckStatusOrder

	data, err := HTTPxFormGETUV(urlSvr, InstitutionReff, InstitutionId)
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

// GetServiceHealthCheck ..
func GetServiceHealthCheck() hcmodels.ServiceHealthCheck {
	redisClient := redis.GetRedisConnection()
	return hcutils.GetServiceHealthCheck(&redisClient, &hcmodels.ServiceEnv{
		Name:    name,
		Address: host,
		// HealthCheckKey: healthCheckKey,
	})
}
