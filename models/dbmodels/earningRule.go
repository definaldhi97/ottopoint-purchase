package dbmodels

import (
	"time"
)

type MEarningRule struct {
	ID               int       `gorm:"id";pk`
	Code             string    `gorm:"code"`
	Description      string    `gorm:"description"`
	MInstitutionID   int       `gorm:"m_institution_id"`
	Dtype            string    `gorm:"dtype"`
	EventName        string    `gorm:"event_name"`
	Active           bool      `gorm:"active"`
	Levels           string    `gorm:"levels"`
	AllTimeActive    bool      `gorm:"all_time_active"`
	StartAt          time.Time `gorm:"start_at"`
	EndAt            time.Time `gorm:"end_at"`
	SkuIds           string    `gorm:"sku_ids"`
	PointsAmount     float64   `gorm:"points_amount"`
	PointValue       float64   `gorm:"point_value"`
	LimitActive      bool      `gorm:"limit_active"`
	LimitLimit       int       `gorm:"limit_limit"`
	LimitPeriod      string    `gorm:"limit_period"`
	RewardCampaignID string    `gorm:"reward_campaign_id"`
	MinOrderValue    int       `gorm:"min_order_value"`
	ExcludedSkus     string    `gorm:"excluded_skus"`
	Multiplier       float64   `gorm:"multiplier"`
	Latitude         float64   `gorm:"latitude"`
	Longitude        float64   `gorm:"longitude"`
	Radius           float32   `gorm:"radisu"`
	RewardType       string    `gorm:"reward_type"`
	CreatedAt        time.Time `gorm:"created_at"`
	UpdatedAt        time.Time `gorm:"updated_at"`
	CreatedBy        string    `gorm:"created_by"`
	UpdatedBy        string    `gorm:"updated_by"`
}

func (t *MEarningRule) TableName() string {
	return "public.m_earning_rule"
}
