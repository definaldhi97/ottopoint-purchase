package models

// Response ..
type Response struct {
	Meta MetaData    `json:"meta"`
	Data interface{} `json:"data,omitempty"`
}

// MetaData ..
type MetaData struct {
	Status  bool   `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// MappingErrorCodes models
type MappingErrorCodes struct {
	Key     string           `json:"key"`
	Content ContentErrorCode `json:"content"`
}

// ContentErrorCode models
type ContentErrorCode struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponseData struct {
	ResponseCode string      `json:"responseCode"`
	ResponseDesc string      `json:"responseDesc"`
	Data         interface{} `json:"data"`
}
