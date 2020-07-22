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

type EarningRuleReq struct {
	Code string `json:"code"`
}

// ============== New Earning ==============
type EarningReq struct {
	Earning        string `json:"earning"`
	ReferenceId    string `json:"referenceId"`
	ProductCode    string `json:"productCode"`
	ProductName    string `json:"productName"`
	AccountNumber1 string `json:"accountNumber1"`
	AccountNumber2 string `json:"accountNumber2"`
	Amount         int64  `json:"amount"`
	Remark         string `json:"remark"`
	// MsgNotif        string `json:"messageNotif"`
	TransactionTime string `json:"transactionTime"`
}

type EarningResp struct {
	ReferenceId string `json:"referenceId"`
	Point       int64  `json:"point"`
}

type ResponseEarning struct {
	// Code           string        `json:"code"`
	// Message        string        `json:"message"`
	ReferenceId    string        `json:"referenceId`
	TransactionId  string        `json:"transactionId"`
	PartnerId      string        `json:"partnerId"`
	Amount         int64         `json:"amount"`
	EarningRule    string        `json:"earningRule"`
	EarningRuleAdd string        `json:"errningRuleAdd"`
	EarningDate    string        `json:"earningDate"`
	AccountData    []DataEarning `json:"accountData"`
}

type DataEarning struct {
	AccountNumber string `json:"accountNumber"`
	Point         int64  `json:"point"`
}

// ============== Check Status Earning ==============
type CheckStatusEarningReq struct {
	ReferenceId string `json:"referenceId"`
}
