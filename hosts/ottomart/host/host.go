package host

import (
	"encoding/json"
	"net/http"
	"ottopoint-purchase/hosts/ottomart/models"

	https "ottopoint-purchase/hosts"
	"ottopoint-purchase/utils"

	"github.com/astaxie/beego/logs"
)

var (
	host string
	name string

	endpointNotifAndInbox string

	// healthCheckKey string
)

func init() {
	host = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_OTTOMART", "http://13.228.25.85:8186/v1.0")
	name = utils.GetEnv("OTTOPOINT_PURCHASE_NAME_OTTOMART", "OTTOMART")

	endpointNotifAndInbox = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_NOTIF_INBOX", "/notifications")

	// healthCheckKey = utils.GetEnv("OTTOPOINT_PURCHASE_KEY_HEALTHCHECK_OTTOMART", "OTTOPOINT-PURCHASE:OTTOPOINT-OTTOMART")
}

// Send Notif & Inbox
func NotifAndInbox(req models.NotifRequest) (*models.NotifResp, error) {
	var resp models.NotifResp

	logs.Info("[Package Host OTTOMART]-[NotifAndInbox]")

	header := make(http.Header)
	header.Set("Content-Type", "application/json")

	urlSvr := host + endpointNotifAndInbox

	data, err := https.HTTPxPOSTwithRequest(urlSvr, req, header)
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
