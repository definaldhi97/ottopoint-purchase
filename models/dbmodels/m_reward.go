package dbmodels

import "time"

type MRewardModel struct {
	RewardID                   string    `gorm:"id" json:"id"`
	Name                       string    `gorm:"column:name" json:"name"`
	Reward                     string    `gorm:"column:reward" json:"reward"`
	ProductID                  string    `gorm:"column:m_product_id" json:"ProductID"`
	Categories                 string    `gorm:"column:categories" json:"categories"`
	MoreInformationLink        string    `gorm:"column:more_information_link" json:"moreInformationLink"`
	PushNotificationText       string    `gorm:"column:push_notification_text" json:"pushNotificationText"`
	ProductValue               int       `gorm:"column:product_value" json:"productValue"`
	CostinPoints               int       `gorm:"column:cost_in_points" json:"costinPoints"`
	Levels                     int       `gorm:"column:levels" json:"levels"`
	Segments                   int       `gorm:"column:segments" json:"segments"`
	Unlimited                  bool      `gorm:"column:unlimited" json:"unlimited"`
	SingleCoupon               bool      `gorm:"column:single_coupon" json:"singleCoupon"`
	UsageLimit                 int       `gorm:"column:usage_limit" json:"usageLimit"`
	LimitPeruser               int       `gorm:"column:limit_per_user" json:"limitPerUser"`
	RewardCodes                string    `gorm:"column:reward_odes" json:"rewardCodes"`
	SupplierCost               int       `gorm:"column:supplier_cost" json:"supplierCost"`
	Tax                        string    `gorm:"column:tax" json:"tax"`
	taxPriceValue              int       `gorm:"column:tax_price_value" json:"taxPriceValue"`
	Labels                     string    `gorm:"column:labels" json:"labels"`
	DaysInactive               int       `gorm:"column:days_inactive" json:"daysInactive"`
	TransactionPercentageValue int       `gorm:"column:transaction_percentage_value" json:"transactionPercentageValue"`
	Featured                   bool      `gorm:"column:featured" json:"Featured"`
	Public                     bool      `gorm:"column:public" json:"public"`
	ConnectType                string    `gorm:"column:connect_type" json:"connectType"`
	CashbackProvider           string    `gorm:"column:cashback_provider" json:"cashbackProvider"`
	FulfillmentTracking        bool      `gorm:"column:fulfillment_tracking" json:"fulfillmentTracking"`
	AuditStatusTracking        bool      `gorm:"column:auditstatustracking" json:"auditStatusTracking"`
	ActivityAllTimeActive      bool      `gorm:"column:activity_all_time_active" json:"activityAllTimeActive"`
	ActivityActiveFrom         time.Time `gorm:"column:activity_active_from" json:"activityActiveFrom"`
	ActivityActiveTo           time.Time `gorm:"column:activity_active_to" json:"activityActiveTo"`
	VisibilityAllTimeVisible   bool      `gorm:"column:visibility_all_time_visible" json:"visibilityAllTimeVisible"`
	VisibilityVisibleFrom      time.Time `gorm:"column:visibility_visible_from" json:"visibilityVisibleFrom"`
	visibilityVisibleTo        time.Time `gorm:"column:visibility_visible_to" json:"visibilityVisibleTo"`
	IsActive                   bool      `gorm:"column:is_active" json:"isActive"`
	CreatedAt                  time.Time `gorm:"column:created_at" json:"createdAt"`
	CreatedBy                  string    `gorm:"column:created_by" json:"createdBy"`
	UpdatedAt                  time.Time `gorm:"column:updated_at" json:"updatedAt"`
	updatedBy                  string    `gorm:"column:updated_by" json:"updatedBy"`
}

func (t *MRewardModel) TableName() string {
	return "product.m_reward"
}
