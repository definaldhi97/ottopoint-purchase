package host

import (
	"encoding/json"
	"fmt"
	"net/http"
	https "ottopoint-purchase/hosts"
	vgmodel "ottopoint-purchase/hosts/voucher_aggregator/models"
	"ottopoint-purchase/models"
	"strconv"
	"time"

	"ottopoint-purchase/utils"

	"github.com/google/go-querystring/query"
	"github.com/sirupsen/logrus"
)

var (
	host string
	name string

	endpointOrderVoucher        string
	endpointOrderVoucherV11     string
	endpointCheckStatusOrder    string
	endpointCheckStatusOrderV21 string
	endpointPaymentInfo         string
	endpointCallbackSepulsa     string
)

func init() {
	host = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_VOUCHERAG", "http://13.228.25.85:8480/transaction")
	name = utils.GetEnv("OTTOPOINT_PURCHASE_NAME_VOUCHERAG", "VOUCHER AGGREGATOR")

	endpointOrderVoucher = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_VOUCHERAR_ORDER", "/v1/order")
	endpointOrderVoucherV11 = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_VOUCHERAR_ORDER_V1.1", "/v1.1/order")
	endpointCheckStatusOrder = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_VOUCHERAR_CHECK_STATUS", "/v1/order/status")
	endpointPaymentInfo = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_VOUCHERAR_PAYMENT_INFO", "/v1/product/payment/info")
	endpointCheckStatusOrderV21 = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_VOUCHERAR_CHECK_STATUS_V2.1", "/v1.1/order/status")
	endpointCallbackSepulsa = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_VOCUHERAR_SEPULSA", "/v1/callback/sepulsa")

}

func OrderVoucher(req vgmodel.RequestOrderVoucherAg, head models.RequestHeader) (*vgmodel.ResponseOrderVoucherAg, error) {
	var resp vgmodel.ResponseOrderVoucherAg

	logrus.Info("[PackageHostVoucherAG]-[OrderVoucher]")

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
		logrus.Error("Check error : ", err.Error())

		return &resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logrus.Error("Failed to unmarshaling response OrderVoucher from VoucherAG ", err.Error())

		return &resp, err
	}

	return &resp, nil
}

func OrderVoucherV11(req vgmodel.RequestOrderVoucherAgV11, head models.RequestHeader) (*vgmodel.ResponseOrderVoucherAg, error) {
	var resp vgmodel.ResponseOrderVoucherAg

	logrus.Info("[PackageHostVoucherAg]-[OrderVoucherV1.1]")

	urlSvr := host + endpointOrderVoucherV11

	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	header.Set("InstitutionId", head.InstitutionID)
	header.Set("DeviceId", head.DeviceID)
	header.Set("Geolocation", head.Geolocation)
	header.Set("ChannelId", head.ChannelID)
	header.Set("Signature", head.Signature)
	header.Set("AppsId", head.AppsID)
	header.Set("Timestamp", head.Timestamp)
	header.Set("Authorization", head.Authorization)

	data, err := https.HTTPxPOSTwithRequest(urlSvr, req, header)
	// data, err := HTTPxFormPostVoucherAg(urlSvr, head, req)
	if err != nil {

		logrus.Error("[PackageHostVoucherAg]-[OrderVoucherV1.1]")
		logrus.Error(fmt.Sprintf("[HTTPxPOSTwithRequest]-[Error : %v]", err.Error()))

		return &resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {

		logrus.Error("[PackageHostVoucherAg]-[OrderVoucherV1.1]")
		logrus.Error(fmt.Sprintf("[Unmarshal]-[Error : %v]", err.Error()))
		logrus.Error("Failed to unmarshaling response OrderVoucher V1.1 from VoucherAG")

		return &resp, err
	}

	return &resp, nil
}

func CheckStatusOrder(req vgmodel.RequestCheckOrderStatus, head models.RequestHeader) (*vgmodel.ResponseCheckOrderStatus, error) {
	var resp vgmodel.ResponseCheckOrderStatus

	logrus.Info("[PackageHostVoucherAg]-[CheckStatusOrder]")

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
		logrus.Error("Check error : ", err.Error())

		return &resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logrus.Error("Failed to unmarshaling response OrderVoucher from VoucherAG ", err.Error())

		return &resp, err
	}

	return &resp, nil

}

func PaymentInfo(productCode string, head models.RequestHeader) (*vgmodel.PaymentInfoResp, error) {
	var resp vgmodel.PaymentInfoResp

	logrus.Info("[PackageHostVoucherAg]-[PaymentInfo]")

	urlSvr := host + endpointPaymentInfo + fmt.Sprintf("?productCode=%v", productCode)

	header := make(http.Header)
	header.Set("InstitutionId", head.InstitutionID)
	header.Set("DeviceId", head.DeviceID)
	header.Set("Geolocation", head.Geolocation)
	header.Set("ChannelId", head.ChannelID)
	header.Set("AppsId", head.AppsID)
	// header.Set("TimeStamp", head.Timestamp)
	// header.Set("Signature", head.Signature)
	header.Set("Authorization", head.Authorization)

	data, err := https.HTTPxGET(urlSvr, header)
	// data, err := HTTPxFormGetVoucherAg(urlSvr, head)
	if err != nil {
		logrus.Error("Check error : ", err.Error())

		return &resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logrus.Error("Failed to unmarshaling response OrderVoucher from VoucherAG ", err.Error())

		return &resp, err
	}

	return &resp, nil

}

func CheckStatusOrderV21(orderId string, head models.RequestHeader) (*vgmodel.ResponseCheckOrderStatus, error) {
	var resp vgmodel.ResponseCheckOrderStatus

	logrus.Info("[PackageHostVoucherAg]-[CheckStatusOrderV21]")

	req := vgmodel.RequestCheckOrderStatusV21{
		OrderID: orderId,
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	signature := utils.CreateSignatureGeneral(timestamp, req, head, 1)

	urlSvr := host + endpointCheckStatusOrderV21

	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	header.Set("InstitutionId", head.InstitutionID)
	header.Set("DeviceId", head.DeviceID)
	header.Set("Geolocation", head.Geolocation)
	header.Set("ChannelId", head.ChannelID)
	header.Set("Signature", signature)
	header.Set("AppsId", head.AppsID)
	header.Set("Timestamp", timestamp)

	data, err := https.HTTPxPOSTwithRequest(urlSvr, req, header)
	// data, err := HTTPxFormGetVoucherAg(urlSvr, head)
	if err != nil {
		logrus.Error("Check error : ", err.Error())

		return &resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logrus.Error("Failed to unmarshaling response OrderVoucher from VoucherAG ", err.Error())

		return &resp, err
	}

	return &resp, nil

}

func CallbackSepulsaVAG(req interface{}) error {
	var resp interface{}

	logrus.Info("[PackageHostVoucherAg]-[CallbackSepulsaVAG.1]")

	urlSvr := host + endpointCallbackSepulsa

	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	header.Set("ChannelId", "H2H")

	data, err := https.HTTPxPOSTwithRequest(urlSvr, req, header)
	// data, err := HTTPxFormPostVoucherAg(urlSvr, head, req)
	if err != nil {

		logrus.Error("[PackageHostVoucherAg]-[OrderVoucherV1.1]")
		logrus.Error(fmt.Sprintf("[HTTPxPOSTwithRequest]-[Error : %v]", err.Error()))

		return err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {

		logrus.Error("[PackageHostVoucherAg]-[OrderVoucherV1.1]")
		logrus.Error(fmt.Sprintf("[Unmarshal]-[Error : %v]", err.Error()))
		logrus.Error("Failed to unmarshaling response OrderVoucher V1.1 from VoucherAG")

		return err
	}

	return nil
}
