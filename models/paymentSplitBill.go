package models

// Request
type PaymentSplitBillReq struct {
	CampaignId         string      `json:"campaignId"`
	FieldValue         interface{} `json:"fieldValue"`
	PaymentMethod      int         `json:"paymentMethod"`
	PaymentAccount     int         `json:"paymentAccount"`
	SavePaymentAccount bool        `json:"savePaymentAccount"`
	Total              int         `json:"total"`
}

// Response
type PaymentSplitBillResp struct {
	Code       string `json:"code"`
	Message    string `json:"msg"`
	Success    int    `json:"success"`
	Failed     int    `json:"failed"`
	Pending    int    `json:"pending"`
	UrlPayment string `json:"urlPayment"`
}
