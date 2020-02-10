package dbmodels

import "time"

type DeductTransaction struct {
	ID            string    `gorm:"column:id"`
	TrxID         string    `gorm:"column:trx_id"`
	AccountID     string    `gorm:"column:account_id"`
	CustomerID    string    `gorm:"column:customer_id"`
	InstitutionID string    `gorm:"column:institution_id"`
	DeductType    int       `gorm:"column:deduct_type"`
	ProductCode   string    `gorm:"column:product_code"`
	ProductName   string    `gorm:"column:product_name"`
	Amount        int       `gorm:"column:amount"`
	Point         int       `gorm:"column:point"`
	Status        string    `gorm:"column:status"`
	CreatedAt     time.Time `gorm:"column:created_at"`
}

func (t *DeductTransaction) TableName() string {
	return "public.t_deduct_transaction"
}
