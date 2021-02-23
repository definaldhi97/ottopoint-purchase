package host

import (
	"crypto/tls"
	"net/http"
	"strconv"
	"strings"
	"time"

	headermodels "ottopoint-purchase/models"

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
	if dch := utils.GetEnv("HTTP_DEBUG_CLIENT_SIGNATURE", "true"); strings.EqualFold(dch, "true") || strings.EqualFold(dch, "false") {
		debugClientHTTP, _ = strconv.ParseBool(strings.ToLower(dch))
	}
	timeout = utils.GetEnv("HTTP_TIMEOUT_SIGNATURE", "60s")
	retrybad = 1
	if rb := utils.GetEnv("HTTP_RETRY_BAD_SIGNATURE", "1"); strings.TrimSpace(rb) != "" {
		if val, err := strconv.Atoi(rb); err == nil {
			retrybad = val
		}
	}

}

func HTTPxFormPostWithHeader(url, key string, data interface{}, header headermodels.RequestHeader) ([]byte, error) {
	// logrus.Info("Token :", dataToken)
	request := gorequest.New()
	request.SetDebug(debugClientHTTP)
	timeout, _ := time.ParseDuration(timeout)
	//_ := errors.New("Connection Problem")
	if strings.HasPrefix(url, "https") {
		request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	reqagent := request.Post(url)

	reqagent.Header.Set("Content-Type", "application/json")
	reqagent.Header.Set("DeviceId", header.DeviceID)
	reqagent.Header.Set("InstitutionId", header.InstitutionID)
	reqagent.Header.Set("Geolocation", header.Geolocation)
	reqagent.Header.Set("ChannelId", header.ChannelID)
	reqagent.Header.Set("AppsId", header.AppsID)
	reqagent.Header.Set("Timestamp", header.Timestamp)
	reqagent.Header.Set("Signature", header.Signature)

	_, body, errs := reqagent.
		Send(data).
		Timeout(timeout).
		Retry(retrybad, time.Second, http.StatusInternalServerError).
		End()

	// healthCheckData, _ := json.Marshal(hcredismodels.ServiceHealthCheckRedis{
	// 	StatusCode: resp.StatusCode,
	// 	UpdatedAt:  time.Now().UTC(),
	// })

	// go redis.SaveRedis(key, healthCheckData)
	if errs != nil {
		logs.Error("Error Sending ", errs)
		return nil, errs[0]
	}
	return []byte(body), nil
}
