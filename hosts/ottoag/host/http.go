package host

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"ottopoint-purchase/utils"

	"github.com/parnurzeal/gorequest"
)

var (
	debugClientHTTP bool
	timeout         string
	retrybad        int
)

func init() {
	debugClientHTTP = true //defaultValue
	if dch := utils.GetEnv("HTTP_DEBUG_CLIENT_OTTOAG", "true"); strings.EqualFold(dch, "true") || strings.EqualFold(dch, "false") {
		debugClientHTTP, _ = strconv.ParseBool(strings.ToLower(dch))
	}
	timeout = utils.GetEnv("HTTP_TIMEOUT_OTTOAG", "60s")
	retrybad = 1
	if rb := utils.GetEnv("HTTP_RETRY_BAD_OTTOAG", "1"); strings.TrimSpace(rb) != "" {
		if val, err := strconv.Atoi(rb); err == nil {
			retrybad = val
		}
	}

}

// HTTPPostWithHeader func
func HTTPPostWithHeader(url string, jsondata interface{}, header http.Header) ([]byte, error) {
	request := gorequest.New()
	request.SetDebug(debugClientHTTP)
	timeout, _ := time.ParseDuration(timeout)
	//_ := errors.New("Connection Problem")
	// if url[:5] == "https" {
	// 	request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	// }
	reqagent := request.Post(url)
	reqagent.Header = header
	_, body, errs := reqagent.
		Send(jsondata).
		Timeout(timeout).
		Retry(retrybad, time.Second, http.StatusInternalServerError).
		End()
	if errs != nil {
		return []byte(body), errs[0]
	}
	return []byte(body), nil
}
