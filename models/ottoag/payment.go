package ottoag

// OttoAGPaymentReq ...
type OttoAGPaymentReq struct {
	Amount      uint64 `json:"amount"`
	CustID      string `json:"custid"`
	MemberID    string `json:"memberid"`
	Period      string `json:"period"`
	Productcode string `json:"productcode"`
	Rrn         string `json:"rrn"`
}

// OttoAGPaymentRes ..
type OttoAGPaymentRes struct {
	Rc          string       `json:"rc"`
	Rrn         string       `json:"rrn"`
	Adminfee    uint64       `json:"adminfee"`
	Amount      uint64       `json:"amount"`
	Custid      string       `json:"custid"`
	Memberid    string       `json:"memberid"`
	Msg         string       `json:"msg"`
	Period      string       `json:"period"`
	Productcode string       `json:"productcode"`
	Uimsg       string       `json:"uimsg"`
	Data        DataGabungan `json:"data"`
}

type DataGabungan struct {
	// PLN
	Custname    string `json:"custname"`
	Meterno     string `json:"meterno"`
	Power       string `json:"power"`
	Amount      int    `json:"amount"`
	Ppn         int    `json:"ppn"`
	Ppj         int    `json:"ppj"`
	Installment int    `json:"installment"`
	Stampduty   int    `json:"stampduty"`
	Stroomtoken string `json:"stroomtoken"`
	Tokenno     string `json:"tokenno"`
	Kwhtotal    string `json:"kwhtotal"`
	Ref         string `json:"ref"`
	Invotext    string `json:"invotext"`
	// Pulsa
	Msisdn       string `json:"msisdn"`
	Timestamp    string `json:"timestamp"`
	Serialnumber string `json:"serialnumber"`
	Billref      string `json:"billref"`
	// Game
	CustName string `json:"custname"`
	BillRef  string `json:"billref"`
	// vidio
	StartDateVidio string `json:"startDate"`
	EndDateVidio   string `json:"endDate"`
	Code           string `json:"code"`
	Description    string `json:"description"`
}

// Data Payment PLN Token
type DataPayPLNTOKEN struct {
	Custname    string `json:"custname"`
	Meterno     string `json:"meterno"`
	Power       string `json:"power"`
	Amount      int    `json:"amount"`
	Ppn         int    `json:"ppn"`
	Ppj         int    `json:"ppj"`
	Installment int    `json:"installment"`
	Stampduty   int    `json:"stampduty"`
	Stroomtoken string `json:"stroomtoken"`
	Tokenno     string `json:"tokenno"`
	Kwhtotal    string `json:"kwhtotal"`
	Ref         string `json:"ref"`
	Invotext    string `json:"invotext"`
}

// Data Payment Game
type DataGame struct {
	CustName string `json:"custname"`
	BillRef  string `json:"billref"`
}

// Data Payment Pulsa
type DataPayPulsa struct {
	Msisdn       string `json:"msisdn"`
	Timestamp    string `json:"timestamp"`
	Serialnumber string `json:"serialnumber"`
	Billref      string `json:"billref"`
}
