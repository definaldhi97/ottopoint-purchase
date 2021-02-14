package host

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/parnurzeal/gorequest"
)

type EnvHttp struct {
	DebugClient bool   `envconfig:"DEBUGCLIENT_LP" default:"true"`
	Timeout     string `envconfig:"TIMEOUT_LP" default:"60s"`
	RetryBad    int    `envconfig:"RETRY_LP" default:"1"`
}

var (
	envHttp         EnvHttp
	debugClientHTTP bool
	timeout         string
	retrybad        int
)

func init() {

	err := envconfig.Process("HTTP", &envHttp)
	if err != nil {
		fmt.Println("Failed to get HTTP env:", err)
	}

}

func HTTP_POST_LP(url string, json interface{}) ([]byte, error) {
	request := gorequest.New()
	request.SetDebug(debugClientHTTP)
	timeout, _ := time.ParseDuration(timeout)
	if strings.HasPrefix(url, "https") {
		request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	reqagent := request.Post(url)
	// reqagent.Header = header
	_, body, errs := reqagent.
		Timeout(timeout).
		Retry(retrybad, time.Second, http.StatusInternalServerError).
		End()
	if errs != nil {
		return []byte(body), errs[0]
	}
	return []byte(body), nil
}
