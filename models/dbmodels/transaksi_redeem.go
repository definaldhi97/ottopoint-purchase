package dbmodels

import "time"

type TransaksiRedeem struct {
	ID              int       `gorm:"id";pk json:"id"`
	AccountNumber   string    `gorm:"account_number"`
	Voucher         string    `gorm:"voucher"`
	MerchantID      string    `gorm:"merchant_id"`
	CustID          string    `gorm:"cust_id"`
	RRN             string    `gorm:"rrn"`
	ProductCode     string    `gorm:"product_code"`
	Amount          int64     `gorm:"amount"`
	TransType       string    `gorm:"trans_type"`
	ProductType     string    `gorm:"product_type"`
	Status          string    `gorm:"status"`
	ExpDate         string    `gorm:"exp_date"`
	Institution     string    `gorm:"institution"`
	CummulativeRef  string    `gorm:"cummulative_ref"`
	DateTime        string    `gorm:"date_time"`
	ResponderData   string    `gorm:"responder_data"`
	Point           int       `gorm:"point"`
	ResponderRc     string    `gorm:"responder_rc"`
	RequestorData   string    `gorm:"requestor_data"`
	ResponderData2  string    `gorm:"responder_data"`
	RequestorOPData string    `gorm:"requestor_op_data"`
	SupplierID      string    `gorm:"supplier_id"`
	CreatedAT       time.Time `gorm:"created_at" json:"created_at"`
	UpdatedAT       time.Time `gorm:"updated_at" json:"updated_at"`
}

func (t *TransaksiRedeem) TableName() string {
	return "public.redeem_transactions"
}
