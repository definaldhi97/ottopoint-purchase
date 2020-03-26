package models

// Request Order
type OrderVoucherReq struct {
	Sku               string `json:"sku"`
	Qty               int    `json:"qty"`
	AccountID         string `json:"account_id"`
	InstitutionRefno  string `json:"institution_refno"`
	ExpireDateVoucher int    `json:"expire_date_voucher"`
	ReceiverName      string `json:"receiverName"`
	ReceiverEmail     string `json:"receiverEmail"`
	ReceiverPhone     string `json:"receiverPhone"`
}

// Response Order
type OrderVoucherResp struct {
	ResponseCode string    `json:"responseCode"`
	ResponseDesc string    `json:"responseDesc"`
	Data         DataOrder `json:"Data"`
}

type DataOrder struct {
	OrderID           string            `json:"orderId"`
	InvoiceOp         string            `json:"invoiceOp"`
	VouchersAvailable int               `json:"vouchersAvailable"`
	VouchersCode      []DataVoucherCode `json:"vouchersCode"`
}

type DataVoucherCode struct {
	Code string `json:"code"`
}
