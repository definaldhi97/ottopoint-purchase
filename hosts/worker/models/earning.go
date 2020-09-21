package models

// === Api Earning ===
type WorkerEarningReq struct {
	InstitutionId   string `json:"institutionId"`
	AccountNumber1  string `json:"accountNumber1"`
	AccountNumber2  string `json:"accountNumber2"`
	Earning         string `json:"earning"`
	ReferenceId     string `json:"referenceId"`
	ProductCode     string `json:"productCode"`
	ProductName     string `json:"productName"`
	Amount          int64  `json:"amount"`
	Remark          string `json:"remark"`
	TransactionTime string `json:"transactionTime"`
}

type WorkerEarningResp struct {
	Data    DataWorkerEarning `json:"data"`
	Code    int               `json:"code"`
	Message string            `json:"message"`
}

type DataWorkerEarning struct {
	Point       int64  `json:"point"`
	ReferenceId string `json:"referenceId"`
}
