package dbmodels

import "time"

type Configs struct {
	ID                int       `gorm:"column:id;pk"`
	TransaksiPPOB     float64   `gorm:"column:transaksi_ppob"`
	TransaksiPayQR    float64   `gorm:"column:transaksi_pay_qr"`
	TransaksiMerchant float64   `gorm:"column:transaksi_merchant"`
	LimitTransaksi    int       `gorm:"column:limit_transaksi"`
	MinimalTransaksi  int64     `gorm:"column:minimal_transaksi"`
	MemberID          string    `gorm:"column:member_id"`
	CreatedAT         time.Time `gorm:"column:created_at"`
	UpdatedAT         time.Time `gorm:"column:updated_at"`
}

func (t *Configs) TableName() string {
	return "public.configs"
}
