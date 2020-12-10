package dbmodels

import "time"

type MProductBrand struct {
	Id        string    `gorm:"primary_key" json:"id"`
	Code      string    `gorm:"column:code" json:"code"`
	Name      string    `gorm:"column:name" json:"name"`
	Path      string    `gorm:"column:path" json:"path"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	CreatedBy string    `gorm:"column:created_by" json:"created_by"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	UpdatedBy string    `gorm:"column:updated_by" json:"updated_by"`
	SortOrder int       `gorm:"column:sort_order" json:"sort_order"`
	IsActive  bool      `gorm:"column:is_active" json:"is_active"`
	IsWidget  bool      `gorm:"column:is_widget" json:"is_widget"`
}

func (t *MProductBrand) TableName() string {
	return "product.m_product_brand"
}
