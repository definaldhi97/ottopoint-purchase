package host

import (
	"encoding/json"
	"ottopoint-purchase/hosts/ottomart/models"

	"github.com/astaxie/beego/logs"
	ODU "ottodigital.id/library/utils"
)

var (
	host string
	name string

	endpointNotifAndInbox string

	healthCheckKey string
)

func init() {
	host = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_OTTOMART", "http://13.228.25.85:8186/v1.0")
	name = ODU.GetEnv("OTTOPOINT_PURCHASE_NAME_OTTOMART", "OTTOMART")

	endpointNotifAndInbox = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_NOTIF_INBOX", "/notifications")

	healthCheckKey = ODU.GetEnv("OTTOPOINT_PURCHASE_KEY_HEALTHCHECK_OTTOMART", "OTTOPOINT-PURCHASE:OTTOPOINT-OTTOMART")
}

// Send Notif & Inbox
func NotifAndInbox(req models.NotifRequest) (*models.NotifResp, error) {
	var resp models.NotifResp

	logs.Info("[Package Host OTTOMART]-[NotifAndInbox]")

	urlSvr := host + endpointNotifAndInbox

	data, err := HTTPxFormOTTOMART(urlSvr, req, healthCheckKey)
	if err != nil {
		logs.Error("Check error", err.Error())

		return &resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logs.Error("Failed to unmarshaling response NotifAndInbox from ottomart ", err.Error())

		return &resp, err
	}

	return &resp, nil
}
