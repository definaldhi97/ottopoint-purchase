package host

import (
	"crypto/tls"
	"net/http"
	"ottopoint-purchase/hosts/voucher_aggregator/models"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/parnurzeal/gorequest"
	ODU "ottodigital.id/library/utils"
)

var (
	debugClientHTTP bool
	timeout         string
	retrybad        int
)

func init() {
	debugClientHTTP = true //defaultValue
	if dch := ODU.GetEnv("HTTP_DEBUG_CLIENT_VOUCHERAG", "true"); strings.EqualFold(dch, "true") || strings.EqualFold(dch, "false") {
		debugClientHTTP, _ = strconv.ParseBool(strings.ToLower(dch))
	}
	timeout = ODU.GetEnv("HTTP_TIMEOUT_VOUCHERAG", "60s")
	retrybad = 1
	if rb := ODU.GetEnv("HTTP_RETRY_BAD_VOUCHERAG", "1"); strings.TrimSpace(rb) != "" {
		if val, err := strconv.Atoi(rb); err == nil {
			retrybad = val
		}
	}

}

func HTTPxFormPostVoucherAg(url string, head models.HeaderHTTP, jsonReq interface{}) ([]byte, error) {

	request := gorequest.New()

	request.SetDebug(debugClientHTTP)

	timeout, _ := time.ParseDuration(timeout)
	if strings.HasPrefix(url, "https") {
		request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	reqagent := request.Post(url)

	reqagent.Header.Set("Content-Type", "application/json")
	reqagent.Header.Set("InstitutionId", head.Institution)
	reqagent.Header.Set("DeviceId", head.DeviceID)
	reqagent.Header.Set("Geolocation", head.Geolocation)
	reqagent.Header.Set("ChannelId", "H2H")
	reqagent.Header.Set("Signature", head.Signature)
	reqagent.Header.Set("AppsId", head.AppsID)
	reqagent.Header.Set("Timestamp", head.Timestamp)

	_, body, errs := reqagent.
		Send(jsonReq).
		Timeout(timeout).
		Retry(retrybad, time.Second, http.StatusInternalServerError).
		End()
	if errs != nil {
		logs.Error("Error Sending ", errs)
		return nil, errs[0]
	}
	return []byte(body), nil

}

func HTTPxFormGetVoucherAg(url string, head models.HeaderHTTP) ([]byte, error) {

	request := gorequest.New()

	request.SetDebug(debugClientHTTP)

	timeout, _ := time.ParseDuration(timeout)
	if strings.HasPrefix(url, "https") {
		request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	reqagent := request.Get(url)

	reqagent.Header.Set("Content-Type", "application/json")
	reqagent.Header.Set("InstitutionId", head.Institution)
	reqagent.Header.Set("DeviceId", head.DeviceID)
	reqagent.Header.Set("Geolocation", head.Geolocation)
	reqagent.Header.Set("ChannelId", "H2H")
	reqagent.Header.Set("Signature", head.Signature)
	reqagent.Header.Set("AppsId", head.AppsID)
	reqagent.Header.Set("Timestamp", head.Timestamp)

	_, body, errs := reqagent.
		Timeout(timeout).
		Retry(retrybad, time.Second, http.StatusInternalServerError).
		End()
	if errs != nil {
		logs.Error("Error Sending ", errs)
		return nil, errs[0]
	}
	return []byte(body), nil

}
