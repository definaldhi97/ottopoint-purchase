package dbmodels

import "time"

type UserMyVocuher struct {
	ID            string    `gorm:"column:id"`
	InstitutionID string    `gorm:"column:institution_id"`
	CouponID      string    `gorm:"column:coupon_id"`    // opl
	VoucherCode   string    `gorm:"column:voucher_code"` // uv
	Phone         string    `gorm:"column:phone"`
	CreatedAT     time.Time `gorm:"column:created_at"`
}

func (t *UserMyVocuher) TableName() string {
	return "public.user_myvoucher"
}
