package host

import (
	"encoding/json"
	"fmt"
	"net/http"
	https "ottopoint-purchase/hosts"
	modelsLP "ottopoint-purchase/hosts/landing_page/models"
)

type EnvLandingPage struct {
	Host            string `envconfig:"HOSTADDRES_AUTH_OTTOPOINT" default:"http://18.136.193.154:8956"`
	EndpointPayment string `envconfig:"ENDPOINT_SIGNATURE_OTTOPOINT" default:"/payment-services/v2.0.0/api/token"`
}

var (
	envLandingPage EnvLandingPage
)

func PaymentLandingPage(email, firstname, lastName, phone, merchantname, trxId string, amount int) (modelsLP.LGResponsePay, error) {
	var res modelsLP.LGResponsePay

	printError := "\033[31m" // merah
	// printSuccess := "\033[32m" // hijau
	// printRes := "\033[34m"     // biru

	req := modelsLP.LGRequestPay{
		Customerdetails: modelsLP.DataCustomerdetails{
			Email:     email,
			Firstname: firstname,
			Lastname:  lastName,
			Phone:     phone,
		},
		Transactiondetails: modelsLP.DataTransactiondetails{
			Amount:       amount,
			Currency:     "idr",
			Merchantname: merchantname,
			Orderid:      trxId,
			// PaymentMethod :
			// Promocode   :
			// Vabca       :
			// Valain      :
			// Vamandiri   :
		},
	}

	urlSvr := envLandingPage.Host + envLandingPage.EndpointPayment

	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	header.Set("Signature", "")
	header.Set("Timestamp", "")
	header.Set("Authorization", "")

	data, err := https.HTTPxPOSTwithRequest(urlSvr, req, header)
	if err != nil {

		fmt.Println(printError, "[PackageHost]-[PaymentLandingPage]")
		fmt.Println(printError, fmt.Sprintf("[HTTP_POST_LP]-[Error : %v]", err))

		return res, err
	}

	err = json.Unmarshal(data, &res)
	if err != nil {

		fmt.Println(printError, "[PackageHost]-[PaymentLandingPage]")
		fmt.Println(printError, fmt.Sprintf("[Unmarshal]-[Error : %v]", err))

		return res, err
	}

	return res, nil
}
