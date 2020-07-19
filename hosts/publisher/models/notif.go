package models

type NotifVoucher struct {
	AccountNumber string `json:"accountNumber"`
	// Institution   string `json:"institution"`
	VoucherName string `json:"voucherName"`
	ExpiredDate string `json:"expiredDate"`
}
