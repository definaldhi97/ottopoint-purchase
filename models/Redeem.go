package models

// Request
type RedeemReq struct {
	CampaignID string `json:"campaign"`
	// Phone      string `json:"phone"` // sementara
	Jumlah int `json:"jumlah"`
}

// Response
type RedeemResp struct {
	CodeVoucher []CouponsRedeem `json:"CodeVoucher"`
}

type CouponsRedeem struct {
	Voucher string `json:"voucher"`
	Code    string `json:"code"`
	ID      string `json:"id"`
}

type Redeem struct {
	Voucher []RedeemVoucher
}
type RedeemVoucher struct {
	Voucher string `json:"voucher"`
}

type CountVoucherRedeemed struct {
	Count int `gorm:"column:count" json:"count"`
}
