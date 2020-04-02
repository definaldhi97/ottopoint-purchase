package models

// CouponCustomer ..
type CouponVoucherCustomerResp struct {
	Points        int    `json:"points"`
	Code          string `json:"code"`
	CampaignID    string `json:"campaignId"`
	CouponID      string `json:"couponId"`
	AccountNumber string `json:"accountNumber"`
	CustID        string `json:"custId"`
	Date          string `json:"date"`
	// Name          string `json:"name"`
}
