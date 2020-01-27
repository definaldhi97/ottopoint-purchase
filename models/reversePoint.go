package models

type ReversePointReq struct {
	TrxID string `json:"trx_id"`
}

type ReversePointResp struct {
	Nama          string `json:"nama"`
	AccountNumber string `json:"accountNumber"`
	Point         int    `json:"points"`
	Text          string `json:"text"`
}
