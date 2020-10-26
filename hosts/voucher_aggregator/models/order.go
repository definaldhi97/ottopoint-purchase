package models

import (
	"ottopoint-purchase/utils"
	"strconv"
	"time"

	ODU "ottodigital.id/library/utils"
)

var (
	serverkey = ODU.GetEnv("OTTOPOINT_PURCHASE_VOUCHERAG_SESSIONKEY", "")
)

type RequestOrderVoucherAg struct {
	ProductCode    string `json:"productCode"`
	Qty            int    `json:"qty"`
	OrderID        string `json:"orderId"`
	CustomerName   string `json:"customerName"`
	CustomerEmail  string `json:"customerEmail"`
	CustomerPhone  string `json:"customerPhone"`
	DeliveryMethod int    `json:"deliveryMethod"`
	RedeemCallback string `json:"redeemCallback"`
}

type ResponseOrderVoucherAg struct {
	ResponseCode string `json:"responseCode"`
	ResponseDesc string `json:"responseDesc"`
}

type ResponseVoucherAg struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	Success int    `json:"success"`
	Failed  int    `json:"failed"`
	Pending int    `json:"pending"`
}

type HeaderHTTP struct {
	Institution string
	DeviceID    string
	Geolocation string
	AppsID      string
	Signature   string
	Timestamp   string
}

func (h *HeaderHTTP) GenerateSignature(req interface{}) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature := utils.VoucherAggregatorSignature(timestamp, req, serverkey)

	h.Signature = signature
	h.Timestamp = timestamp
}

type RequestCheckOrderStatus struct {
	OrderID       string `url:"orderId"`
	CurrentPage   string `url:"currentPage"`
	RecordPerPage string `url:"recordPerPage"`
}

type ResponseCheckOrderStatus struct {
	ResponseCode string `json:"responseCode"`
	ResponseDesc string `json:"responseDesc"`
	Data         Data   `json:"data"`
}

type Data struct {
	ProductCode    string    `json:"productCode"`
	Qty            int       `json:"qty"`
	OrderID        string    `json:"orderId"`
	OrderDate      string    `json:"orderDate"`
	TransactionID  string    `json:"transactionId"`
	DeliveryMethod string    `json:"deliveryMethod"`
	Vouchers       []Voucher `json:"voucher"`
	TotalRecord    int       `json:"totalRecord"`
	RecordPerPage  int       `json:"recordPerPage"`
	CurrentPage    int       `json:"currentPage"`
	TotalPage      int       `json:"totalPage"`
}

type Voucher struct {
	VoucherID   string `json:"voucherId"`
	VoucherCode string `json:"voucherCode"`
	VoucherName string `json:"voucherName"`
	ExpiredDate string `json:"expiredDate,omitempty"`
	Link        string `json:"link"`
}
