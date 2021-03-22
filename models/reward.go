package models

import (
	"time"

	"github.com/lib/pq"
)

type VoucherDetailsManagement1 struct {
	RewardID            string         `gorm:"column:reward_id" json:"rewardID"`
	VoucherName         string         `gorm:"column:name" json:"voucherName"`
	CostPoints          float64        `gorm:"column:cost_in_points" json:"costPoints"`
	UsageLimit          int            `gorm:"column:usage_limit" json:"usageLimit"`
	BrandName           string         `gorm:"column:brand_name" json:"brandName"`
	ActivityActiveFrom  time.Time      `gorm:"column:activity_active_from" json:"activityActiveFrom"`
	ActivityActiveTo    string         `gorm:"column:activity_active_to" json:"activityActiveTo"`
	CategoriesID        string         `gorm:"column:categories_id" json:"categoriesID"`
	CodeSuplier         string         `gorm:"column:code_suplier" json:"codeSuplier"`
	RewardCodes         string         `gorm:"column:reward_codes" json:"rewardCodes"`
	ExternalProductCode string         `gorm:"column:external_code" json:"productCodeExternal"`
	InternalProductCode string         `gorm:"column:internal_code" json:"productCodeInternal"`
	ProductID           string         `gorm:"column:m_product_id" json:"productID"`
	Fields              pq.StringArray `gorm:"m_product_brand_id"`
	// Fields string `gorm:"m_product_brand_id"`
}

type VoucherDetailsManagement struct {
	RewardID            string         `gorm:"column:reward_id" json:"rewardID"`
	VoucherName         string         `gorm:"column:name" json:"voucherName"`
	CostPoints          float64        `gorm:"column:cost_in_points" json:"costPoints"`
	UsageLimit          int            `gorm:"column:usage_limit" json:"usageLimit"`
	BrandName           string         `gorm:"column:brand_name" json:"brandName"`
	ActivityActiveFrom  time.Time      `gorm:"column:activity_active_from" json:"activityActiveFrom"`
	ActivityActiveTo    string         `gorm:"column:activity_active_to" json:"activityActiveTo"`
	CategoriesID        []string       `gorm:"column:categories_id" json:"categoriesID"`
	CodeSuplier         string         `gorm:"column:code_suplier" json:"codeSuplier"`
	RewardCodes         string         `gorm:"column:reward_codes" json:"rewardCodes"`
	ExternalProductCode string         `gorm:"column:external_code" json:"productCodeExternal"`
	InternalProductCode string         `gorm:"column:internal_code" json:"productCodeInternal"`
	ProductID           string         `gorm:"column:m_product_id" json:"productID"`
	Fields              pq.StringArray `gorm:"m_product_brand_id"`
}

// type BuyVocuherResp struct {
// 	Coupons []GetCoupons `json:"coupons"`
// 	Error   string       `json:"error"` // jika error
// 	Code    int          `json:"code"`  // jika error
// 	Message string       `json:"message"`
// }

type SpendingPointVoucher struct {
	Rc              string
	Rd              string
	PointTransferID string
	CouponseVouch   []CouponsVoucher
}
type CouponsVoucher struct {
	CouponsCode string
	CouponsID   string
}
