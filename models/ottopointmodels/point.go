package ottopointmodels

type ResponsePoint struct {
	Data TransferPointResp `json:"data"`
	Meta interface{}       `json:"meta"`
}

type TransferPointResp struct {
	Nama          string `json:"nama"`
	AccountNumber string `json:"accountNumber"`
	Point         int    `json:"points"`
	Text          string `json:"text"`
}
