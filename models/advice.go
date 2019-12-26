package models

import (
	ottomart "ottopoint-purchase/hosts/ottomart/models"
	"ottopoint-purchase/models/ottoagmodels"
	"time"
)

type TrxHistory struct {
	ID                        int       `gorm:"column:id"`
	Amount                    int       `gorm:"column:amount"`
	Name                      string    `gorm:"column:name"`
	PaymentMethod             string    `gorm:"column:payment_method"`
	ReferenceNumber           string    `gorm:"column:reference_number"`
	Description               string    `gorm:"column:description"`
	Detail                    string    `gorm:"column:detail"`
	TransactionAt             time.Time `gorm:"column:transaction_at"`
	UserID                    int       `gorm:"column:user_id"`
	CreatedAt                 time.Time `gorm:"column:created_at"`
	UpdatedAt                 time.Time `gorm:"column:updated_at"`
	Status                    int       `gorm:"column:status"`
	NoResi                    string    `gorm:"column:no_resi"`
	MiddlewareReferenceNumber string    `gorm:"column:middleware_reference_number"`
	BillerReference           string    `gorm:"column:biller_reference"`
	CustomerName              string    `gorm:"column:customer_name"`
	Commission                string    `gorm:"column:commission"`
	StroomToken               string    `gorm:"column:stroom_token"`
	Voucher                   string    `gorm:"column:voucher"`
}

// PpobPaymentAdvice ...
type PpobPaymentAdvice struct {
	Req         ottoagmodels.OttoAGPaymentReq
	Db          TrxHistory
	RrnInquiry  string
	RrnReversal string
	DataToken   ottomart.AccessToken
	Layanan     string
}

// OttoAGPaymentReq ...
type OttoAGPaymentReq struct {
	Amount      uint64 `json:"amount"`
	CustID      string `json:"custid"`
	MemberID    string `json:"memberid"`
	Period      string `json:"period"`
	Productcode string `json:"productcode"`
	Rrn         string `json:"rrn"`
}
