package models

type GetVoucherAgResp struct {
	Voucher     string `json:"voucher"`
	VoucherCode string `json:"voucherCode"`
	Link        string `json:"link"`
}

type CallbackRequestVoucherAg struct {
	InstitutionID    string              `json:"institutionId"`
	NotificationType string              `json:"notificationType"`
	NotificationTo   string              `json:"notificationTo"`
	TransactionID    string              `json:"transactionId"`
	Data             CallbackRequestData `json:"data"`
}

type CallbackRequestData struct {
	OrderID      string `json:"orderId"`
	VoucherID    string `json:"voucherId"`
	VoucherCode  string `json:"voucherCode"`
	VoucherName  string `json:"voucherName"`
	Status       string `json:"status"`
	IsRedeemed   bool   `json:"isRedeemed"`
	RedeemedDate string `json:"redeemedDate"`
	UsedDate     string `json:"usedDate"`
}
