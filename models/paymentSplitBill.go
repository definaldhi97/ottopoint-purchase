package models

// Request
type PaymentSplitBillReq struct {
	CampaignId    string      `json:"campaignId"`
	FieldValue    interface{} `json:"fieldValue"`
	PaymentMethod int         `json:"paymentMethod"` // 0. Full point, 1. Spilt bill Ottocash

	// PaymentAccount     int         `json:"paymentAccount"`
	// SavePaymentAccount bool        `json:"savePaymentAccount"`
	// Total              int         `json:"total"`
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

type CallBackSGReq struct {
	Amount              string `json:"amount"`
	Issuer              string `json:"issuer"`
	IssuerRefNo         string `json:"issuerRefNo"`
	OttoRefNo           string `json:"ottoRefNo"`
	ResponseCode        string `json:"responseCode"`
	ResponseDescription string `json:"responseDescription"`
	TrxRef              string `json:"trxRef"`
	TransactionType     string `json:"transactionType"`
	UserId              string `json:"userId"`
}
