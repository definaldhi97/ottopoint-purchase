package models

type BuyVocuherResp struct {
	Coupons []GetCoupons `json:"coupons"`
	Error   string       `json:"error"`
	Code    int          `json:"code"`
	Message string       `json:"message"`
}

type GetCoupons struct {
	Id   string `json:"id"`
	Code string `json:"code"`
}
