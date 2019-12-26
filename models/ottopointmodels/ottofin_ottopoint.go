package ottopointmodels

// Request Inquiry
type OttofinOttopointReq struct {
	IssuerID      string `json:"issuerId"`
	AccountNumber string `json:"accountNumber"`
	Amount        int    `json:"amount"`
	Refnumber     string `json:"refnumber"`
	Datetime      string `json:"datetime"`
}

type OttofinOttopointTopupReq struct {
	IssuerID      string `json:"issuerId"`
	AccountNumber string `json:"accountNumber"`
	Amount        int    `json:"amount"`
	Refnumber     string `json:"refnumber"`
	Datetime      string `json:"datetime"`
	Member        string `json:"member"`
}

// Request Reversal
type ReversalOttofinOttopointReq struct {
	IssuerID      string `json:"issuerId"`
	AccountNumber string `json:"accountNumber"`
	Refnumber     string `json:"refnumber"`
}

// Response
type OttofinOttopointResp struct {
	Rc       string `json:"rc"`
	Rrn      string `json:"rrn"`
	Msg      string `json:"msg"`
	Datetime string `json:"datetime"`
	Data     struct {
		AccountNumber string `json:"accountNumber"`
		AccountName   string `json:"accountName"`
		Saldo         string `json:"saldo"`
	}
}
