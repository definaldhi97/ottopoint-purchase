package ottoag

// BillerInquiryDataReq ..
type BillerInquiryDataReq struct {
	CustID      string `valid:"Required" json:"custid"`
	MemberID    string `valid:"Required" json:"memberid"`
	Period      string `valid:"Required" json:"period"`
	ProductCode string `valid:"Required" json:"productcode"`
}

// InquiryRequest ..
type OttoAGInquiryRequest struct {
	TypeTrans      string      `json:"typeTrans"`
	Datetime       string      `json:"datetime"`
	IssuerID       string      `json:"issuerId"`
	Amount         int64       `json:"amount"`
	AccountNumber  string      `json:"accountNumber,omitempty"`
	Token          string      `json:"token,omitempty"`
	TransactionID  string      `json:"transactionId,omitempty"`
	Data           interface{} `json:"data,omitempty"`
	MerchantEmoney string      `json:"merchant_emoney"`
}

// InquiryResponse ..
type OttoAGInquiryResponse struct {
	Rc          string      `json:"rc"`
	Rrn         string      `json:"rrn"`
	Adminfee    int         `json:"adminfee"`
	Amount      int64       `json:"amount"`
	CustID      string      `json:"custid"`
	MemberID    string      `json:"memberid"`
	Msg         string      `json:"msg"`
	Periode     string      `json:"period"`
	ProductCode string      `json:"productcode"`
	Uimsg       string      `json:"uimsg"`
	Data        interface{} `json:"data"`
}
