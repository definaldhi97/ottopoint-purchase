package models

type PaymentInfoResp struct {
	ResponseCode string `json:"responseCode"`
	ResponseDesc string `json:"responseDesc"`
	Data         struct {
		VoucherType       string `json:"voucherType"`
		IsAsyncOrder      bool   `json:"isAsyncOrder"`
		IsCumulativeOrder bool   `json:"isCumulativeOrder"`
		Fields            []struct {
			Key         string `json:"key"`
			Label       string `json:"label"`
			IsNumeric   bool   `json:"isNumeric"`
			Maxlength   int    `json:"maxlength"`
			CharSplit   string `json:"charSplit"`
			LengthSplit int    `json:"lengthSplit"`
			Function    string `json:"function"`
		} `json:"fields"`
	} `json:"data"`
}
