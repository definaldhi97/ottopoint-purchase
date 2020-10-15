package models

var ResponseCode = map[string]string{
	"00": "Success",
	"01": "Pending",
	"20": "Wrong Number / Number Blocked / Number Expired",
	"21": "Product Issue",
	"22": "Duplicate Transaction",
	"23": "Connection Timeout",
	"98": "Order Canceled by ops",
	"99": "General Error",
}

func GetErrorMsg(responseCode string) string {
	return ResponseCode[responseCode]
}

type SepulsaRes struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	Success int    `json:"success"`
	Failed  int    `json:"failed"`
	Pending int    `json:"pending"`
}
