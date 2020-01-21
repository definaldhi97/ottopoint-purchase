package models

type SignatureReq struct {
	Data interface{} `json:"data"`
}

type SignatureResp struct {
	ResponseCode string      `json:"response_code"`
	Message      string      `json:"message"`
	Data         interface{} `json:"data"`
}
