package check_status

import (
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	vag "ottopoint-purchase/hosts/voucher_aggregator/host"
	"ottopoint-purchase/models"

	"github.com/sirupsen/logrus"
)

func SchedulerCheckStatusServiceV21() interface{} {
	res := models.ResponseData{}

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

	var sepulsaTotal, voucherAgTotal, ottoAgTotal, uvTotal, jkTotal, gvTotal int
	var supplier string
	var total int

	for i := 0; i < count; i++ {

		var institution string
		dataTrxId, errTrxId := db.CheckTrxbyTrxID(getData[i].TransactionID)
		if errTrxId != nil {

			logrus.Error(savericename)
			logrus.Error(fmt.Sprintf("[CheckTrxbyTrxID]-[Error : %v]", errTrxId))

			institution = ""

		}

		switch getData[i].Code {
		case constants.CodeSchedulerSepulsa:
			supplier = constants.Sepulsa
			sepulsaTotal++
		case constants.CodeSchedulerVoucherAG:
			supplier = constants.VoucherAg
			voucherAgTotal++
		case constants.CodeSchedulerOttoAG:
			supplier = constants.OttoAG
			ottoAgTotal++
		case constants.CodeSchedulerUV:
			supplier = constants.UV
			uvTotal++
		case constants.CodeSchedulerJempolKios:
			supplier = constants.JempolKios
			jkTotal++
		case constants.CodeSchedulerGudangVoucher:
			supplier = constants.GudangVoucher
			gvTotal++
		}

		institution = dataTrxId.Institution

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

		fmt.Println(">>> Supplier : ", supplier)
		_, errStatus := vag.CheckStatusOrderV21(getData[i].TransactionID, head)
		if errStatus != nil {

			fmt.Println(">>> Supplier : ", supplier)
			logrus.Error(savericename)
			logrus.Error(fmt.Sprintf("[CheckStatusOrderV21]-[Error : %v]", errStatus))

			go db.UpdateSchedulerStatus(false, total, getData[i].TransactionID)

			continue
		}

	}

	resData := models.SchedulerCheckStatusResp{
		Data: models.SchedulerCheckStatusDataSupplier{
			Sepulsa:       sepulsaTotal,
			OttoAG:        ottoAgTotal,
			UltraVoucher:  uvTotal,
			JempolKios:    jkTotal,
			GudangVoucher: gvTotal,
			VouicherAG:    voucherAgTotal,
		},
		Total: count,
	}

	return resData

}
