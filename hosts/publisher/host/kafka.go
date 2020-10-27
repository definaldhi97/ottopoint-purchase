package publisher

import (
	"fmt"
	"time"

	hcmodels "ottodigital.id/library/healthcheck/models"

	jsoniter "github.com/json-iterator/go"
	ODU "ottodigital.id/library/utils"
)

type PublishReq struct {
	Topic string `json:"topic"`
	Value []byte `json:"value"`
}

var (
	host              string
	name              string
	endpointPublisher string
)

func init() {
	host = ODU.GetEnv("HOST_PUBLISHER", "http://13.228.25.85:8703")
	name = ODU.GetEnv("NAME_PUBLISHER", "Ottopoint-worker-earning")
	endpointPublisher = ODU.GetEnv("endpoint.publish", "/ottopoint/v0.1.0/kafka/publish")
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

// GetServiceHealthCheck ..
func GetServiceHealthCheckPublisher() hcmodels.ServiceHealthCheck {
	return hcmodels.ServiceHealthCheck{
		Name:    name,
		Address: host,
		Status:  "OK",
		// Description: ,
		UpdatedAt: time.Now(),
	}
}
