package host

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"ottopoint-purchase/utils"
	"strconv"
	"strings"
	"time"

	redishost "ottopoint-purchase/hosts/redis_token/host"

	"github.com/astaxie/beego/logs"
	"github.com/parnurzeal/gorequest"
	"github.com/sirupsen/logrus"
)

var (
	debugClientHTTP bool
	timeout         string
	retrybad        int
)

func init() {
	debugClientHTTP = true //defaultValue
	if dch := utils.GetEnv("HTTP_DEBUG_CLIENT", "true"); strings.EqualFold(dch, "true") || strings.EqualFold(dch, "false") {
		debugClientHTTP, _ = strconv.ParseBool(strings.ToLower(dch))
	}
	timeout = utils.GetEnv("HTTP_TIMEOUT_OPL", "60s")
	retrybad = 1
	if rb := utils.GetEnv("HTTP_RETRY_BAD", "1"); strings.TrimSpace(rb) != "" {
		if val, err := strconv.Atoi(rb); err == nil {
			retrybad = val
		}
	}

}

// Post (Tanpa Request), Token Customer
func HTTPxFormPostCustomerWithoutRequest(url, phone string) ([]byte, error) {
	logrus.Info("PhoneNumber :", phone)
	token, _ := redishost.GetToken(fmt.Sprintf("Ottopoint-Token-Customer-%s :", phone))
	data := strings.Replace(token.Data, `"`, "", 2)
	dataToken := "Bearer" + " " + data
	logrus.Info("Token :", dataToken)
	request := gorequest.New()
	request.SetDebug(debugClientHTTP)
	timeout, _ := time.ParseDuration(timeout)
	//_ := errors.New("Connection Problem")
	if strings.HasPrefix(url, "https") {
		request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	reqagent := request.Post(url)
	// reqagent.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqagent.Header.Set("Authorization", dataToken)
	_, body, errs := reqagent.
		// Send(jsondata).
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

// Post (Request), Token Customer
func HTTPxFormPostCustomerWithRequest(url, phone string, jsondata interface{}) ([]byte, error) {
	logrus.Info("PhoneNumber :", phone)
	token, _ := redishost.GetToken(fmt.Sprintf("Ottopoint-Token-Customer-%s :", phone))
	data := strings.Replace(token.Data, `"`, "", 2)
	dataToken := "Bearer" + " " + data
	logrus.Info("Token :", dataToken)
	request := gorequest.New()
	request.SetDebug(debugClientHTTP)
	timeout, _ := time.ParseDuration(timeout)
	//_ := errors.New("Connection Problem")
	if strings.HasPrefix(url, "https") {
		request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	reqagent := request.Post(url)
	reqagent.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqagent.Header.Set("Authorization", dataToken)
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

// GET, Token Admin
func HTTPxFormGETAdmin(url string) ([]byte, error) {
	token, _ := redishost.GetToken(utils.RedisKeyAuth)
	data := strings.Replace(token.Data, `"`, "", 2)
	dataToken := "Bearer" + " " + data
	request := gorequest.New()
	request.SetDebug(debugClientHTTP)
	timeout, _ := time.ParseDuration(timeout)
	//_ := errors.New("Connection Problem")
	if strings.HasPrefix(url, "https") {
		request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	reqagent := request.Get(url)
	// reqagent.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqagent.Header.Set("Authorization", dataToken)
	_, body, errs := reqagent.
		// Send(jsondata).
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

// GET, Token Customer
func HTTPxFormGETCustomer(url, phone string) ([]byte, error) {

	token, _ := redishost.GetToken(fmt.Sprintf("Ottopoint-Token-Customer-%s :", phone))
	data := strings.Replace(token.Data, `"`, "", 2)
	dataToken := "Bearer" + " " + data
	request := gorequest.New()
	request.SetDebug(debugClientHTTP)
	timeout, _ := time.ParseDuration(timeout)
	//_ := errors.New("Connection Problem")
	if strings.HasPrefix(url, "https") {
		request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	reqagent := request.Get(url)
	// reqagent.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqagent.Header.Set("Authorization", dataToken)

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

// Post (Request), Token Admin
func HTTPxFormPostAdminWithRequest(url string, jsondata interface{}) ([]byte, error) {
	token, _ := redishost.GetToken(utils.RedisKeyAuth)
	data := strings.Replace(token.Data, `"`, "", 2)
	dataToken := "Bearer" + " " + data
	// logrus.Info("Token :", dataToken)
	request := gorequest.New()
	request.SetDebug(debugClientHTTP)
	timeout, _ := time.ParseDuration(timeout)
	//_ := errors.New("Connection Problem")
	if strings.HasPrefix(url, "https") {
		request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	reqagent := request.Post(url)
	reqagent.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqagent.Header.Set("Authorization", dataToken)
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

// Post (Request), Token Admin
func HTTPxFormPostAdminWithoutRequest(url string, jsondata interface{}) ([]byte, error) {
	token, _ := redishost.GetToken(utils.RedisKeyAuth)
	data := strings.Replace(token.Data, `"`, "", 2)
	dataToken := "Bearer" + " " + data
	// logrus.Info("Token :", dataToken)
	request := gorequest.New()
	request.SetDebug(debugClientHTTP)
	timeout, _ := time.ParseDuration(timeout)
	//_ := errors.New("Connection Problem")
	if strings.HasPrefix(url, "https") {
		request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	reqagent := request.Get(url)
	// reqagent.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqagent.Header.Set("Authorization", dataToken)
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

// Customer Status
func HTTPxFormCustomerStatus(url string, customer string) ([]byte, error) {
	token, _ := redishost.GetToken(utils.RedisKeyAuth)
	data := strings.Replace(token.Data, `"`, "", 2)
	dataToken := "Bearer" + " " + data
	log.Print("ini token status : ", dataToken)
	request := gorequest.New()
	request.SetDebug(debugClientHTTP)
	timeout, _ := time.ParseDuration(timeout)
	//_ := errors.New("Connection Problem")
	if strings.HasPrefix(url, "https") {
		request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	reqagent := request.Get(url + customer)
	log.Print("url status : ", reqagent)
	// reqagent := request.Get(url)
	// reqagent.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reqagent.Header.Set("Authorization", dataToken)
	// reqagent.Param("customer", customer)
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
