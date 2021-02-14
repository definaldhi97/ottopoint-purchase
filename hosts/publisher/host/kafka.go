package publisher

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego/logs"
	jsoniter "github.com/json-iterator/go"

	https "ottopoint-purchase/hosts"
	"ottopoint-purchase/utils"
)

type PublishReq struct {
	Topic string `json:"topic"`
	Value []byte `json:"value"`
}

type PublisherResp struct {
	ResponseCode string `json:"responseCode"`
	ResponseDesc string `json:"responseDesc"`
}

var (
	host              string
	endpointPublisher string
)

func init() {
	host = utils.GetEnv("HOST_PUBLISHER", "http://13.228.25.85:8703")
	endpointPublisher = utils.GetEnv("endpoint.publish", "/ottopoint/v0.1.0/kafka/publish")
}

// SendPublishKafka ...
func SendPublishKafka(request PublishReq) (PublisherResp, error) {

	var resp PublisherResp

	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	datareq, _ := json.Marshal(request)

	url := host + endpointPublisher

	header := make(http.Header)
	header.Set("Content-Type", "application/json")

	data, err := https.HTTPxPOSTwithRequest(url, request, header)
	fmt.Println("xxxx-----------xxxx")
	fmt.Println("urlSvr", url)
	fmt.Println("msgreq", request)
	fmt.Println("datareq", string(datareq))
	fmt.Println("err", err)
	fmt.Println("xxxx-----------xxxx")

	err = json.Unmarshal(data, &resp)
	if err != nil {
		logs.Error("Failed to unmarshaling response SendPublishKafka from Publisher ", err.Error())

		return resp, err
	}

	return resp, err
}
