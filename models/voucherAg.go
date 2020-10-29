package models

type GetVoucherAgResp struct {
	Voucher     string `json:"voucher"`
	VoucherCode string `json:"voucher_code"`
	Link        string `json:"link"`
}
