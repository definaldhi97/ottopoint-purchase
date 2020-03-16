package models

type NotifRequest struct {
	AccountNumber    string `json:"accountNumber"`
	Title            string `json:"title"`
	Message          string `json:"message"`
	NotificationType int    `json:"notificationType"`
}

type NotifResp struct {
	RC      string `json:"rc"`
	Message string `json:"msg"`
}
