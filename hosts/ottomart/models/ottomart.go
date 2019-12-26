package models

type RedeemVoucherRequest struct {
	AccountNumber string `json:"accountNumber"`
	CustID        string `json:"cust_id"`
	ProductCode   string `json:"product_code"`
}

// Response
type RedeemVoucherResp struct {
	Data RedeemVoucherResp1 `json:"data"`
	Meta Metas              `json:"meta"`
}

type Metas struct {
	Status  bool   `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type RedeemVoucherResp1 struct {
	Rc            string `json:"rc"`
	AccountNumber string `json:"accountNumber"`
	CustID        string `json:"cust_id"`
	CampaignID    string `json:"campaign_id"`
	Name          string `json:"name"`
	ProductCode   string `json:"product_code"`
	Amount        int64  `json:"amount"`
}

// Voucher
type CheckingResponse struct {
	CampaignID string `json:"CampaignID"`
}

// Token
type ResponseToken struct {
	Data struct {
		AccountNumber string `json:"accountNumber"`
		Token         string `json:"token"`
		DeviceID      string `json:"device-id"`
		FirebaseToken string `json:"firebase_token"`
		MerchantID    string `json:"merchantId"`
		Name          string `json:"name" gorm:"column:name"`
		MerchantName  string `json:"merchantName"`
		UserID        int    `json:"userId"`
	} `json:"data"`
	Meta struct {
		Status  bool   `json:"status"`
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"meta"`
}

type AccessToken struct {
	UserID        int    `json:"userId"`
	FirebaseToken string `json:"firebase_token"`
}
