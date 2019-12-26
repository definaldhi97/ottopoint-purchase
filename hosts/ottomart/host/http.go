package host

import (
	"net/http"
	"ottopoint-purchase/models"
	"strconv"
	"strings"

	"time"

	"github.com/astaxie/beego/logs"
	"github.com/parnurzeal/gorequest"
	ODU "ottodigital.id/library/utils"
)

var (
	debugClientHTTP bool
	instuition      string
	timeout         string
	retrybad        int
)

func init() {
	debugClientHTTP = true //defaultValue
	if dch := ODU.GetEnv("HTTP_DEBUG_CLIENT", "true"); strings.EqualFold(dch, "true") || strings.EqualFold(dch, "false") {
		debugClientHTTP, _ = strconv.ParseBool(strings.ToLower(dch))
	}
	timeout = ODU.GetEnv("HTTP_TIMEOUT", "60s")
	retrybad = 1
	if rb := ODU.GetEnv("HTTP_RETRY_BAD", "1"); strings.TrimSpace(rb) != "" {
		if val, err := strconv.Atoi(rb); err == nil {
			retrybad = val
		}
	}

}

func HTTPxFormOttomart(url string, header models.RequestHeader, jsondata interface{}) ([]byte, error) {
	request := gorequest.New()
	request.SetDebug(debugClientHTTP)
	timeout, _ := time.ParseDuration(timeout)
	//_ := errors.New("Connection Problem")
	reqagent := request.Post(url)
	reqagent.Header.Set("Content-Type", "application/json")
	reqagent.Header.Set("Authorization", header.Authorization)
	reqagent.Header.Set("Device-Id", header.DeviceID)
	_, body, errs := reqagent.
		Send(jsondata).
		Timeout(timeout).
		Retry(retrybad, time.Second, http.StatusInternalServerError).
		End()
	if errs != nil {
		logs.Error("Error Sending ", errs)
		return nil, errs[0]
	}
	return []byte(body), nil
}

func HTTPxFormOttomartNotif(url string, jsondata interface{}, header models.RequestHeader) ([]byte, error) {
	request := gorequest.New()
	request.SetDebug(debugClientHTTP)
	timeout, _ := time.ParseDuration(timeout)
	//_ := errors.New("Connection Problem")
	reqagent := request.Post(url)
	reqagent.Header.Set("Content-Type", "application/json")
	reqagent.Header.Set("Authorization", header.Authorization)
	reqagent.Header.Set("Device-Id", header.DeviceID)
	_, body, errs := reqagent.
		Send(jsondata).
		Timeout(timeout).
		Retry(retrybad, time.Second, http.StatusInternalServerError).
		End()
	if errs != nil {
		logs.Error("Error Sending ", errs)
		return nil, errs[0]
	}
	return []byte(body), nil
}

func HTTPxFormOttomartVoucher(url string, header models.RequestHeader, jsondata interface{}) ([]byte, error) {
	request := gorequest.New()
	request.SetDebug(debugClientHTTP)
	timeout, _ := time.ParseDuration(timeout)
	//_ := errors.New("Connection Problem")
	reqagent := request.Post(url)
	reqagent.Header.Set("Content-Type", "application/json")
	reqagent.Header.Set("Authorization", header.Authorization)
	reqagent.Header.Set("Device-Id", header.DeviceID)
	_, body, errs := reqagent.
		Send(jsondata).
		Timeout(timeout).
		Retry(retrybad, time.Second, http.StatusInternalServerError).
		End()
	if errs != nil {
		logs.Error("Error Sending ", errs)
		return nil, errs[0]
	}
	return []byte(body), nil
}

func HTTPxFormTokenOttomart(url string, header models.RequestHeader) ([]byte, error) {
	request := gorequest.New()
	request.SetDebug(debugClientHTTP)
	timeout, _ := time.ParseDuration(timeout)
	//_ := errors.New("Connection Problem")
	reqagent := request.Get(url)
	// reqagent.Header.Set("Content-Type", "application/json")
	reqagent.Header.Set("Authorization", header.Authorization)
	reqagent.Header.Set("Device-Id", header.DeviceID)
	_, body, errs := reqagent.
		// Send(jsondata).
		Timeout(timeout).
		Retry(retrybad, time.Second, http.StatusInternalServerError).
		End()
	if errs != nil {
		logs.Error("Error Sending ", errs)
		return nil, errs[0]
	}
	return []byte(body), nil
}
