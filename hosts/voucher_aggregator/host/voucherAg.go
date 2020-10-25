package host

import (
	"encoding/json"
	"fmt"
	"ottopoint-purchase/hosts/voucher_aggregator/models"

	"github.com/google/go-querystring/query"

	"github.com/astaxie/beego/logs"
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

func OrderVoucher(req models.RequestOrderVoucherAg, head models.HeaderHTTP) (*models.ResponseOrderVoucherAg, error) {
	var resp models.ResponseOrderVoucherAg

	logs.Info("[PackageHostUV]-[OrderVoucher]")

	urlSvr := host + endpointOrderVoucher

	head.GenerateSignature(req)

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

func CheckStatusOrder(req models.RequestCheckOrderStatus, head models.HeaderHTTP) (*models.ResponseCheckOrderStatus, error) {
	var resp models.ResponseCheckOrderStatus

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
