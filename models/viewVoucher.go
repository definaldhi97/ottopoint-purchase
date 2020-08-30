package models

type ViewVocuherVidio struct {
	VoucherName string `json:"voucherName"`
	ExpiredDate string `json:"expiredDate"`
	VoucherCode string `json:"voucherCode"`
	ImageUrl    string `json:"imageUrl"`
}
