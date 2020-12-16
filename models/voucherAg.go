package models

type GetVoucherAgResp struct {
	Voucher     string `json:"voucher"`
	VoucherCode string `json:"voucherCode"`
	Link        string `json:"link"`
}
