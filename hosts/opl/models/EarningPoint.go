package models

type RulePointResponse struct {
	Point int `json:"points"`
}

type LisrRulePointResponse struct {
	EarningRules []GetEarningRulesResponse `json:"earningRules"`
	Currency     string                    `json:"currency"`
}

type GetEarningRulesResponse struct {
	Name         string `json:"name"`
	EventName    string `json:"eventName"`
	PointsAmount int    `json:"pointsAmount"`
}
