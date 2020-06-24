package publisher

import (
	"fmt"

	"github.com/astaxie/beego"
	jsoniter "github.com/json-iterator/go"
)

type PublishReq struct {
	Topic string `json:"topic"`
	Value []byte `json:"value"`
}

var (
	host              string
	endpointPublisher string
)

func init() {
	host = beego.AppConfig.DefaultString("host.publish", "http://13.228.25.85:8703")
	endpointPublisher = beego.AppConfig.DefaultString("endpoint.publish", "/ottopoint/v0.1.0/kafka/publish")
}

// SendPublishKafka ...
func SendPublishKafka(request PublishReq) ([]byte, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	datareq, _ := json.Marshal(request)

	url := host + endpointPublisher

	data, err := HTTPPostKafka(url, request)
	fmt.Println("xxxx-----------xxxx")
	fmt.Println("urlSvr", url)
	fmt.Println("msgreq", request)
	fmt.Println("datareq", string(datareq))
	fmt.Println("err", err)
	fmt.Println("xxxx-----------xxxx")

	return data, err
}
