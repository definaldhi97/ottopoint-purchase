package host

import (
	"encoding/json"
	"fmt"
	"net/http"
	https "ottopoint-purchase/hosts"
	modelsLP "ottopoint-purchase/hosts/landing_page/models"
	"ottopoint-purchase/utils"

	"github.com/sirupsen/logrus"
)

var (
	host            string
	endpointPayment string
	auth            string
	timestamp       string
	apikey          string
)

func init() {
	host = utils.GetEnv("HOST_SECUREPAGE", "http://18.136.193.154:8956")
	endpointPayment = utils.GetEnv("ENDPOINT_PAYMENT_SECUREPAGE", "/payment-services/v2.0.0/api/token")
	auth = utils.GetEnv("AUTH_SECUREPAGE", "Basic T1RUT1BBWQ==")
	timestamp = utils.GetEnv("TIMESTAMP_SECUREPAGE", "1613338906")
	apikey = utils.GetEnv("APIKEY_SECUREPAGE", "E60E0K0PAP00Y001PPK00A1IK0EA3P00")
}

func PaymentLandingPage(req modelsLP.LGRequestPay) (modelsLP.LGResponsePay, error) {
	var res modelsLP.LGResponsePay

	signature := createSignature(req, apikey, timestamp)

	urlSvr := host + endpointPayment

	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	header.Set("Signature", signature)
	header.Set("Timestamp", timestamp)
	header.Set("Authorization", auth)

	data, err := https.HTTPxPOSTwithRequest(urlSvr, req, header)
	if err != nil {

		logrus.Error("[PackageHost]-[PaymentLandingPage]")
		logrus.Error(fmt.Sprintf("[HTTP_POST_LP]-[Error : %v]", err))

		return res, err
	}

	err = json.Unmarshal(data, &res)
	if err != nil {

		logrus.Error("[PackageHost]-[PaymentLandingPage]")
		logrus.Error(fmt.Sprintf("[Unmarshal]-[Error : %v]", err))

		return res, err
	}

	return res, nil
}

func CheckStatusLandingPage(trxId string) (modelsLP.CheckStatusLPResp, error) {
	var res modelsLP.CheckStatusLPResp

	req := modelsLP.CheckStatusLPReq{
		TrxRef: trxId,
	}

	signature := createSignature(req, apikey, timestamp)

	urlSvr := host + endpointPayment

	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	header.Set("Signature", signature)
	header.Set("Timestamp", "1613338906")
	header.Set("Authorization", auth)

	data, err := https.HTTPxPOSTwithRequest(urlSvr, req, header)
	if err != nil {

		logrus.Error("[PackageHost]-[CheckStatusLandingPage]")
		logrus.Error(fmt.Sprintf("[HTTP_POST_LP]-[Error : %v]", err))

		return res, err
	}

	err = json.Unmarshal(data, &res)
	if err != nil {

		logrus.Error("[PackageHost]-[CheckStatusLandingPage]")
		logrus.Error(fmt.Sprintf("[Unmarshal]-[Error : %v]", err))

		return res, err
	}

	return res, nil
}

func createSignature(data interface{}, apikey string, timestamp string) string {
	var signature string

	jsonReq, _ := json.Marshal(data)
	bodyMsg := string(jsonReq)

	logrus.Error("[Create-Signature]")
	logrus.Error("Request : ", bodyMsg)
	// logrus.Error("Header : ", header)

	jsonReqString := utils.SignReplaceAll(bodyMsg)
	plainText := jsonReqString + "&" + timestamp + "&" + apikey
	logrus.Error("request data signature system : ", plainText)

	signatureSystem := utils.HashSha512(apikey, plainText)
	logrus.Error("signature : ", signatureSystem)

	signature = signatureSystem

	return signature
}
