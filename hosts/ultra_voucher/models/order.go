package models

// Request Order
type OrderVoucherReq struct {
	Sku               string `json:"sku"`
	Qty               int    `json:"qty"`
	AccountID         string `json:"accountId"`
	InstitutionRefno  string `json:"institutionRefno"`
	ExpireDateVoucher int    `json:"expireDateVoucher"`
	ReceiverName      string `json:"receiverName"`
	ReceiverEmail     string `json:"receiverEmail"`
	ReceiverPhone     string `json:"receiverPhone"`
}

// Response Order
type OrderVoucherResp struct {
	ResponseCode string    `json:"responseCode"`
	ResponseDesc string    `json:"responseDesc"`
	Data         DataOrder `json:"data"`
}

type DataOrder struct {
	OrderID           string            `json:"orderId"`
	InstReffnum       string            `json:"instReffnum"`
	InvoiceOp         string            `json:"invoiceOp"`
	InvoiceUV         string            `json:"InvoiceUV"`
	VouchersAvailable string            `json:"vouchersAvailable"`
	VouchersCode      []DataVoucherCode `json:"vouchersCode"`
}

type DataVoucherCode struct {
	Code string `json:"code"`
}
