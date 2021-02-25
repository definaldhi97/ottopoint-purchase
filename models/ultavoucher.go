package models

type UltraVoucherResp struct {
	// Success int    `json:"success"`
	// Failed  int    `json:"failed"`
	// Total   int    `json:"total"`
	// Voucher string `json:"voucher"`

	Code    string `json:"code"`
	Msg     string `json:"msg"`
	Success int    `json:"success"`
	Failed  int    `json:"failed"`
	Pending int    `json:"pending"`
}
