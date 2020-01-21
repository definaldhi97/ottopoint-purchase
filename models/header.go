package models

type RequestHeader struct {
	DeviceID      string
	InstitutionID string
	Geolocation   string
	ChannelID     string
	AppsID        string
	Timestamp     string
	Authorization string
	Signature     string
}

// Response
type ResponseToken struct {
	Data struct {
		AccountNumber string `json:"accountNumber"`
	} `json:"data"`
	Meta struct {
		Status  bool   `json:"status"`
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"meta"`
}
