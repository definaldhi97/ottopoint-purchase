package trxmodels

// PurchasePayReq ...
type PurchasePayReq struct {
	AccountNumber           string `json:"accountNumber"`
	Amount                  int64  `json:"amount"`
	Fee                     int64  `json:"fee"`
	BillerID                string `json:"billerId"`
	CustomerReferenceNumber string `json:"customerReferenceNumber"`
	ProductName             string `json:"productName"` //:"Top up fr sales",
	ProductCode             string `json:"productCode"` //:"Top Up",
	PartnerCode             string `json:"partnerCode"` //:"rrn OttoAg"
}
