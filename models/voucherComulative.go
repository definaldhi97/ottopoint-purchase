package models

type VoucherComultaiveReq struct {
	Jumlah     int    `json:"total"`
	CampaignID string `json:"campaignId"`
	CustID     string `json:"custId"`
	CustID2    string `json:"custId2"`
}

type CommulativeResp struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	Success int    `json:"success"`
	Failed  int    `json:"failed"`
	Pending int    `json:"pending"`
	//RedeemRes RedeemResponse `json:"redeem_res"`
}

// type RedeemRequest struct {
// 	AccountNumber string `json:"accountNumber"`
// 	CustID        string `json:"custId"`
// 	CustID2       string `json:"custId2"`
// 	ProductCode   string `json:"productCode"`
// }

type RedeemResponse struct {
	Rc          string      `json:"rc"`
	Rrn         string      `json:"rrn"`
	CustID      string      `json:"custId"`
	ProductCode string      `json:"productCode"`
	Amount      int64       `json:"amount"`
	Msg         string      `json:"msg"`
	Uimsg       string      `json:"uimsg"`
	Datetime    string      `json:"datetime"`
	Data        interface{} `json:"data"`
}

type RedeemComuResp struct {
	Code            string         `json:"code"`
	Message         string         `json:"message"`
	CouponID        string         `json:"couponId"`
	CouponCode      string         `json:"couponCode"`
	PointTransferID string         `json:"pointTransferID"`
	Comment         string         `json:"comment"`
	Redeem          RedeemResponse `json:"redeem"`
}
