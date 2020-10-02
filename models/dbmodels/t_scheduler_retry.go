package dbmodels

import "time"

type TSchedulerRetry struct {
	ID            int       `gorm:"id"`
	Code          string    `gorm:"code"`
	TransactionID string    `gorm:"transaction_id"`
	Count         int       `gorm:"count"`
	IsDone        bool      `gorm:"is_done"`
	CreatedAT     time.Time `gorm:"column:created_at"`
	UpdatedAT     time.Time `gorm:"column:updated_at"`
}

func (t *TSchedulerRetry) TableName() string {
	return "public.t_scheduler_retry"
}
