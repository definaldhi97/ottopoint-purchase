package models

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
	TrxID           string
	VoucherCode     string // voucher code vidio
	ExpireDateVidio string // Expire date Voucher Vidio
	VoucherLink     string // voucher link
	CategoryID      string
	// percobaan
	Total               int
	DataSupplier        Supplier
	UsageLimitVoucher   int
	ProductCodeInternal string
	ProductID           string
	Comment             string
	RewardID            string
}

type Supplier struct {
	Request  string // Request dijadikan byte[] >> string
	Response string // Response dijadikan byte[] >> string
	Rd       string // msg response
	Rc       string // rc response
}
