package check_status

import (
	"fmt"
	"ottopoint-purchase/db"
	vag "ottopoint-purchase/hosts/voucher_aggregator/host"
	"ottopoint-purchase/models"

	"github.com/sirupsen/logrus"
)

func SchedulerCheckStatusServiceV21() interface{} {
	res := models.ResponseData{}

	res.ResponseCode = "00"
	res.ResponseDesc = "Success"

	savericename := "[PackageCheckStatus]-[SchedulerCheckStatusServiceV21]"

	getData, errData := db.GetDataScheduler()
	if errData != nil || len(getData) == 0 {

		logrus.Error(savericename)
		logrus.Error(fmt.Sprintf("[GetDataScheduler]-[Error : %v]", errData))

		res.ResponseCode = "153"
		res.ResponseDesc = "Data Not Found"

		return res
	}

	count := len(getData)

	var total, failed, success int

	for i := 0; i < count; i++ {

		total = getData[i].Count

		var institution string
		dataTrxId, errTrxId := db.CheckTrxbyTrxID(getData[i].TransactionID)

		institution = dataTrxId.Institution

		if errTrxId != nil {

			logrus.Error(savericename)
			logrus.Error(fmt.Sprintf("[CheckTrxbyTrxID]-[Error : %v]", errTrxId))

			institution = ""

		}

		head := models.RequestHeader{
			DeviceID:      "-",
			InstitutionID: institution,
			Geolocation:   "-",
			ChannelID:     "H2H",
			AppsID:        "-",
			// Timestamp    : "-",
			// Authorization: "-",
			// Signature    : "-",
		}

		fmt.Println(">>> Supplier : ", dataTrxId.SupplierID)
		checkStatus, errStatus := vag.CheckStatusOrderV21(getData[i].TransactionID, head)
		if errStatus != nil || checkStatus.ResponseCode != "00" {

			total = total + 1

			fmt.Println(">>> Supplier : ", dataTrxId.SupplierID)
			logrus.Error(savericename)
			logrus.Error(fmt.Sprintf("[CheckStatusOrderV21]-[Error : %v]", errStatus))

			go db.UpdateSchedulerStatus(false, total, getData[i].TransactionID)

			failed++
			continue
		}

		success++
		go db.UpdateSchedulerStatus(true, total, getData[i].TransactionID)

	}

	res = models.ResponseData{
		Data: map[string]interface{}{
			"Failed":  failed,
			"Success": success,
			"Total":   count,
		},
	}

	return res

}
