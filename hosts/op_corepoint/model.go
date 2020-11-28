package op_corepoint

import "time"

type AddingPointReq struct {
	AccountID     string `json:"accountId"`
	TransactionID string `json:"transactionId"`
	ExpiredDays   string `json:"expiredDays"`
	Point         int    `json:"point"`
	Comment       string `json:"comment"`
}

type TrxPointRes struct {
	ResponseCode string      `json:"responseCode"`
	ResponseDesc string      `json:"responseDesc"`
	Data         interface{} `json:"data"`
}

type DataAddingPoint struct {
	PointsTransferID string    `json:"pointsTransferId"`
	ExpiredPoint     time.Time `json:"expiredPoint"`
}

type SpendingPointReq struct {
	AccountID     string `json:"accountId"`
	TransactionID string `json:"transactionId"`
	Point         int    `json:"point"`
	Comment       string `json:"comment"`
}

type DataSpendingPoint struct {
	PointsTransferID string `json:"pointsTransferId"`
}
