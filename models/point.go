package models

type PointReq struct {
	Type  string `json:"type"`
	Point int    `json:"point"`
	Text  string `json:"text"`
}

type PointResp struct {
	Nama          string `json:"nama"`
	AccountNumber string `json:"accountNumber"`
	Point         int    `json:"points"`
	Text          string `json:"text"`
}
