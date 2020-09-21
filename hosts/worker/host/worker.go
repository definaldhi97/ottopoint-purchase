package host

import (
	"encoding/json"
	"fmt"
	"ottopoint-purchase/hosts/worker/models"

	ODU "ottodigital.id/library/utils"
)

var (
	host string
	name string

	endpointEarning string

	// healthCheckKey string
)

func init() {
	host = ODU.GetEnv("OTTOPOINT_HOST_WORKER", "http://13.228.25.85:8011")
	name = ODU.GetEnv("OTTOPOINT_NAME_WORKER", "OTTOMART")

	endpointEarning = ODU.GetEnv("OTTOPOINT_ENDPOINT_WORKER", "/ottopoint-worker-earning/earningPoint")

}

func WorkerEarning(req models.WorkerEarningReq) (*models.WorkerEarningResp, error) {
	var resp models.WorkerEarningResp

	fmt.Println("[Package Host Worker]-[WorkerEarning]")

	urlSvr := host + endpointEarning

	data, err := HTTPxFormWithBody(urlSvr, req)
	if err != nil {
		fmt.Println("Check error : ", err)

		return &resp, err
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {

		fmt.Println("Failed to unmarshaling response WorkerEarning from worker ", err.Error())

		return &resp, err
	}

	return &resp, nil
}
