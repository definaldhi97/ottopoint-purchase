package ottoagmodels

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
	Rc          string      `json:"rc"`
	Rrn         string      `json:"rrn"`
	Adminfee    uint64      `json:"adminfee"`
	Amount      uint64      `json:"amount"`
	Custid      string      `json:"custid"`
	Memberid    string      `json:"memberid"`
	Msg         string      `json:"msg"`
	Period      string      `json:"period"`
	Productcode string      `json:"productcode"`
	Uimsg       string      `json:"uimsg"`
	Data        interface{} `json:"data"`
}
