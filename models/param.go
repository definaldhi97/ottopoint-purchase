package models

type Params struct {
	AccountNumber string
	MerchantID    string
	InstitutionID string
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
	Total int
}
