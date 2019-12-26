package models

// Request
type RedeemVoucherRequest struct {
	VoucherName string `json:"voucher_name"`
	Category    string `json:"category"`
	CampaignID  string `json:"campaign_id"`
	CustID      string `json:"cust_id"`
	CustID2     string `json:"cust_id2"`
	ProductCode string `json:"product_code"`
}

// Response
type RedeemVoucherResp struct {
	Rc            string `json:"rc"`
	AccountNumber string `json:"accountNumber"`
	CustID        string `json:"cust_id"`
	CampaignID    string `json:"campaign_id"`
	Name          string `json:"name"`
	ProductCode   string `json:"product_code"`
	Amount        int64  `json:"amount"`
}
