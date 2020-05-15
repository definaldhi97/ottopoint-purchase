package dbmodels

import (
	"time"
)

type MEarningRule struct {
	ID               int       `gorm:"id";pk`
	Code             string    `gorm:"code"`
	Desc             string    `gorm:"description"`
	InstitutionID    int       `gorm:"m_institution_id"`
	Type             string    `gorm:"type"`
	EventName        string    `gorm:"event_name"`
	Active           bool      `gorm:"active"`
	Levels           string    `gorm:"levels"`
	Segements        string    `gorm:"segments"`
	Pos              string    `gorm:"pos"`
	AllTimeActive    bool      `gorm:"all_time_active"`
	StartAt          time.Time `gorm:"start_at"`
	EndAt            time.Time `gorm:"end_at"`
	SkuIds           string    `gorm:"sku_ids"`
	PointsAmount     float64   `gorm:"points_amount"`
	PointValue       float64   `gorm:"point_value"`
	LabelMultipilers string    `gorm:"label_multipilers"`
	LimitActive      bool      `gorm:"limit_active"`
	LimitLimit       int       `gorm:"limit_limit"`
	LimitPeriod      string    `gorm:"limit_period"`
	// RewardCampaignID
	MinOrderValue        float64   `gorm:"min_order_value"`
	ExcludeDeliveryCost  bool      `gorm:"exclude_delivery_cost"`
	excludedSkus         string    `gorm:"excluded_skus"`
	LabelsInclusionType  string    `gorm:"labels_inclusion_type"`
	ExcludedLabels       string    `gorm:"exclude_labels"`
	IncludedLabels       string    `gorm:"included_labels"`
	Multipiler           float64   `gorm:"multipiler"`
	Labels               string    `gorm:"labels"`
	Latitude             int       `gorm:"latitude"`
	Longitude            int       `gorm:"longitude"`
	Radius               float32   `gorm:"radisu"`
	RewardType           string    `gorm:"reward_type"`
	EarningRulePhotoPath string    `gorm:"earning_rule_photo_path"`
	EarningRulePhotoMime string    `gorm:"earning_rule_photo_mime"`
	CreatedAt            time.Time `gorm:"created_at"`
	UpdatedAt            time.Time `gorm:"updated_at"`
	CreatedBy            string    `gorm:"created_by"`
	UpdatedBy            string    `gorm:"updated_by"`
}

func (t *MEarningRule) TableName() string {
	return "public.m_earning_rule"
}
