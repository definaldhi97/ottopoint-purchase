package models

// Voucher Detail
type VoucherDetailResp struct {
	Name                          string                   `json:"name"`
	ShortDescription              string                   `json:"shortDescription"`
	ConditionsDescription         string                   `json:"conditionsDescription"`
	BrandDescription              string                   `json:"brandDescription"`
	Levels                        []string                 `json:"levels"`
	Segments                      []string                 `json:"segments"`
	Categories                    []string                 `json:"categories"`
	Coupons                       []CouponDetails          `json:"coupons"`
	BrandIcon                     bool                     `json:"brandIcon"`
	CampaignID                    string                   `json:"campaignId"`
	Reward                        string                   `json:"reward"`
	Active                        bool                     `json:"active"`
	CostInPoints                  int                      `json:"costInPoints"`
	SingleCoupon                  bool                     `json:"singleCoupon"`
	Unlimited                     bool                     `json:"unlimited"`
	Limit                         int                      `json:"limit,omitempty"`
	LimitPerUser                  int                      `json:"limitPerUser,omitempty"`
	CampaignActivity              CampaignActivityDetail   `json:"campaignActivity"`
	CampaignVisibility            CampaignVisibilityDetail `json:"campaignVisibility"`
	RewardValue                   int                      `json:"rewardValue,omitempty"`
	Labels                        []LabelsDetail           `json:"labels"`
	DaysInactive                  int                      `json:"daysInactive,omitempty"`
	DaysValid                     int                      `json:"daysValid,omitempty"`
	Featured                      bool                     `json:"featured"`
	Photos                        []PhotosDetail           `json:"photos"`
	Public                        bool                     `json:"public"`
	FulfillmentTracking           bool                     `json:"fulfillmentTracking"`
	Translations                  []TranslationsDetail     `json:"translations"`
	SegmentNames                  interface{}              `json:"segmentNames"`
	LevelNames                    interface{}              `json:"levelNames"`
	CategoryNames                 interface{}              `json:"categoryNames"`
	UsageLeft                     int                      `json:"usageLeft"`
	VisibleForCustomersCount      int                      `json:"visibleForCustomersCount"`
	UsersWhoUsedThisCampaignCount int                      `json:"usersWhoUsedThisCampaignCount"`
	UsageInstruction              string                   `json:"usageInstruction"`
}

type CouponDetails struct {
	Coupon string `json:"coupon"`
}

type CampaignActivityDetail struct {
	AllTimeActive bool   `json:"allTimeActive"`
	ActiveFrom    string `json:"activeFrom"`
	ActiveTo      string `json:"activeTo"`
}

type CampaignVisibilityDetail struct {
	AllTimeVisible bool   `json:"allTimeVisible"`
	VisibleFrom    string `json:"visibleFrom"`
	VisibleTo      string `json:"visibleTo"`
}

type LabelsDetail struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type TranslationsDetail struct {
	Name                  string `json:"name"`
	ShortDescription      string `json:"shortDescription"`
	ConditionsDescription string `json:"conditionsDescription,omitempty"`
	UsageInstruction      string `json:"usageInstruction,omitempty"`
	BrandDescription      string `json:"brandDescription,omitempty"`
	BrandName             string `json:"brandName,omitempty"`
	ID                    int    `json:"id"`
	Locale                string `json:"locale"`
}
