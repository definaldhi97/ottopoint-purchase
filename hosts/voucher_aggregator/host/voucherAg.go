package host

import (
	"encoding/json"
	"fmt"
	"net/http"
	https "ottopoint-purchase/hosts"
	vgmodel "ottopoint-purchase/hosts/voucher_aggregator/models"
	"ottopoint-purchase/models"

	"ottopoint-purchase/utils"

	"github.com/astaxie/beego/logs"
	"github.com/google/go-querystring/query"
)

var (
	host string
	name string

	endpointOrderVoucher     string
	endpointCheckStatusOrder string
)

func init() {
	host = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_VOUCHERAG", "http://13.228.25.85:8480/api")
	name = utils.GetEnv("OTTOPOINT_PURCHASE_NAME_VOUCHERAG", "VOUCHER AGGREGATOR")

	endpointOrderVoucher = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_VOUCHERAR_ORDER", "/v1/order")
	endpointCheckStatusOrder = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_VOUCHERAR_CHECK_STATUS", "/v1/order/status")

}

func OrderVoucher(req vgmodel.RequestOrderVoucherAg, head models.RequestHeader) (*vgmodel.ResponseOrderVoucherAg, error) {
	var resp vgmodel.ResponseOrderVoucherAg

	logs.Info("[PackageHostUV]-[OrderVoucher]")

	urlSvr := host + endpointOrderVoucher

	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	header.Set("InstitutionId", head.InstitutionID)
	header.Set("DeviceId", head.DeviceID)
	header.Set("Geolocation", head.Geolocation)
	header.Set("ChannelId", head.ChannelID)
	header.Set("Signature", head.Signature)
	header.Set("AppsId", head.AppsID)
	header.Set("Timestamp", head.Timestamp)

	data, err := https.HTTPxPOSTwithRequest(urlSvr, req, header)
	// data, err := HTTPxFormPostVoucherAg(urlSvr, head, req)
	if err != nil {
		logs.Error("Check error : ", err.Error())

		return &resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logs.Error("Failed to unmarshaling response OrderVoucher from Ultra Voucher ", err.Error())

		return &resp, err
	}

	return &resp, nil
}

func CheckStatusOrder(req vgmodel.RequestCheckOrderStatus, head models.RequestHeader) (*vgmodel.ResponseCheckOrderStatus, error) {
	var resp vgmodel.ResponseCheckOrderStatus

	logs.Info("[PackageHostVoucherAg]-[CheckStatusOrder]")

	v, _ := query.Values(req)

	urlSvr := host + endpointCheckStatusOrder + fmt.Sprintf("?%s", v.Encode())

	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	header.Set("InstitutionId", head.InstitutionID)
	header.Set("DeviceId", head.DeviceID)
	header.Set("Geolocation", head.Geolocation)
	header.Set("ChannelId", head.ChannelID)
	header.Set("Signature", head.Signature)
	header.Set("AppsId", head.AppsID)
	header.Set("Timestamp", head.Timestamp)

	data, err := https.HTTPxGET(urlSvr, header)
	// data, err := HTTPxFormGetVoucherAg(urlSvr, head)
	if err != nil {
		logs.Error("Check error : ", err.Error())

		return &resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logs.Error("Failed to unmarshaling response OrderVoucher from Ultra Voucher ", err.Error())

		return &resp, err
	}

	return &resp, nil

}
