package models

type EwalletInsertTrxReq struct {
	CustomerNumber string `json:"customer_number"`
	OrderID        string `json:"order_id"`
	ProductID      int    `json:"product_id"`
}

type EwalletInsertTrxRes struct {
	Transaction
}

type CallbackTrxReq struct {
	TransactionID  string `json:"transaction_id"`
	Type           string `json:"type"`
	Created        string `json:"created"`
	Changed        string `json:"changed"`
	CustomerNumber string `json:"customer_number"`
	OrderID        string `json:"order_id"`
	Price          string `json:"price"`
	Status         string `json:"status"`
	ResponseCode   string `json:"response_code"`
	SerialNumber   string `json:"serial_number"`
	Amount         string `json:"amount"`
	ProductID      string `json:"product_id"`
	Token          string `json:"token"`
	Data           string `json:"data"`
}

type Transaction struct {
	TransactionID  string  `json:"transaction_id"`
	Type           string  `json:"type"`
	Created        string  `json:"created"`
	Changed        string  `json:"changed"`
	CustomerNumber string  `json:"customer_number"`
	OrderID        string  `json:"order_id"`
	Price          string  `json:"price"`
	Status         string  `json:"status"`
	ResponseCode   string  `json:"response_code"`
	SerialNumber   string  `json:"serial_number"`
	Amount         string  `json:"amount"`
	ProductID      Product `json:"product_id"`
	Token          string  `json:"token"`
	Data           string  `json:"data"`
}

type Product struct {
	ProductID string `json:"product_id"`
	Type      string `json:"type"`
	Label     string `json:"label"`
	Operator  string `json:"operator"`
	Nominal   string `json:"nominal"`
	Price     int    `json:"price"`
	Enabled   string `json:"enabled"`
}
