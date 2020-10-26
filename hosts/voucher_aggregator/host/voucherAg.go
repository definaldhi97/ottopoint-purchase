package host

import (
	"encoding/json"
	"fmt"
	vgmodel "ottopoint-purchase/hosts/voucher_aggregator/models"
	"ottopoint-purchase/models"

	"github.com/astaxie/beego/logs"
	"github.com/google/go-querystring/query"
	ODU "ottodigital.id/library/utils"
)

var (
	host string
	name string

	endpointOrderVoucher     string
	endpointCheckStatusOrder string
)

func init() {
	host = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_VOUCHERAG", "http://13.228.25.85:8480/api")
	name = ODU.GetEnv("OTTOPOINT_PURCHASE_NAME_VOUCHERAG", "VOUCHER AGGREGATOR")

	endpointOrderVoucher = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_VOUCHERAR_ORDER", "/v1/order")
	endpointCheckStatusOrder = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_VOUCHERAR_CHECK_STATUS", "/v1/order/status")

}

func OrderVoucher(req vgmodel.RequestOrderVoucherAg, head models.RequestHeader) (*vgmodel.ResponseOrderVoucherAg, error) {
	var resp vgmodel.ResponseOrderVoucherAg

	logs.Info("[PackageHostUV]-[OrderVoucher]")

	urlSvr := host + endpointOrderVoucher

	data, err := HTTPxFormPostVoucherAg(urlSvr, head, req)
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

	data, err := HTTPxFormGetVoucherAg(urlSvr, head)
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
