package host

import (
	"crypto/tls"
	"errors"
	"strconv"
	"strings"

	"ottopoint-purchase/utils"

	"github.com/astaxie/beego/logs"
	"github.com/parnurzeal/gorequest"
	"github.com/prometheus/common/log"
)

var (
	debugClientHTTP bool
	timeout         string
	retrybad        int
	username        string
	password        string
)

func init() {
	debugClientHTTP = true //defaultValue
	if dch := utils.GetEnv("HTTP_DEBUG_CLIENT_SEPULSA", "true"); strings.EqualFold(dch, "true") || strings.EqualFold(dch, "false") {
		debugClientHTTP, _ = strconv.ParseBool(strings.ToLower(dch))
	}
	timeout = utils.GetEnv("HTTP_TIMEOUT_SEPULSA", "60s")
	retrybad = 1
	if rb := utils.GetEnv("HTTP_RETRY_BAD_SEPULSA", "1"); strings.TrimSpace(rb) != "" {
		if val, err := strconv.Atoi(rb); err == nil {
			retrybad = val
		}
	}

	username = utils.GetEnv("OTTOPOINT_PURCHASE_USERNAME_SEPULSA", "ottopoint")
	password = utils.GetEnv("OTTOPOINT_PURCHASE_PASSWORD_SEPULSA", "DoXptIfK6rOeqxMDa0uxEGXjlzvEVfST")
}

func HTTPxFormPostSepulsa(url string, jsonReq interface{}) ([]byte, error) {
	request := gorequest.New()
	request.SetDebug(debugClientHTTP)

	if strings.HasPrefix(url, "https") {
		request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	reqagent := request.Post(url)
	reqagent.Header.Set("User-Agent", "ottopoint")
	reqagent.Header.Set("Accept", "application/json")
	reqagent.Header.Set("Content-Type", "application/json")
	reqagent.SetBasicAuth(username, password)

	resp, body, errs := reqagent.
		Send(jsonReq).
		End()
	if errs != nil {
		logs.Error("Error Sending ", errs)
		return nil, errs[0]
	}

	if resp.StatusCode != 201 {
		log.Errorf("Failed Create Trx: %v", resp.Status)
		return nil, errors.New(resp.Status)
	}

	return []byte(body), nil
}

func HTTPxFormGETSepulsa(url string, jsonReq interface{}) ([]byte, error) {
	request := gorequest.New()
	request.SetDebug(debugClientHTTP)

	if strings.HasPrefix(url, "https") {
		request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	reqagent := request.Get(url)
	reqagent.Header.Set("User-Agent", "ottopoint")
	reqagent.Header.Set("Accept", "application/json")
	reqagent.Header.Set("Content-Type", "application/json")
	reqagent.SetBasicAuth(username, password)

	_, body, errs := reqagent.
		Send(jsonReq).
		End()
	if errs != nil {
		logs.Error("Error Sending ", errs)
		return nil, errs[0]
	}
	return []byte(body), nil
}
