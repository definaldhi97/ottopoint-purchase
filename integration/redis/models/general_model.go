package models

type RedisBody struct {
	Value string `json:"value"`
}

type RedisHead struct {
	Key string `json:"key"`
	// Action string `json:"action"`
	Expire string `json:"expire"`
}

// type ResponseRedis struct {
// 	ResponseCode string    `json:"responseCode"`
// 	ResponseDesc string    `json:"responseDesc"`
// 	Data         DataValue `json:"value"`
// }

type DataValue struct {
	AccountNumber string `json:"account_number"`
	MerchantID    string `json:"merchant_id"`
}

type ResponseApi struct {
	ResponseCode string `json:"responseCode"`
	ResponseDesc string `json:"responseDesc"`
}

type ResponseRedis1 struct {
	ResponseCode string `json:"responseCode"`
	ResponseDesc string `json:"responseDesc"`
	Value        string `json:"value"`
}
