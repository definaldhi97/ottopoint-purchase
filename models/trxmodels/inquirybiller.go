package trxmodels

// BillerInquiryRequest ..
type BillerInquiryRequest struct {
	TypeTrans     string      `json:"typeTrans"`
	Datetime      string      `json:"datetime"`
	IssuerID      string      `json:"issuerId"`
	Amount        int64       `json:"amount"`
	AccountNumber string      `json:"accountNumber,omitempty"`
	Token         string      `json:"token,omitempty"`
	TransactionID string      `json:"transactionId,omitempty"`
	Data          interface{} `json:"data,omitempty"`
}

// BillerInquiryDataReq ..
type BillerInquiryDataReq struct {
	ProductCode string `valid:"Required" json:"productcode"`
	MemberID    string `valid:"Required" json:"memberid"`
	CustID      string `valid:"Required" json:"custid"`
	Period      string `valid:"Required" json:"period"`
}

// BillerInquiryResponse ..
type BillerInquiryResponse struct {
	Rc            string      `json:"rc"`
	Rrn           string      `json:"rrn"`
	Datetime      string      `json:"datetime"`
	Adminfee      int         `json:"adminfee"`
	Amount        int64       `json:"amount"`
	Balance       int64       `json:"balance"`
	Token         string      `json:"token"`
	TransactionID string      `json:"transactionId"`
	Msg           string      `json:"msg"`
	Uimsg         string      `json:"uimsg"`
	Data          interface{} `json:"data"`
}

// BillerInquiryResponse ..
type BillerInquiryResponsePulsa struct {
	Rc            string       `json:"rc"`
	Rrn           string       `json:"rrn"`
	Datetime      string       `json:"datetime"`
	Adminfee      int          `json:"adminfee"`
	Amount        int64        `json:"amount"`
	Balance       int64        `json:"balance"`
	Token         string       `json:"token"`
	TransactionID string       `json:"transactionId"`
	Msg           string       `json:"msg"`
	Uimsg         string       `json:"uimsg"`
	Data          DataInqPulsa `json:"data"`
}

// Data Inquiry Pulsa
type DataInqPulsa struct {
	Name        string `json:"Name"`
	ProductCode string `json:"ProductCode"`
	Denom       string `json:"Denom"`
	SalesPrice  string `json:"SalesPrice"`
	Company     string `json:"Company"`
}
