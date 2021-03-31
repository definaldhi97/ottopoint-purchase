package payment

type SpendingPaymentReq struct {
	AccountNumber   string `json:"accountNumber"`
	TransType       string `json:"transType"`
	ProductName     string `json:"productName"`
	ReferenceId     string `json:"referenceId"`
	TransactionTime string `json:"transactionTime"`
	Point           int    `json:"point"`
	Cash            int    `json:"cash"`
	Amount          int    `json:"amount"`
	Comment         string `json:"comment"`
	PaymentMethod   int    `json:"paymentMethod"`
}
