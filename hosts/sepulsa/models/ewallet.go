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
	TransactionID  string      `json:"transaction_id"`
	Type           string      `json:"type"`
	Created        string      `json:"created"`
	Changed        string      `json:"changed"`
	CustomerNumber string      `json:"customer_number"`
	OrderID        string      `json:"order_id"`
	Price          string      `json:"price"`
	Status         string      `json:"status"`
	ResponseCode   string      `json:"response_code"`
	SerialNumber   string      `json:"serial_number"`
	Amount         string      `json:"amount"`
	ProductID      string      `json:"product_id"`
	Token          string      `json:"token"`
	Data           interface{} `json:"data"`
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

type CheckStatusSepulsaResp struct {
	Amount         string      `json:"amount"`
	Changed        string      `json:"changed"`
	Created        string      `json:"created"`
	CustomerNumber string      `json:"customer_number"`
	Data           interface{} `json:"data"`
	OrderID        string      `json:"order_id"`
	Price          string      `json:"price"`
	ProductID      struct {
		Enabled   string `json:"enabled"`
		Label     string `json:"label"`
		Nominal   string `json:"nominal"`
		Operator  string `json:"operator"`
		Price     int    `json:"price"`
		ProductID string `json:"product_id"`
		Type      string `json:"type"`
	} `json:"product_id"`
	ResponseCode  string      `json:"response_code"`
	SerialNumber  interface{} `json:"serial_number"`
	Status        string      `json:"status"`
	Token         interface{} `json:"token"`
	TransactionID string      `json:"transaction_id"`
	Type          string      `json:"type"`
}
