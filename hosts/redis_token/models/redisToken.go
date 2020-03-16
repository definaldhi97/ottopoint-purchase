package models

type TokenResp struct {
	ResponseCode string `json:"responseCode"`
	ResponseDesc string `json:"responseDesc"`
	Data         string `json:"value"`
	MerchantID   string `json:"merchantId"`
}

type DataValue struct {
	AccountNumber string `json:"account_number"`
	MerchantID    string `json:"merchant_id"`
}
