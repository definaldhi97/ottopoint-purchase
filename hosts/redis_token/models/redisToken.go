package models

type TokenResp struct {
	ResponseCode string    `json:"responseCode"`
	ResponseDesc string    `json:"responseDesc"`
	Data         DataValue `json:"data"`
}

type DataValue struct {
	AccountNumber string `json:"account_number"`
	MerchantID    string `json:"merchant_id"`
}
