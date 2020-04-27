package models

type UltraVoucherResp struct {
	Success int    `json:"success"`
	Failed  int    `json:"failed"`
	Total   int    `json:"total"`
	Voucher string `json:"voucher"`
}

type ParamUV struct {
	Nama    string
	Email   string
	Phone   string
	Expired string
}
