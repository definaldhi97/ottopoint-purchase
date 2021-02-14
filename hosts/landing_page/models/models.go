package landing_page

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
