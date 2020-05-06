package models

// RulePoint
type RulePointReq struct {
	Amount    int    `json:"amount"`
	EventName string `json:"rule"`
}

type RulePointResp struct {
	Point       int    `json:"point"`
	Product     string `json:"product"`
	Institution string `json:"institution"`
}

type LisrRulePointResp struct {
	EarningRules []GetEarningRulesResp `json:"earningRules"`
	Currency     string                `json:"currency"`
}

type GetEarningRulesResp struct {
	Name         string `json:"name"`
	EventName    string `json:"eventName"`
	PointsAmount int    `json:"pointsAmount"`
}
