package redismodels

import (
	"ottopoint-purchase/models"
)

type PayQrRedis struct {
	QrData        string
	InqRes        models.InquiryRes
	Resp          models.PayQrInquiryDataResponse
	AccountNumber string
	MerchantID    string
	YoutapStatus  bool
	Amount        int64
}
