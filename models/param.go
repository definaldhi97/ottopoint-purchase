package models

type Params struct {
	AccountNumber string
	MerchantID    string
	InstitutionID string
	TransType     string
	AccountId     string //opl
	Reffnum       string // internal
	RRN           string // eksternal
	CustID        string
	Amount        int64
	NamaVoucher   string
	ProductType   string
	ProductCode   string
	Category      string
	Point         int
	ExpDate       string
	CouponID      string
	CouponCode    string
	CampaignID    string
	SupplierID    string
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
