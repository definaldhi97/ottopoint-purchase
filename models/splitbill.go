package models

type SplitBillReq struct {
	Type     string `json:"type"`
	Rrn      string `json:"rrn"`
	AdminFee int64  `json:"admin_fee"`
	Amount   int64  `json:"amount"`
	Point    int    `json:"point"`
	WalletID int    `json:"wallet_id"`
	Pin      string `json:"pin"`
}
