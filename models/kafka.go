package models

import "time"

type PublishEarningReq struct {
	Header          RequestHeader `json:"header"`
	Earning         string        `json:"earning"`
	ReferenceId     string        `json:"referenceId"`
	ProductCode     string        `json:"productCode"`
	ProductName     string        `json:"productName"`
	AccountNumber1  string        `json:"accountNumber1"`
	AccountNumber2  string        `json:"accountNumber2"`
	Amount          int64         `json:"amount"`
	Remark          string        `json:"remark"`
	TransactionTime time.Time     `json:"transactionTime"`
}

type NotifPubreq struct {
	Type           string    `json:"notificationType"` // PLN, Earning, Reversal
	NotificationTo string    `json:"notificationTo"`   // AccountNumber
	Institution    string    `json:"institutionId"`
	ReferenceId    string    `json:"referenceId"`
	TransactionId  string    `json:"transactionId"`
	Data           DataValue `json:"data"`
}

// type NotifPubreq struct {
// 	Type          string `json:"type"`
// 	AccountNumber string `json:"accountNumber"`
// 	Institution   string `json:"institution"`
// 	Point         int    `json:"point"`
// 	Product       string `json:"product"`
// }

type DataValue struct {
	RewardValue string `json:"rewardType"` // type point & voucher
	Value       string `json:"value"`      // point & nama voucher
}
