package host

import (
	"crypto/tls"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/parnurzeal/gorequest"
	ODU "ottodigital.id/library/utils"
)

var (
	debugClientHTTP bool = true
	username        string
	password        string
)

func init() {
	debugClientHTTP = true
	username = ODU.GetEnv("OTTOPOINT_PURCHASE_USERNAME_SEPULSA", "ottopoint")
	password = ODU.GetEnv("OTTOPOINT_PURCHASE_PASSWORD_SEPULSA", "DoXptIfK6rOeqxMDa0uxEGXjlzvEVfST")
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

	_, body, errs := reqagent.
		Send(jsonReq).
		End()
	if errs != nil {
		logs.Error("Error Sending ", errs)
		return nil, errs[0]
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
