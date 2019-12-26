package http

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"time"

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

// HTTPGet func
func HTTPGet(url string, header http.Header) ([]byte, error) {
	request := gorequest.New()
	request.SetDebug(debugClientHTTP)
	timeout, _ := time.ParseDuration(timeout)
	//_ := errors.New("Connection Problem")
	// if url[:5] == "https" {
	// 	request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	// }
	reqagent := request.Get(url)
	reqagent.Header = header
	_, body, errs := reqagent.
		Timeout(timeout).
		Retry(retrybad, time.Second, http.StatusInternalServerError).
		End()
	if errs != nil {
		return []byte(body), errs[0]
	}
	return []byte(body), nil
}

// HTTPPost func
func HTTPPost(url string, jsondata interface{}) ([]byte, error) {
	request := gorequest.New()
	request.SetDebug(debugClientHTTP)
	timeout, _ := time.ParseDuration(timeout)
	reqagent := request.Post(url)
	reqagent.Header.Set("Content-Type", "application/json")
	_, body, errs := reqagent.
		Send(jsondata).
		Timeout(timeout).
		Retry(retrybad, time.Second, http.StatusInternalServerError).
		End()
	if errs != nil {
		fmt.Println("[ErrorSending]")
		fmt.Println("[Utils]")
		fmt.Println("[HTTPPost]")
		fmt.Println(fmt.Sprintf("url = %s", url))
		fmt.Println(fmt.Sprintf("Error Sending : %v", errs))
		return nil, errs[0]
	}
	return []byte(body), nil
}

// HTTPPostWithHeader func
func HTTPPostWithHeader(url string, jsondata interface{}, header http.Header) ([]byte, error) {
	request := gorequest.New()
	request.SetDebug(debugClientHTTP)
	timeout, _ := time.ParseDuration(timeout)
	//_ := errors.New("Connection Problem")
	if strings.HasPrefix(url, "https") {
		request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
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

// HTTPPutWithHeader func
func HTTPPutWithHeader(url string, jsondata interface{}, header http.Header) ([]byte, error) {
	request := gorequest.New()
	request.SetDebug(debugClientHTTP)
	timeout, _ := time.ParseDuration(timeout)
	//_ := errors.New("Connection Problem")
	if url[:5] == "https" {
		request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	reqagent := request.Put(url)
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

// HTTPDeleteWithHeader func
func HTTPDeleteWithHeader(url string, jsondata interface{}, header http.Header) ([]byte, error) {
	request := gorequest.New()
	request.SetDebug(debugClientHTTP)
	timeout, _ := time.ParseDuration(timeout)
	//_ := errors.New("Connection Problem")
	if url[:5] == "https" {
		request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	reqagent := request.Delete(url)
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
