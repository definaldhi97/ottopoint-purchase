package dbmodels

import "time"

type ErrLogEarning struct {
	ID            int       `gorm:"column:id"`
	AccountNumber string    `gorm:"column:account_number"`
	PartnerId     string    `gorm:"column:partner_id"`
	ReferenceId   string    `gorm:"column:reference_id"`
	RequestorData string    `gorm:"column:requestor_data"`
	StatusCode    string    `gorm:"column:status_code"`
	StatusMessage string    `gorm:"column:status_message"`
	CreatedAt     time.Time `gorm:"column:created_at"`
}

func (t *ErrLogEarning) TableName() string {
	return "logging.err_log_earning"
}
