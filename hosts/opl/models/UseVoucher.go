package models

type CouponVoucherCustomerResp struct {
	Coupons []CouponsVoucherResp `json:"coupons"`
}

type CouponsVoucherResp struct {
	Code       string `json:"code"`
	CouponID   string `json:"couponId"`
	Used       bool   `json:"used"`
	CampaignID string `json:"campaignId"`
	CustomerID string `json:"customerId"`
}
