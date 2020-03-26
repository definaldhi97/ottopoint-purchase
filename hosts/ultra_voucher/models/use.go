package models

// Request Use Voucher UV
type UseVoucherUVReq struct {
	Account     string `json:"account"`
	VoucherCode string `json:"voucherCode"`
}

// Response Use Voucher UV
type UseVoucherUVResp struct {
	ResponseCode string    `json:"responseCode"`
	ResponseDesc string    `json:"responseDesc"`
	Data         DataUseUV `json:"data"`
}

type DataUseUV struct {
	Sku                string `json:"sku"`
	ReffCode           string `json:"reffCode"`
	SuplierVoucherCode string `json:"suplierVoucherCode"`
	SuplierVoucherID   string `json:"suplierVoucherId"`
	Link               string `json:"link"`
	OrderNo            string `json:"orderNo"`
}
