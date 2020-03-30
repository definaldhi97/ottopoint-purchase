package models

type BuyVocuherResp struct {
	Coupons []GetCoupons `json:"coupons"`
	Error   string       `json:"error"`
}

type GetCoupons struct {
	Id   string `json:"id"`
	Code string `json:"code"`
}
