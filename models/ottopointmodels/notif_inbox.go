package ottopointmodels

type NotifOttopointReq struct {
	Phone string `json:"phone"`
	Point int    `json:"point"`
	Rc    string `json:"rc"`
}

type NotifOttopointResp struct {
	Point         int    `json:"point"`
	Nama          string `json:"nama"`
	AccountNumber string `json:"account_number"`
}
