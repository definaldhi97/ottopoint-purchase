package trxmodels

// BillerPaymentRequest ..
type BillerPaymentRequest struct {
	Rrn             string                `json:"rrn"`
	TypeTrans       string                `json:"typeTrans"`
	Datetime        string                `json:"datetime"`
	QRCode          string                `json:"qrCode,omitempty"`
	IssuerID        string                `json:"issuerId"`
	AcquirerID      string                `json:"acquirerId"`
	MerchantID      string                `json:"merchantId,omitempty"`
	MerchantPhone   string                `json:"merchantPhone,omitempty"`
	MerchantName    string                `json:"merchantName,omitempty"`
	AccountNumber   string                `json:"accountNumber"`
	Amount          int64                 `json:"amount"`
	Fee             int64                 `json:"fee"`
	Commission      int64                 `json:"commission"`
	ReferenceNumber string                `json:"referenceNumber"`
	BillerID        string                `json:"billerId"`
	ProductName     string                `json:"productName"`
	ProductCode     string                `json:"productCode"`
	PartnerCode     string                `json:"partnerCode"`
	BillerData      PurchaseBillerDataReq `json:"billerData,omitempty"`
	MerchantEmoney  string                `json:"merchantEmoney"`
}

// PurchaseBillerDataReq ..
type PurchaseBillerDataReq struct {
	BillerCode   string `json:"billercode,omitempty"`  //mandatory
	ProductCode  string `json:"productcode,omitempty"` //mandatory
	BillRef      string `json:"billref,omitempty"`     //mandatory
	SubscriberID string `json:"subscriber,omitempty"`  //mandatory
	CustID       string `json:"custid,omitempty"`
	Input1       string `json:"input1,omitempty"` //mandatory
	Input2       string `json:"input2,omitempty"` //mandatory
}

// BillerPaymentResponse ..
type BillerPaymentResponse struct {
	Rc             string                 `json:"rc"`
	Rrn            string                 `json:"rrn"`
	Msg            string                 `json:"msg"`
	Datetime       string                 `json:"datetime"`
	Data           PurchaseDateResponse   `json:"data,omitempty"`
	BillerData     PurchaseBillerResponse `json:"billerData,omitempty"`
	BillerDataResp BillerDataResponse     `json:"billerDataResp,omitempty"`
	UIMsg          string                 `json:"uimsg,omitempty"`
}

// PurchaseDateResponse ..
type PurchaseDateResponse struct {
	Amount    int64  `json:"amount"`
	Refnumber string `json:"rrn"`
	IssuerRef string `json:"issuerref"`
	SisaSaldo int64  `json:"saldo"`
}

// PpobDetail ...
type PpobDetail struct {
	Amount     int    `json:"amount"`
	Penalty    int    `json:"penalty"`
	Period     string `json:"period"`
	Watermeter string `json:"watermeter"`
}

// BillerDataResponse ...
type BillerDataResponse struct {
	Msisdn       string `json:"msisdn"`
	Timestamp    string `json:"timestamp"`
	SerialNumber string `json:"serialnumber"`
	BillerRef    string `json:"billref"`
	Ref          string `json:"ref"`

	KwhTotal string `json:"kwhtotal"`
	TokenNo  string `json:"tokenno"`

	NomorTransaksi string `json:"nomortransaksi"`
	PinKode        string `json:"pinkode"`
	Serial         string `json:"serial"`

	// BEGIN define BPJS purchase response attributes
	Months			int64			`json:"months,omitempty"`
	HeadVa			string			`json:"headva,omitempty"`
	TotalPremium	int64			`json:"totalpremium,omitempty"`
	NumFm			int64			`json:"numfm,omitempty"`
	FamilyMembers	[]FamilyMember	`json:"familymembers,omitempty"`
	JpaRef			string			`json:"jparef,omitempty"`
	PaidDate		string			`json:"paiddate,omitempty"`
	// END define BPJS purchase response attributes
}

// PurchaseBillerResponse ..
type PurchaseBillerResponse struct {
	BillRef string      `json:"billref,omitempty"` //mandatory
	CustID  string      `json:"custid,omitempty"`
	Input1  string      `json:"input1,omitempty"` //mandatory
	Input2  string      `json:"input2,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// FamilyMember ...
type FamilyMember struct {
	BranchCode		string	`json:"branchcode,omitempty"`
	BranchName		string	`json:"branchname,omitempty"`
	FmPremiumCode	string	`json:"fmpremiumcode,omitempty"`
	FmVa			string	`json:"fmva,omitempty"`
	Name			string	`json:"name,omitempty"`
	Premium			int64	`json:"premium,omitempty"`
	PremiumDp		int64	`json:"premiumdp,omitempty"`
	PremiumMonth	int64	`json:"premiummonth,omitempty"`
}
