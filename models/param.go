package models

type Params struct {
	AccountNumber string
	MerchantID    string
	InstitutionID string
	CustID        string
	Reffnum       string // internal
	RRN           string // eksternal
	Amount        int64
	NamaVoucher   string
	ProductType   string
	ProductCode   string
	Category      string
	Point         int
	ExpDate       string
	SupplierID    string
}
