package models

type RespClearCacheBalance struct {
	ResponseCode string      `json:"response_code"`
	Messages     string      `json:"messages"`
	Data         interface{} `json:"data"`
}
