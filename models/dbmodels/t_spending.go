package dbmodels

import (
	"time"
)

type TSpending struct {
	ID                string     `gorm:"id"`
	AccountNumber     string     `gorm:"account_number"`
	Voucher           string     `gorm:"voucher"`
	MerchantID        string     `gorm:"merchant_id"`
	CustID            string     `gorm:"cust_id"`
	RRN               string     `gorm:"rrn"`
	TransactionId     string     `grom:"transaction_id"`
	ProductCode       string     `gorm:"product_code"`
	Amount            int64      `gorm:"amount"`
	TransType         string     `gorm:"trans_type"`
	IsUsed            bool       `gorm:"is_used"`
	ProductType       string     `gorm:"product_type"`
	Status            string     `gorm:"status"`
	ExpDate           *time.Time `gorm:"exp_date"`
	Institution       string     `gorm:"institution"`
	CummulativeRef    string     `gorm:"cummulative_ref"`
	DateTime          string     `gorm:"date_time"`
	ResponderData     string     `gorm:"responder_data"`
	Point             int        `gorm:"point"`
	ResponderRc       string     `gorm:"responder_rc"`
	ResponderRd       string     `gorm:"responder_rd"`
	RequestorData     string     `gorm:"requestor_data"`
	RequestorOPData   string     `gorm:"requestor_op_data"`
	SupplierID        string     `gorm:"supplier_id"`
	CouponId          string     `gorm:"coupon_id"`
	CampaignId        string     `gorm:"campaign_id"`
	AccountId         string     `gorm:"account_id"`
	RedeemAt          *time.Time `gorm:"redeem_at"`
	UsedAt            *time.Time `gorm:"used_at"`
	CreatedAT         time.Time  `gorm:"created_at" json:"created_at"`
	UpdatedAT         time.Time  `gorm:"updated_at" json:"updated_at"`
	VoucherCode       string     `gorm:"voucher_code"`
	ProductCategoryID *string    `gorm:"product_category_id"`
	Comment           string     `gorm:"comment"`
	MRewardID         *string    `gorm:"m_reward_id"`
	MProductID        *string    `gorm:"m_product_id"`
	VoucherLink       string     `gorm:"voucher_link"`
	PointsTransferID  string     `gorm:"points_transfer_id"`
	InvoiceNumber     string     `gorm:"invoice_number"`
	PaymentMethod     int        `gorm:"payment_method"`
	IsCallback        bool       `gorm:"is_callback"`
}

func (t *TSpending) TableName() string {
	return "public.t_spending"
}
