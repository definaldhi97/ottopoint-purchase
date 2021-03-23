package host

import (
	"encoding/json"
	"fmt"
	"log"
	"ottopoint-purchase/hosts/opl/models"
	"ottopoint-purchase/utils"

	"github.com/astaxie/beego/logs"
	"github.com/sirupsen/logrus"
)

var (
	host string
	name string

	endpointVoucherDetail            string
	endpointRedeemVoucher            string
	endpointCouponVoucherCustomer    string
	endpointHistoryVoucherCustomer   string
	endpointRulePoint                string
	endpointListRulePoint            string
	endpointGetBalance               string
	endpointRedeemCumulativeVoucher  string
	endpointRedeemCumulativeVoucher2 string

	endpointAddedPoint string
	endpointSpendPoint string

	endpointSetting string

	// HealthCheckKey string
)

func init() {
	host = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_OPL", "https://openloyalty-stg.ottopoint.id")
	name = utils.GetEnv("OTTOPOINT_PURCHASE_NAME_OPL", "OPENLOYALTY")

	endpointRedeemVoucher = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_VOUCHER_REDEEM", "/api/customer/campaign/")
	endpointRedeemCumulativeVoucher2 = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_VOUCHER_REDEEM_CUMULATIVE", "/api/customer/campaign/")
	endpointRedeemCumulativeVoucher = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_VOUCHER_REDEEM_CUMULATIVE_v2", "/api/admin/customer/")

	endpointVoucherDetail = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_VOUCHER_DETAIL", "/api/campaign/")
	endpointHistoryVoucherCustomer = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_HISTORY_VOUCHER", "/api/customer/campaign/bought")
	endpointCouponVoucherCustomer = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_COUPONVOUCHER", "/api/admin/campaign/coupons/mark_as_used")

	endpointRulePoint = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_RULEPOINT", "/api/customer/earnRule/")
	endpointListRulePoint = utils.GetEnv("OTTOPOINT_PURCHASE_LIST_RULEPOINT", "/api/customer/earningRule")

	endpointAddedPoint = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_ADD_POINT", "/api/points/transfer/add")
	endpointSpendPoint = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_SPEND_POINT", "/api/points/transfer/spend")

	endpointSetting = utils.GetEnv("OTTOPOINT_PURCHASE_SETTING_OPL", "/api/settings")

	endpointGetBalance = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_GET_BALANCE", "/api/admin/customer")

	// HealthCheckKey = utils.GetEnv("OTTOPOINT_PURCHASE_KEY_HEALTHCHECK_OPL", "OTTOPOINT-PURCHASE:OTTOPOINT")
}

// Redeem Voucher
func RedeemVoucher(campaignID, phone string) (*models.BuyVocuherResp, error) {
	var resp models.BuyVocuherResp

	logrus.Info("[Package Host OPL]-[RedeemVoucher]")

	api := campaignID + "/buy"
	urlSvr := host + endpointRedeemVoucher + api

	data, err := HTTPxFormPostCustomerWithoutRequest(urlSvr, phone)
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

// Redeem Voucher Cumulative
func RedeemVoucherCumulative(campaignID, custId, total, status string) (*models.BuyVocuherResp, error) {
	var resp models.BuyVocuherResp

	logrus.Info("[Package Host OPL]-[RedeemVoucherCumulative]")

	jsonData := map[string]interface{}{
		"quantity":      total,
		"withoutPoints": status, // 1 (true) tidak pake point, 0(false) pake pint
	}

	// api := campaignID + "/buy"
	api := custId + "/campaign/" + campaignID + "/buy"
	urlSvr := host + endpointRedeemCumulativeVoucher + api

	data, err := HTTPxFormPostAdminWithRequest(urlSvr, jsonData)
	if err != nil {
		logs.Error("Check error", err.Error())

		return &resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logs.Error("Failed to unmarshaling response RedeemVoucherCumulative from open-loyalty ", err.Error())

		return &resp, err
	}

	return &resp, nil
}

// RulePoint
func RulePoint(eventName, phone string) (models.RulePointResponse, error) {
	var resp models.RulePointResponse

	logrus.Info("[Package Host OPL]-[RulePoint]")

	todo := endpointRulePoint + eventName
	logrus.Info("Request EranRule :", todo)

	urlSvr := host + todo

	data, err := HTTPxFormPostCustomerWithoutRequest(urlSvr, phone)
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

	logrus.Info("[Package Host OPL]-[ListRulePoint]")

	todo := endpointListRulePoint

	urlSvr := host + todo

	data, err := HTTPxFormGETCustomer(urlSvr, phone)
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

	logrus.Info("[Package Host OPL]-[VoucherDetail]")

	urlSvr := host + endpointVoucherDetail + campaign

	data, err := HTTPxFormGETAdmin(urlSvr)
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

	logrus.Info("[Package Host OPL]-[HistoryVoucherCustomer]")

	param := fmt.Sprintf("?includeDetails=1&page=%s&perPage=1000&sort&direction", page)
	urlSvr := host + endpointHistoryVoucherCustomer + param
	data, err := HTTPxFormGETCustomer(urlSvr, phone)
	if err != nil {
		logs.Error("Check error ", err.Error())
		return &resp, err
	}

	logrus.Info("Response OPL")

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

	logrus.Info("[Package Host OPL]-[CouponVoucherCustomer]")

	urlSvr := host + endpointCouponVoucherCustomer

	jsonData := map[string]interface{}{
		"coupons[0][campaignId]": campaign,   //"coupons[0][campaignId]": campaign,
		"coupons[0][couponId]":   couponId,   //"coupons[0][couponId]":   couponId,
		"coupons[0][code]":       couponCode, //"coupons[0][code]":       couponCode,
		"coupons[0][used]":       useVoucher, //"coupons[0][used]":       "true"}
		"coupons[0][customerId]": custID}

	logrus.Info("===== Use Voucher True / False =====")
	data, err := HTTPxFormPostAdminWithRequest(urlSvr, jsonData)
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

	logrus.Info("[Package Host OPL]-[TransferPoint]")
	urlSvr := host + endpointAddedPoint

	jsonData := map[string]interface{}{
		"transfer[customer]": customer,
		"transfer[points]":   point,
		"transfer[comment]":  text,
	}

	logrus.Info("Request to OPL : ", jsonData)
	data, err := HTTPxFormPostAdminWithRequest(urlSvr, jsonData)
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

	logrus.Info("[Package Host OPL]-[SpendPoint]")

	urlSvr := host + endpointSpendPoint

	jsonData := map[string]interface{}{
		"transfer[customer]": customer,
		"transfer[points]":   point,
		"transfer[comment]":  text,
	}

	logrus.Info("Request to OPL : ", jsonData)
	data, err := HTTPxFormPostAdminWithRequest(urlSvr, jsonData)
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

	logrus.Info("[Package Host OPL]-[GetBalance]")

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

// Settings ..
func SettingsOPL() (*models.SettingOPL, error) {
	var resp models.SettingOPL

	fmt.Println("[Package Host OPL]-[SettingsOPL]")
	urlSvr := host + endpointSetting

	data, err := HTTPxFormGETAdmin(urlSvr)
	if err != nil {
		logs.Error("Check error ", err.Error())
		return &resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logs.Error("Failed to unmarshaling response SettingsOPL from open-loyalty ", err.Error())
		return &resp, err
	}
	return &resp, nil
}
