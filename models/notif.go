package models

// Request Notif
type NotifReq struct {
	CollapseKey string `json:"collapse_key"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	Target      string `json:"target"`
	Phone       string `json:"phone"`
	Rc          string `json:"rc"`
}

// Response Notif
type NotifResp struct {
	Data Notifs      `json:"data"`
	Meta interface{} `json:"meta"`
}

type Notifs struct {
	Nama          string `json:"nama"`
	AccountNumber string `json:"account_number"`
	Body          string `json:"body"`
}

type NotifPubreq struct {
	Type          string `json:"type"`
	AccountNumber string `json:"accountNumber"`
	Institution   string `json:"institution"`
	Point         int    `json:"point"`
	Product       string `json:"product"`
}
