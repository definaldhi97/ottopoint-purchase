package host

import (
	"encoding/json"
	"fmt"
	"log"
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
	endpointListRulePoint          string
	endpointGetBalance             string

	endpointAddedPoint string
	endpointSpendPoint string

	HealthCheckKey string
)

func init() {
	host = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_OPL", "http://18.138.173.105")
	name = ODU.GetEnv("OTTOPOINT_PURCHASE_NAME_OPL", "OPENLOYALTY")

	endpointRedeemVoucher = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_VOUCHER_REDEEM", "/api/customer/campaign/")
	endpointVoucherDetail = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_VOUCHER_DETAIL", "/api/campaign/")
	endpointHistoryVoucherCustomer = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_HISTORY_VOUCHER", "/api/customer/campaign/bought")
	endpointCouponVoucherCustomer = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_COUPONVOUCHER", "/api/admin/campaign/coupons/mark_as_used")

	endpointRulePoint = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_RULEPOINT", "/api/customer/earnRule/")
	endpointListRulePoint = ODU.GetEnv("OTTOPOINT_PURCHASE_LIST_RULEPOINT", "/api/customer/earningRule")

	endpointAddedPoint = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_ADD_POINT", "/api/points/transfer/add")
	endpointSpendPoint = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_SPEND_POINT", "/api/points/transfer/spend")

	endpointGetBalance = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_GET_BALANCE", "/api/admin/customer")

	HealthCheckKey = ODU.GetEnv("OTTOPOINT_PURCHASE_KEY_HEALTHCHECK_OPL", "OTTOPOINT-PURCHASE:OTTOPOINT")
}

// Redeem Voucher
func RedeemVoucher(campaignID, phone string) (*models.BuyVocuherResp, error) {
	var resp models.BuyVocuherResp

	logs.Info("[Package Host OPL]-[RedeemVoucher]")

	api := campaignID + "/buy"
	urlSvr := host + endpointRedeemVoucher + api

	data, err := HTTPxFormPostCustomer1(urlSvr, phone, HealthCheckKey)
	if err != nil {
		logs.Error("Check error", err.Error())

		return &resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logs.Error("Failed to unmarshaling response RedeemVoucher from open-loyalty ", err.Error())

		return &resp, err
	}

	return &resp, nil
}

// RulePoint
func RulePoint(eventName, phone string) (models.RulePointResponse, error) {
	var resp models.RulePointResponse

	logs.Info("[Package Host OPL]-[RulePoint]")

	todo := endpointRulePoint + eventName
	logs.Info("Request EranRule :", todo)

	urlSvr := host + todo

	data, err := HTTPxFormPostCustomer1(urlSvr, phone, HealthCheckKey)
	if err != nil {
		logs.Error("Check error Rule Point", err.Error())
		//fmt.Printf("Check error %v", err.Error())
		return resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logs.Error("Failed to unmarshaling response RulePoint from open-loyalty ", err.Error())

		return resp, err
	}

	return resp, nil

}

// ListRulePoint
func ListRulePoint(phone string) (models.LisrRulePointResponse, error) {
	var resp models.LisrRulePointResponse

	logs.Info("[Package Host OPL]-[ListRulePoint]")

	todo := endpointListRulePoint

	urlSvr := host + todo

	data, err := HTTPxFormGETCustomer(urlSvr, phone, HealthCheckKey)
	if err != nil {
		logs.Error("Check error Rule Point", err.Error())
		//fmt.Printf("Check error %v", err.Error())
		return resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logs.Error("Failed to unmarshaling response ListRulePoint from open-loyalty ", err.Error())

		return resp, err
	}

	return resp, nil

}

// Voucher Detail
func VoucherDetail(campaign string) (models.VoucherDetailResp, error) {
	var resp models.VoucherDetailResp

	logs.Info("[Package Host OPL]-[VoucherDetail]")

	urlSvr := host + endpointVoucherDetail + campaign

	data, err := HTTPxFormGETAdmin(urlSvr, HealthCheckKey)
	if err != nil {
		logs.Error("Check error Voucher Detail ", err.Error())
		//fmt.Printf("Check error %v", err.Error())
		return resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logs.Error("Failed to unmarshaling response VoucherDetail from open-loyalty ", err.Error())

		return resp, err
	}

	return resp, nil

}

// History Voucher Customer
func HistoryVoucherCustomer(phone, page string) (*models.HistoryVoucherCustomerResponse, error) {
	var resp models.HistoryVoucherCustomerResponse

	logs.Info("[Package Host OPL]-[HistoryVoucherCustomer]")

	param := fmt.Sprintf("?includeDetails=1&page=%s&perPage=1000&sort&direction", page)
	urlSvr := host + endpointHistoryVoucherCustomer + param
	data, err := HTTPxFormGETCustomer(urlSvr, phone, HealthCheckKey)
	if err != nil {
		logs.Error("Check error ", err.Error())
		return &resp, err
	}

	logs.Info("Response OPL")

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logs.Error("Failed to unmarshaling response HistoryVoucherCustomer from open-loyalty ", err.Error())
		return &resp, err
	}
	return &resp, nil
}

// CouponVoucherCustomer ..
func CouponVoucherCustomer(campaign, couponId, couponCode, custID string, useVoucher int) (*models.CouponVoucherCustomerResp, error) {
	var resp models.CouponVoucherCustomerResp

	logs.Info("[Package Host OPL]-[CouponVoucherCustomer]")

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
		logs.Error("Failed to unmarshaling response CouponVoucherCustomer from open-loyalty ", err.Error())
		return &resp, err
	}
	return &resp, nil
}

// Transfer Point ..
func TransferPoint(customer string, point string, text string) (*models.PointResponse, error) {
	var resp models.PointResponse

	logs.Info("[Package Host OPL]-[TransferPoint]")
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
		logs.Error("Failed to unmarshaling response TransferPoint from open-loyalty ", err.Error())
		return &resp, err
	}
	return &resp, nil
}

// Spend Point ..
func SpendPoint(customer, point, text string) (*models.PointResponse, error) {
	var resp models.PointResponse

	logs.Info("[Package Host OPL]-[SpendPoint]")

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
		logs.Error("Failed to unmarshaling response SpendPoint from open-loyalty ", err.Error())
		return &resp, err
	}
	return &resp, nil
}

// GetBalance
func GetBalance(customer string) (*models.BalanceResponse, error) {
	var result models.BalanceResponse

	logs.Info("[Package Host OPL]-[GetBalance]")

	cust := "/" + customer + "/status"
	urlSvr := host + endpointGetBalance
	log.Print("url endpoind status : ", urlSvr)
	data, err := HTTPxFormCustomerStatus(urlSvr, cust)
	if err != nil {
		logs.Error("CustomerStatus ", err.Error())
		return &result, err
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		logs.Error("Failed to unmarshaling response GetBalance from open-loyalty ", err.Error())

		return &result, err
	}
	return &result, nil
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
