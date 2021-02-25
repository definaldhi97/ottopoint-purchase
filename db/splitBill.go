package db

import (
	"fmt"
	"ottopoint-purchase/models/dbmodels"
	"time"

	"github.com/sirupsen/logrus"
)

type GroupSPlitBill struct {
	ID                string     `gorm:"id"`
	TSpendingID       string     `gorm:"t_spending_id"`
	ExternalReffId    string     `gorm:"external_reff_id"`
	TransType         string     `gorm:"trans_type"`
	Value             int64      `grom:"value"`
	ValueType         string     `gorm:"value_type"`
	Status            string     `gorm:"status"`
	ResponseRc        string     `gorm:"response_rc"`
	ResponseRd        string     `gorm:"response_rd"`
	AccountNumber     string     `gorm:"account_number"`
	Voucher           string     `gorm:"voucher"`
	MerchantID        string     `gorm:"merchant_id"`
	CustID            string     `gorm:"cust_id"`
	RRN               string     `gorm:"rrn"`
	TransactionId     string     `grom:"transaction_id"`
	ProductCode       string     `gorm:"product_code"`
	Amount            int64      `gorm:"amount"`
	IsUsed            bool       `gorm:"is_used"`
	ProductType       string     `gorm:"product_type"`
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
	ProductCategoryId string     `gorm:"product_category_id"`
	Comment           string     `gorm:"comment"`
	MRewardID         string     `gorm:"m_reward_id"`
	MProductID        string     `gorm:"m_product_id"`
	VoucherLink       string     `gorm:"voucher_link"`
	PointsTransferID  string     `gorm:"points_transfer_id"`
}

func GetDataSplitBillbyTrxID(trxId string) ([]GroupSPlitBill, error) {
	res := []GroupSPlitBill{}

	err := DbCon.Raw(`select * from t_payment as a join t_spending as b on a.t_spending_id = b.id where a.external_reff_id = ? and trans_type = 'PCS01'`, trxId).Scan(&res).Error
	if err != nil {

		logrus.Error("[PackageDB]-[GetDataSplitBillbyTrxID]")
		logrus.Error(fmt.Sprintf("[TrxID : %v]-[Error : %v]", trxId, err))

		return res, err
	}

	return res, nil
}

func UpdateTransactionSplitBill(used bool, trxID, status, rc, rd string, respVendor, reqVendor, reqOP interface{}) (dbmodels.TSpending, error) {

	res := dbmodels.TSpending{}
	tx := DbCon.Begin()

	queryString := fmt.Sprintf(`update t_spending set is_used = %v , status = %v, responder_data = %v, requestor_data = %v, responder_rc = %v, responder_rd = %v, requestor_op_data = %v where transaction_id = %v`, used, status, respVendor, reqVendor, rc, rd, reqOP)
	if err := tx.Exec(queryString).Error; err != nil {

		logrus.Error("[PackageDB]-[UpdateTransactionSplitBill]")
		logrus.Error(fmt.Sprintf("[TrxID : %v]-[Error : %v]", trxID, err))

		return res, err
	}

	return res, nil

}
