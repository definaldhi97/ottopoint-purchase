package dbmodels

import "time"

type TPayment struct {
	ID             string     `gorm:"id"`
	TSpendingID    string     `gorm:"t_spending_id"`
	ExternalReffId string     `gorm:"external_reff_id"`
	TransType      string     `gorm:"trans_type"`
	Value          int64      `grom:"value"`
	ValueType      string     `gorm:"value_type"`
	Status         string     `gorm:"status"`
	ResponseRc     string     `gorm:"response_rc"`
	ResponseRd     string     `gorm:"response_rd"`
	CreatedBy      string     `gorm:"created_by"`
	UpdatedBy      string     `gorm:"updated_by"`
	CreatedAt      time.Time  `gorm:"created_at"`
	UpdatedAt      *time.Time `gorm:"updated_at"`
}

func (t *TPayment) TableName() string {
	return "public.t_payment"
}
