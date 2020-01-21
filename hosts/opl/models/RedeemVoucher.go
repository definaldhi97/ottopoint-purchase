package models

type BuyVocuherResp struct {
	Coupons []GetCoupons `json:"coupons"`
}

type GetCoupons struct {
	Code string `json:"code"`
}
