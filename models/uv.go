package models

type GetVoucherUVResp struct {
	Voucher     string `json:"voucher"`
	VoucherCode string `json:"voucherCode"`
	Link        string `json:"link"`
}

type UseVoucherUVResp struct {
	Voucher string `json:"voucher"`
}
