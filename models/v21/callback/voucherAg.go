package callback

type CallbackVoucherAGReq struct {
	InstitutionId    string              `json:"institutionId"`
	NotificationType string              `json:"notificationType"`
	NotificationTo   string              `json:"notificationTo"`
	TransactionId    string              `json:"transactionId"`
	VoucherType      string              `json:"voucherType"`
	Data             DataVoucherTypePPOB `json:"data"`
}

type CallbackVoucherAGReq1 struct {
	InstitutionId    string               `json:"institutionId"`
	NotificationType string               `json:"notificationType"`
	NotificationTo   string               `json:"notificationTo"`
	TransactionId    string               `json:"transactionId"`
	VoucherType      string               `json:"voucherType"`
	Data             DataVoucherTypeMerge `json:"data"`
}

// Voucher type PPOB
type DataVoucherTypePPOB struct {
	OrderId      string `json:"orderId"`
	ResponseCode string `json:"responseCode"`
	ResponseDesc string `json:"responseDesc"`
}

// Voucher type Voucher Code
type DataVoucherTypeVoucherCode struct {
	OrderId      string `json:"orderId"`
	VoucherId    string `json:"voucherId"`
	VoucherCode  string `json:"voucherCode"`
	VoucherName  string `json:"voucherName"`
	IsRedeemed   bool   `json:"isRedeemed"`
	RedeemedDate string `json:"redeemedDate"`
}

// Merge
type DataVoucherTypeMerge struct {
	OrderId      string `json:"orderId"`
	VoucherId    string `json:"voucherId"`
	VoucherCode  string `json:"voucherCode"`
	VoucherName  string `json:"voucherName"`
	IsRedeemed   bool   `json:"isRedeemed"`
	RedeemedDate string `json:"redeemedDate"`
	ResponseCode string `json:"responseCode"`
	ResponseDesc string `json:"responseDesc"`
}
