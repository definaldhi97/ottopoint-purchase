package trxmodels

// ReversalRequest ..
type ReversalRequest struct {
	Datetime        string `json:"datetime"`
	IssuerID        string `json:"issuerId"`
	ReferenceNumber string `json:"referenceNumber"`
}

// ReversalResponse ..
type ReversalResponse struct {
	Rc       string `json:"rc"`
	Rrn      string `json:"rrn"`
	Msg      string `json:"msg"`
	Datetime string `json:"datetime"`
}
