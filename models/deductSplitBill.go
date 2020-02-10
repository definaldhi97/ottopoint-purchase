package models

type DeductPointReq struct {
	AccountNumber string `json:"accountNumber"`
	Point         int    `json:"point"`
	DeductType    int    `json:"deductType` // 1 (FullPoint), 2(SplitBill)
	TrxID         string `json:"trxID"`
	Amount        int    `json:"amount"`
	ProductCode   string `json:"productCode"`
	ProductName   string `json:"productName"`
}

// type DeductPointResp struct {

// }
