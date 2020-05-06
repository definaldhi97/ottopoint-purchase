package models

import (
	ottoagmodels "ottopoint-purchase/models/ottoag"
)

type UseVoucherReq struct {
	CampaignID string `json:"campaignId"`
	CustID     string `json:"custId"`
	CustID2    string `json:"custId2"`
}

type UseVoucherResp struct {
	Points        int    `json:"points"`
	Code          string `json:"code"`
	CampaignID    string `json:"campaignId"`
	CouponID      string `json:"couponId"`
	AccountNumber string `json:"accountNumber"`
	CustID        string `json:"custId"`
	Date          string `json:"date"`
}

type UseRedeemRequest struct {
	AccountNumber string `json:"account_number"`
	CustID        string `json:"custId"`
	CustID2       string `json:"custId2"`
	ProductCode   string `json:"product_code"`
	Jumlah        int
}

type UseRedeemResponse struct {
	Rc          string                    `json:"rc"`
	Rrn         string                    `json:"rrn"`
	Category    string                    `json:"category"`
	CustID      string                    `json:"custId"`
	CustID2     string                    `json:"custId2"`
	ProductCode string                    `json:"product_code"`
	Amount      int64                     `json:"amount"`
	Msg         string                    `json:"msg"`
	Uimsg       string                    `json:"uimsg"`
	Datetime    string                    `json:"datetime"`
	Data        ottoagmodels.DataGabungan `json:"data"`
}

type ResponseUseVoucher struct {
	Voucher     string `json:"voucher"`
	CustID      string `json:"custId"`
	CustID2     string `json:"custId2"`
	ProductCode string `json:"product_code"`
	Amount      int64  `json:"amount"`
}

type ResponseUseVoucherPLN struct {
	Voucher     string `json:"voucher"`
	CustID      string `json:"custId"`
	CustID2     string `json:"custId2"`
	ProductCode string `json:"product_code"`
	Amount      int64  `json:"amount"`
	Token       string `json:"token"`
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
