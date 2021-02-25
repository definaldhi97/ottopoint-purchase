package check_status

import (
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	"ottopoint-purchase/models"

	"github.com/sirupsen/logrus"
)

func SchedulerCheckStatusService() interface{} {
	res := models.ResponseData{}

	savericename := "[PackageCheckStatus]-[SchedulerCheckStatusService]"

	getData, errData := db.GetDataScheduler()
	if errData != nil || len(getData) == 0 {

		logrus.Error(savericename)
		logrus.Error(fmt.Sprintf("[GetDataScheduler]-[Error : %v]", errData))

		res.ResponseCode = "153"
		res.ResponseDesc = "Data Not Found"

		return res
	}

	count := len(getData)

	csd := []models.SchedulerCheckStatusData{}

	var sp, fp, tp int
	var supplierSepulsa, supplierVoucherAG string

	for i := 0; i < count; i++ {

		if getData[i].Code == constants.CodeSchedulerSepulsa {
			supplierSepulsa = constants.Sepulsa
			total := getData[i].Count

			// errSepulsa := t.CheckStatusSepulsaServices(utils.Before(getData[i].TransactionID, "PSM"))
			errSepulsa := CheckStatusSepulsaServices(getData[i].TransactionID)
			if errSepulsa != nil {

				total = total + 1

				logrus.Error(savericename)
				logrus.Error(fmt.Sprintf("[CheckStatusSepulsaServices]-[Error : %v]", errSepulsa))

				go db.UpdateSchedulerStatus(false, total, getData[i].TransactionID)

				fp++
				tp++
				continue
			}

			go db.UpdateSchedulerStatus(true, total, getData[i].TransactionID)

			sp++
			tp++
			continue
		}

		if getData[i].Code == constants.CodeSchedulerVoucherAG {
			supplierVoucherAG = constants.VoucherAg
			total := getData[i].Count

			// errVaG := t.CheckStatusVoucherAgService(utils.Before(getData[i].TransactionID, "PSM"))
			errVaG := CheckStatusVoucherAgService(getData[i].TransactionID)
			if errVaG != nil {

				total = total + 1

				logrus.Error(savericename)
				logrus.Error(fmt.Sprintf("[CheckStatusVoucherAgService]-[Error : %v]", errVaG))

				go db.UpdateSchedulerStatus(false, total, getData[i].TransactionID)

				fp++
				tp++
				continue
			}

			go db.UpdateSchedulerStatus(true, total, getData[i].TransactionID)

			sp++
			tp++
			continue
		}

	}

	if supplierSepulsa != "" {
		dataSepulsa := models.SchedulerCheckStatusData{
			Supplier: supplierSepulsa,
			Success:  sp,
			Failed:   fp,
			Total:    tp,
		}

		csd = append(csd, dataSepulsa)
	}

	if supplierVoucherAG != "" {
		dataVoucherAG := models.SchedulerCheckStatusData{
			Supplier: supplierVoucherAG,
			Success:  sp,
			Failed:   fp,
			Total:    tp,
		}

		csd = append(csd, dataVoucherAG)
	}

	resData := models.SchedulerCheckStatusResp{
		Data:  csd,
		Total: count,
	}

	return resData

}
