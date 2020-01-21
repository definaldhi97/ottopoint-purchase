package host

import (
	"encoding/json"
	"fmt"
	"ottopoint-purchase/hosts/opl/models"
	"ottopoint-purchase/redis"

	"github.com/astaxie/beego/logs"
	hcmodels "ottodigital.id/library/healthcheck/models"
	hcutils "ottodigital.id/library/healthcheck/utils"
	ODU "ottodigital.id/library/utils"
)

var (
	host string
	name string

	endpointVoucherDetail          string
	endpointRedeemVoucher          string
	endpointCouponVoucherCustomer  string
	endpointHistoryVoucherCustomer string
	endpointRulePoint              string
	endpointAddedPoint             string
	endpointSpendPoint             string

	HealthCheckKey string
)

func init() {
	host = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_OPL", "http://54.179.186.194")
	name = ODU.GetEnv("OTTOPOINT_PURCHASE_NAME_OPL", "OPENLOYALTY")

	endpointRedeemVoucher = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_VOUCHER_REDEEM", "/api/customer/campaign/")
	endpointVoucherDetail = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_VOUCHER_DETAIL", "/api/campaign/")
	endpointHistoryVoucherCustomer = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_HISTORY_VOUCHER", "/api/customer/campaign/bought")
	endpointCouponVoucherCustomer = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_COUPONVOUCHER", "/api/admin/campaign/coupons/mark_as_used")
	endpointRulePoint = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_RULEPOINT", "/api/customer/earnRule/")
	endpointAddedPoint = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_ADD_POINT", "/api/points/transfer/add")
	endpointSpendPoint = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_SPEND_POINT", "/api/points/transfer/spend")

	HealthCheckKey = ODU.GetEnv("OTTOPOINT_PURCHASE_KEY_HEALTHCHECK_OPL", "OTTOPOINT-PURCHASE:OTTOPOINT")
}

// Redeem Voucher
func RedeemVoucher(campaignID, phone string) (*models.BuyVocuherResp, error) {
	var resp models.BuyVocuherResp

	api := campaignID + "/buy"
	urlSvr := host + endpointRedeemVoucher + api

	data, err := HTTPxFormPostCustomer1(urlSvr, phone, HealthCheckKey)
	if err != nil {
		logs.Error("Check error", err.Error())

		return &resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logs.Error("Failed to unmarshaling response from open-loyalty ", err.Error())

		return &resp, err
	}

	return &resp, nil
}

// RulePoint
func RulePoint(eventName, phone string) (models.RulePointResponse, error) {
	var resp models.RulePointResponse

	todo := endpointRulePoint + eventName
	logs.Info("Request EranRule :", todo)

	logs.Info("==========")
	logs.Info("Phone :", phone)
	logs.Info("==========")

	urlSvr := host + todo

	data, err := HTTPxFormPostCustomer1(urlSvr, phone, HealthCheckKey)
	if err != nil {
		logs.Error("Check error Rule Point", err.Error())
		//fmt.Printf("Check error %v", err.Error())
		return resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logs.Error("Failed to unmarshaling response from open-loyalty ", err.Error())

		return resp, err
	}

	return resp, nil

}

// Voucher Detail
func VoucherDetail(campaign string) (*models.VoucherDetailResp, error) {
	var resp models.VoucherDetailResp

	urlSvr := host + endpointVoucherDetail + campaign

	data, err := HTTPxFormGETAdmin(urlSvr, HealthCheckKey)
	if err != nil {
		logs.Error("Check error Voucher Detail ", err.Error())
		//fmt.Printf("Check error %v", err.Error())
		return &resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logs.Error("Failed to unmarshaling response from open-loyalty ", err.Error())

		return &resp, err
	}

	return &resp, nil

}

// Voucher Detail
func VoucherDetail2(campaign string) ([]models.VoucherDetailResp, error) {
	var resp []models.VoucherDetailResp

	urlSvr := host + endpointVoucherDetail + campaign

	data, err := HTTPxFormGETAdmin(urlSvr, HealthCheckKey)
	if err != nil {
		logs.Error("Check error Voucher Detail ", err.Error())
		//fmt.Printf("Check error %v", err.Error())
		return resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logs.Error("Failed to unmarshaling response from open-loyalty ", err.Error())

		return resp, err
	}

	return resp, nil

}

// History Voucher Customer
func HistoryVoucherCustomer(phone, page string) (*models.HistoryVoucherCustomerResponse, error) {
	var resp models.HistoryVoucherCustomerResponse

	param := fmt.Sprintf("?includeDetails=1&page=%s&perPage=100&sort&direction", page)
	urlSvr := host + endpointHistoryVoucherCustomer + param
	data, err := HTTPxFormGETCustomer(urlSvr, phone, HealthCheckKey)
	if err != nil {
		logs.Error("Check error ", err.Error())
		return &resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logs.Error("Failed to unmarshaling response from open-loyalty ", err.Error())
		return &resp, err
	}
	return &resp, nil
}

// CouponVoucherCustomer ..
func CouponVoucherCustomer(campaign, couponId, couponCode, custID string, useVoucher int) (*models.CouponVoucherCustomerResp, error) {

	var resp models.CouponVoucherCustomerResp
	urlSvr := host + endpointCouponVoucherCustomer

	jsonData := map[string]interface{}{
		"coupons[0][campaignId]": campaign,   //"coupons[0][campaignId]": campaign,
		"coupons[0][couponId]":   couponId,   //"coupons[0][couponId]":   couponId,
		"coupons[0][code]":       couponCode, //"coupons[0][code]":       couponCode,
		"coupons[0][used]":       useVoucher, //"coupons[0][used]":       "true"}
		"coupons[0][customerId]": custID}

	logs.Info("===== Use Voucher True / False =====")
	data, err := HTTPxFormPostAdmin2(urlSvr, jsonData, HealthCheckKey)
	if err != nil {
		logs.Error("Check error ", err.Error())
		return &resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logs.Error("Failed to unmarshaling response from open-loyalty ", err.Error())
		return &resp, err
	}
	return &resp, nil
}

// GetServiceHealthCheck ..
func GetServiceHealthCheck() hcmodels.ServiceHealthCheck {
	redisClient := redis.GetRedisConnection()
	return hcutils.GetServiceHealthCheck(&redisClient, &hcmodels.ServiceEnv{
		Name:           name,
		Address:        host,
		HealthCheckKey: HealthCheckKey,
	})
}

// Transfer Point ..
func TransferPoint(customer string, point string, text string) (*models.PointResponse, error) {

	var resp models.PointResponse
	urlSvr := host + endpointAddedPoint

	jsonData := map[string]interface{}{
		"transfer[customer]": customer,
		"transfer[points]":   point,
		"transfer[comment]":  text,
	}

	logs.Info("Request to OPL : ", jsonData)
	data, err := HTTPxFormPostAdmin2(urlSvr, jsonData, HealthCheckKey)
	if err != nil {
		logs.Error("Check error ", err.Error())
		return &resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logs.Error("Failed to unmarshaling response from open-loyalty ", err.Error())
		return &resp, err
	}
	return &resp, nil
}

// Spend Point ..
func SpendPoint(customer, point, text string) (*models.PointResponse, error) {

	var resp models.PointResponse
	urlSvr := host + endpointSpendPoint

	jsonData := map[string]interface{}{
		"transfer[customer]": customer,
		"transfer[points]":   point,
		"transfer[comment]":  text,
	}

	logs.Info("Request to OPL : ", jsonData)
	data, err := HTTPxFormPostAdmin2(urlSvr, jsonData, HealthCheckKey)
	if err != nil {
		logs.Error("Check error ", err.Error())
		return &resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logs.Error("Failed to unmarshaling response from open-loyalty ", err.Error())
		return &resp, err
	}
	return &resp, nil
}
