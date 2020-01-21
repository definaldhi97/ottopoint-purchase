package models

type UseVoucherReq struct {
	Category    string `json:"category"`
	CampaignID  string `json:"campaignId"`
	CustID      string `json:"cust_id"`
	CustID2     string `json:"cust_id2"`
	ProductCode string `json:"product_code"`
	Date        string `json:"date"`
}

type UseVoucherResp struct {
	Points        int    `json:"points"`
	Code          string `json:"code"`
	CampaignID    string `json:"campaignId"`
	CouponID      string `json:"couponId"`
	AccountNumber string `json:"account_number"`
	CustID        string `json:"cust_id"`
	Date          string `json:"date"`
}

type UseRedeemRequest struct {
	AccountNumber string `json:"account_number"`
	CustID        string `json:"cust_id"`
	CustID2       string `json:"cust_id2"`
	ProductCode   string `json:"product_code"`
}

type UseRedeemResponse struct {
	Rc          string      `json:"rc"`
	Rrn         string      `json:"rrn"`
	CustID      string      `json:"cust_id"`
	ProductCode string      `json:"product_code"`
	Amount      int64       `json:"amount"`
	Msg         string      `json:"msg"`
	Uimsg       string      `json:"uimsg"`
	Datetime    string      `json:"datetime"`
	Data        interface{} `json:"data"`
}

type CampaignsDetail struct {
	Name         string       `json:"name"`
	BrandName    string       `json:"brand_name"`
	CanBeUsed    bool         `json:"canBeUsed"`
	PurchaseAt   string       `json:"purchaseAt"`
	CostInPoints int          `json:"costInPoints"`
	CampaignID   string       `json:"campaignId"`
	Used         bool         `json:"used"`
	Coupon       CouponDetail `json:"coupon"`
	Status       string       `json:"status"`
	ActiveTo     string       `json:"activeTo"`
	// DeliveryStatus string       `json:"deliveryStatus"`
	UrlPhoto    []Photo `json:"url_photo"`
	UrlLogo     string  `json:"url_logo"`
	VoucherTime string  `json:"voucher_time"`
}

type CouponDetail struct {
	ID   string `json:"id"`
	Code string `json:"code"`
}

type Photo struct {
	URLPhoto string `json:"url_photo"`
}
