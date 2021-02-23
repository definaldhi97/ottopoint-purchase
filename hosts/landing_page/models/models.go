package landing_page

import "time"

type LGRequestPay struct {
	Customerdetails    DataCustomerdetails    `json:"customerdetails"`
	Transactiondetails DataTransactiondetails `json:"transactiondetails"`
}

type DataCustomerdetails struct {
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Phone     string `json:"phone"`
}

type DataTransactiondetails struct {
	Amount        int    `json:"amount"`
	Currency      string `json:"currency"`
	Merchantname  string `json:"merchantname"`
	Orderid       string `json:"orderid"`
	PaymentMethod int    `json:"paymentMethod"`
	Promocode     string `json:"promocode"`
	Vabca         string `json:"vabca"`
	Valain        string `json:"valain"`
	Vamandiri     string `json:"vamandiri"`
}

type LGResponsePay struct {
	// Failed
	ResponseCode string `json:"responseCode"`
	ResponseDesc string `json:"responseDesc"`

	// Success
	ResponseAuth struct {
		Signature string `json:"signature"`
	} `json:"responseAuth"`
	ResponseData struct {
		StatusCode    string `json:"statusCode"`
		StatusMessage string `json:"statusMessage"`
		OrderID       string `json:"orderId"`
		EndpointURL   string `json:"endpointUrl"`
	} `json:"responseData"`
}

// Check Status
type CheckStatusLPReq struct {
	TrxRef string `json:"trxRef"`
}

type CheckStatusLPResp struct {
	ResponseCode          string    `json:"responseCode"`
	ResponseDesc          string    `json:"responseDesc"`
	TrxRef                string    `json:"trxRef"`
	Issuer                string    `json:"issuer"`
	IssuerRefNo           string    `json:"issuerRefNo"`
	OttoRefNo             string    `json:"ottoRefNo"`
	TransactionStatusCode string    `json:"transactionStatusCode"`
	TransactionStatusDesc string    `json:"transactionStatusDesc"`
	Amount                string    `json:"amount"`
	TransactionTime       time.Time `json:"transactionTime"`
	CustomerId            string    `json:"customerId"`
}

// Callback
type CallBackLPReq struct {
	Amount              string `json:"amount"`
	Issuer              string `json:"issuer"`
	IssuerRefNo         string `json:"issuerRefNo"`
	OttoRefNo           string `json:"ottoRefNo"`
	ResponseCode        string `json:"responseCode"`
	ResponseDescription string `json:"responseDescription"`
	TrxRef              string `json:"trxRef"`
	TransactionType     string `json:"transactionType"`
	UserId              string `json:"userId"`
}

type CallBackLPResp struct {
	ResponseCode string `json:"responseCode"`
	ResponseDesc string `json:"responseDesc"`
}
