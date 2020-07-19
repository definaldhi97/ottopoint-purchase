package publisher

import (
	"crypto/tls"
	"net/http"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/parnurzeal/gorequest"
)

var (
	debugClientHTTP bool
	instuition      string
	timeout         string
	retrybad        int
)

func init() {
	debugClientHTTP = beego.AppConfig.DefaultBool("debugClientHTTP", true)
	timeout = beego.AppConfig.DefaultString("timeout", "60s")
	retrybad = beego.AppConfig.DefaultInt("retrybad", 1)
}

func HTTPPostKafka(url string, jsondata interface{}) ([]byte, error) {
	request := gorequest.New()
	request.SetDebug(debugClientHTTP)
	timeout, _ := time.ParseDuration(timeout)
	if strings.HasPrefix(url, "https") {
		request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	reqagent := request.Post(url)
	reqagent.Header.Set("Content-Type", "application/json")
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
