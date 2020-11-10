package models

import "time"

type Params struct {
	AccountNumber   string
	MerchantID      string
	InstitutionID   string
	TransType       string
	AccountId       string //opl
	Reffnum         string // internal (generate ottopoint)
	CumReffnum      string // internal (generate ottopoint) untuk pembelian kelipatan
	RRN             string // eksternal (from supplier exp : ottoag, uv)
	CustID          string
	Amount          int64
	NamaVoucher     string
	ProductType     string
	ProductCode     string
	Category        string
	Point           int
	ExpDate         string
	CouponID        string
	CouponCode      string
	CampaignID      string
	SupplierID      string
	VoucherCode     string // voucher code vidio
	VoucherLink     string // voucher link
	ExpireDateVidio string // Expire date Voucher Vidio
	CategoryID      string
	TrxID           string
	TrxTime         time.Time
	// percobaan
	Total        int
	DataSupplier Supplier
}

type Supplier struct {
	Request  string // Request dijadikan byte[] >> string
	Response string // Response dijadikan byte[] >> string
	Rd       string // msg response
	Rc       string // rc response
}

type DataErrorOPL struct {
	Error string `json:"error"`
}
