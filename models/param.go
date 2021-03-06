package models

import "time"

type Params struct {
	ResponseCode  int
	AccountNumber string
	Email         string
	FirstName     string
	LastName      string
	MerchantID    string
	InstitutionID string
	TransType     string
	AccountId     string
	TrxID         string // internal
	CumReffnum    string // internal (generate ottopoint) untuk pembelian kelipatan
	RRN           string // eksternal (from supplier exp : ottoag, uv)
	InvoiceNumber string
	CustID        string
	Amount        int64
	NamaVoucher   string
	ProductType   string
	ProductCode   string
	ProductName   string
	Category      string
	Point         int
	ExpDate       string
	CouponID      string
	CouponCode    string
	CampaignID    string
	SupplierID    string

	VoucherCode     string // voucher code vidio
	VoucherLink     string // voucher link
	ExpireDateVidio string // Expire date Voucher Vidio
	TrxTime         time.Time
	CategoryID      *string
	// percobaan
	Total               int
	DataSupplier        Supplier
	UsageLimitVoucher   int
	ProductCodeInternal string
	ProductID           string
	Comment             string
	RewardID            string
	PointTransferID     string
	Fields              []string
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

type ParamOrder struct {
	Nama    string
	Email   string
	Phone   string
	Expired string
}
