package dbmodels

import "time"

type MProduct struct {
	ID              string    `gorm:"primary_key", json:"id"`
	Code            string    `gorm:"code" json:"code"`
	Name            string    `gorm:"name" json:"name"`
	Denom           string    `gorm:"denom" json:"denom"`
	RewardType      string    `gorm:"reward_type" json:"reward_type"`
	MProductBrandID string    `gorm:"m_product_brand_id" json:"m_product_brand_id"`
	MVendorID       string    `gorm:"m_vendor_id" json:"m_vendor_id"`
	IsActive        bool      `gorm:"is_active" json:"is_active"`
	CreatedAt       time.Time `gorm:"created_at" json:"created_at"`
	CreatedBy       string    `gorm:"created_by" json:"created_by"`
	UpdatedAt       time.Time `gorm:"updated_at" json:"updated_at"`
	UpdatedBy       string    `gorm:"updated_by" json:"updated_by"`
}

func (t *MProduct) TableName() string {
	return "product.m_product"
}
