package dbmodels

import "time"

type TEarning struct {
	ID               string    `gorm:"column:id"`
	EarningRule      string    `gorm:"column:earning_rule"`
	EarningRuleAdd   string    `gorm:"column:earning_rule_add"`
	PartnerId        string    `gorm:"column:partner_id"`
	ReferenceId      string    `gorm:"column:reference_id"`
	TransactionId    string    `grom:"column:transaction_id"`
	ProductCode      string    `gorm:"column:product_code"`
	ProductName      string    `gorm:"column:product_name"`
	AccountNumber    string    `gorm:"column:account_number"`
	Amount           int64     `gorm:"column:amount"`
	Point            int64     `gorm:"column:point"`
	Remark           string    `gorm:"column:remark"`
	Status           string    `gorm:"column:status"`
	StatusMessage    string    `gorm:"column:status_message"`
	PointsTransferId string    `gorm:"column:points_transfer_id"`
	RequestorData    string    `gorm:"column:requestor_data"`
	ResponderData    string    `gorm:"column:responder_data"`
	TransType        string    `gorm:"trans_type"`
	ExpiredPoint     string    `gorm:"expired_point"`
	AccountId        string    `gorm:"account_id"`
	CreatedAt        time.Time `gorm:"column:created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
	TransactionTime  time.Time `json:"transaction_time"`
}

func (t *TEarning) TableName() string {
	return "public.t_earning"
}
