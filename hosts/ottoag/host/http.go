package host

import (
	"encoding/json"
	"net/http"
	"ottopoint-purchase/redis"
	"strconv"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"

	hcredismodels "ottodigital.id/library/healthcheck/models/redismodels"

	ODU "ottodigital.id/library/utils"
)

var (
	debugClientHTTP bool
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

// HTTPPostWithHeader func
func HTTPPostWithHeader(url string, jsondata interface{}, header http.Header, key string) ([]byte, error) {
	request := gorequest.New()
	request.SetDebug(debugClientHTTP)
	timeout, _ := time.ParseDuration(timeout)
	//_ := errors.New("Connection Problem")
	// if url[:5] == "https" {
	// 	request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	// }
	reqagent := request.Post(url)
	reqagent.Header = header
	resp, body, errs := reqagent.
		Send(jsondata).
		Timeout(timeout).
		Retry(retrybad, time.Second, http.StatusInternalServerError).
		End()

	healthCheckData, _ := json.Marshal(hcredismodels.ServiceHealthCheckRedis{
		StatusCode: resp.StatusCode,
		UpdatedAt:  time.Now().UTC(),
	})

	go redis.SaveRedis(key, healthCheckData)
	if errs != nil {
		return []byte(body), errs[0]
	}
	return []byte(body), nil
}
