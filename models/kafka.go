package models

type PublishEarningReq struct {
	Header         RequestHeader `json:"header"`
	Earning        string        `json:"earning"`
	ReferenceId    string        `json:"referenceId"`
	ProductCode    string        `json:"productCode"`
	ProductName    string        `json:"productName"`
	AccountNumber1 string        `json:"accountNumber1"`
	AccountNumber2 string        `json:"accountNumber2"`
	Amount         int64         `json:"amount"`
	Remark         string        `json:"remark"`
}

type NotifPubreq struct {
	Type          string `json:"type"`
	AccountNumber string `json:"accountNumber"`
	Institution   string `json:"institution"`
	Point         int    `json:"point"`
	Product       string `json:"product"`
}
