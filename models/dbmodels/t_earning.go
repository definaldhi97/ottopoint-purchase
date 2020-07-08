package dbmodels

import "time"

type TEarning struct {
	ID              string    `gorm:"column:id"`
	EarningRule     string    `gorm:"column:earning_earning"`
	EarningRuleAdd  string    `gorm:"column:earning_earning_add"`
	PartnerId       string    `gorm:"column:partner_id"`
	ReferenceId     string    `gorm:"column:reference_id"`
	TransactionId   string    `grom:"column:transaction_id"`
	ProductCode     string    `gorm:"column:product_code"`
	ProductName     string    `gorm:"column:product_name"`
	AccountNumber1  string    `gorm:"column:account_number1"`
	AccountNumber2  string    `gorm:"column:account_number2"`
	Amount          int64     `gorm:"column:amount"`
	Point           int64     `gorm:"column:point"`
	Remark          string    `gorm:"column:remark"`
	StatusCode      string    `gorm:"column:status_code"`
	StatusMessage   string    `gorm:"column:status_message"`
	PointTransferId string    `gorm:"column:point_transfer_id"`
	RequestorData   string    `gorm:"column:requestor_data"`
	ResponderData   string    `gorm:"column:responder_data"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
}

func (t *TEarning) TableName() string {
	return "public.t_earning"
}
