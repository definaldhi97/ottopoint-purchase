package models

// Voucher Detail
type VoucherDetailResp struct {
	Name       string `json:"name"`
	CampaignID string `json:"campaignId"`
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
