package models

// RulePoint
type RulePointReq struct {
	Amount    int    `json:"amount"`
	EventName string `json:"rule"`
	RC        string `json:"rc"`
}

type RulePointResp struct {
	Point       string `json:"point"`
	Institution string `json:"institution"`
}
