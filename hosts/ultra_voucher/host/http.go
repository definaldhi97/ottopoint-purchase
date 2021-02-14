package host

import (
	"crypto/tls"
	"net/http"
	"strconv"
	"strings"
	"time"

	"ottopoint-purchase/utils"

	"github.com/astaxie/beego/logs"
	"github.com/parnurzeal/gorequest"
)

var (
	debugClientHTTP bool
	timeout         string
	retrybad        int
)

func init() {
	debugClientHTTP = true //defaultValue
	if dch := utils.GetEnv("HTTP_DEBUG_CLIENT_UV", "true"); strings.EqualFold(dch, "true") || strings.EqualFold(dch, "false") {
		debugClientHTTP, _ = strconv.ParseBool(strings.ToLower(dch))
	}
	timeout = utils.GetEnv("HTTP_TIMEOUT_UV", "60s")
	retrybad = 1
	if rb := utils.GetEnv("HTTP_RETRY_BAD_UV", "1"); strings.TrimSpace(rb) != "" {
		if val, err := strconv.Atoi(rb); err == nil {
			retrybad = val
		}
	}

}

func HTTPxFormPostUV(url, InstitutionID string, jsonReq interface{}) ([]byte, error) {
	request := gorequest.New()
	request.SetDebug(debugClientHTTP)
	timeout, _ := time.ParseDuration(timeout)
	//_ := errors.New("Connection Problem")
	if strings.HasPrefix(url, "https") {
		request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	reqagent := request.Post(url)
	if InstitutionID == "" {
		reqagent.Header.Set("Content-Type", "application/json")
	} else {
		reqagent.Header.Set("Content-Type", "application/json")
		reqagent.Header.Set("InstitutionId", InstitutionID)
	}
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

func HTTPxFormGETUV(url, InstitutionReff, InstitutionId string) ([]byte, error) {
	request := gorequest.New()
	request.SetDebug(debugClientHTTP)
	timeout, _ := time.ParseDuration(timeout)
	//_ := errors.New("Connection Problem")
	if strings.HasPrefix(url, "https") {
		request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	reqagent := request.Get(url)
	reqagent.Header.Set("Content-Type", "application/json")
	reqagent.Header.Set("InstitutionId", InstitutionId)
	reqagent.Header.Set("InstitutionRefno", InstitutionReff)

	_, body, errs := reqagent.
		// Send(jsonReq).
		Timeout(timeout).
		Retry(retrybad, time.Second, http.StatusInternalServerError).
		End()
	if errs != nil {
		logs.Error("Error Sending ", errs)
		return nil, errs[0]
	}
	return []byte(body), nil
}
