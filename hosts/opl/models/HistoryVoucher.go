package models

type HistoryVoucherCustomerResponse struct {
	Campaigns []CampaignsDetail `json:"campaigns"`
	Total     int               `json:"total"`
}

type CampaignsDetail struct {
	CanBeUsed      bool            `json:"canBeUsed"`
	PurchaseAt     string          `json:"purchaseAt"`
	UsageDate      string          `json:"usageDate"`
	CostInPoints   int             `json:"costInPoints"`
	CampaignID     string          `json:"campaignId"`
	Used           bool            `json:"used"`
	Campaign       CampaignDetails `json:"campaign"`
	Coupon         CouponDetail    `json:"coupon"`
	Status         string          `json:"status"`
	ActiveTo       string          `json:"activeTo"`
	DeliveryStatus interface{}     `json:"deliveryStatus"`
}

type CampaignDetails struct {
	Name string `json:"name"`
	// BrandIcon                     bool                     `json:"brandIcon"`
	BrandName string `json:"brandName,omitempty"`
	// CampaignID                    string                   `json:"campaignId"`
	// Reward                        string                   `json:"reward"`
	// Active                        bool                     `json:"active"`
	// CostInPoints                  int                      `json:"costInPoints"`
	// SingleCoupon                  bool                     `json:"singleCoupon"`
	// Unlimited                     bool                     `json:"unlimited"`
	// Limit                         int                      `json:"limit"`
	// LimitPerUser                  int                      `json:"limitPerUser"`
	// CampaignActivity              CampaignActivityDetail   `json:"campaignActivity"`
	// CampaignVisibility            CampaignVisibilityDetail `json:"campaignVisibility"`
	// Labels                        []LabelsDetail           `json:"labels"`
	// DaysInactive                  int                      `json:"daysInactive"`
	// DaysValid                     int                      `json:"daysValid"`
	// Featured                      bool                     `json:"featured"`
	Photos []PhotosDetail `json:"photos"`
	// Public                        bool                     `json:"public"`
	// FulfillmentTracking           bool                     `json:"fulfillmentTracking"`
	// Translations                  []TranslationsDetail     `json:"translations"`
	// SegmentNames                  interface{}              `json:"segmentNames"`
	// LevelNames                    interface{}              `json:"levelNames"`
	// CategoryNames                 interface{}              `json:"categoryNames"`
	// UsageLeft                     int                      `json:"usageLeft"`
	// UsageLeftForCustomer          int                      `json:"usageLeftForCustomer"`
	// CanBeBoughtByCustomer         bool                     `json:"canBeBoughtByCustomer"`
	// VisibleForCustomersCount      int                      `json:"visibleForCustomersCount"`
	// UsersWhoUsedThisCampaignCount int                      `json:"usersWhoUsedThisCampaignCount"`
	// BrandDescription              interface{}              `json:"brandDescription"`
	// ShortDescription              interface{}              `json:"shortDescription"`
	// ConditionsDescription         interface{}              `json:"conditionsDescription"`
	// UsageInstruction              interface{}              `json:"usageInstruction"`
	UsageDate string `json:"usageDate"`
}

type CouponDetail struct {
	ID   string `json:"id"`
	Code string `json:"code"`
}

type PhotosDetail struct {
	PhotoID      PhotoIDDetail `json:"photoId"`
	Path         ValueDetail   `json:"path"`
	OriginalName ValueDetail   `json:"originalName"`
	MimeType     ValueDetail   `json:"mimeType"`
}

type PhotoIDDetail struct {
	ID string `json:"id"`
}

type ValueDetail struct {
	Value string `json:"value"`
}
