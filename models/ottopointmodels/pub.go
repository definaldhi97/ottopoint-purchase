package ottopointmodels

// Publish
type PubOttopoint struct {
	Rc             string `json:"rc"`
	RRN            string `json:"rrn"`
	ReffNumber     string `json:"reff_number"`
	ProductCode    string `json:"product_code"`
	AcccountNumber string `json:"acccount_number"` // pengirim
	Phone          string `json:"phone"`           // penerima
	TypeTrans      string `json:"typeTrans"`
	TypeTRX        string `json:"type_trx"`
	ProductType    string `json:"product_type"`
	Amount         int64  `json:"amount"`
	FeeAmount      string `json:"fee_amount"`
	Datetime       string `json:"datetime"`
	MechantID1     string `json:"merchant_id1"` // pengerim
	MechantID2     string `json:"merchant_id2"` // penerima
	ProductName    string `json:"productName"`
	CustID         string `json:"cust_id"`
	BillerID       string `json:"billerId"`
	IssuerID       string `json:"issuerId"`
	Token          string `json:"token"`
	DeviceID       string `json:"device-id"`
}
