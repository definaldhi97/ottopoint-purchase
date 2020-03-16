package models

type BuyVocuherResp struct {
	Coupons []GetCoupons `json:"coupons"`
}

type GetCoupons struct {
	Id   string `json:"id"`
	Code string `json:"code"`
}
