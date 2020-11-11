package schedulers

import (
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"

	"github.com/opentracing/opentracing-go"
)

type SchedulerCheckStatusService struct {
	General models.GeneralModel
}

func (t SchedulerCheckStatusService) NewSchedulerCheckStatusService() interface{} {
	res := models.ResponseData{}

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[SchedulerCheckStatusService]")

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[SchedulerCheckStatusService]")
	defer span.Finish()

	fmt.Println(">>> [Start]-[SchedulerCheckStatusService] <<<")

	getData, errData := db.GetDataScheduler()
	if errData != nil || len(getData) == 0 {

		fmt.Println(fmt.Sprintf("[Failed to GetDataScheduler]-[Error : %v]", errData))
		fmt.Println("[PackageServices]-[GetDataScheduler]")

		sugarLogger.Info(fmt.Sprintf("[Failed to GetDataScheduler]-[Error : %v]", errData))
		sugarLogger.Info("[PackageServices]-[GetDataScheduler]")

		res.ResponseCode = "153"
		res.ResponseDesc = "Data Not Found"

		return res
	}

	count := len(getData)

	csd := []models.SchedulerCheckStatusData{}

	var sp, fp, tp int
	var supplierSepulsa string
	// var supplierUV, supplierOttoAG string

	for i := 0; i < count; i++ {

		if getData[i].Code == constants.CodeSchedulerSepulsa {
			supplierSepulsa = "Sepulsa"
			total := getData[i].Count

			errSepulsa := t.CheckStatusSepulsaServices(utils.Before(getData[i].TransactionID, "PSM"))
			if errSepulsa != nil {

				total = total + 1

				fmt.Println(fmt.Sprintf("[Error from CheckStatusSepulsaServices]-[Error : %v]", errSepulsa))
				fmt.Println("[PackageServices]-[CheckStatusSepulsaServices]")

				sugarLogger.Info(fmt.Sprintf("[Error from CheckStatusSepulsaServices]-[Error : %v]", errSepulsa))
				sugarLogger.Info("[PackageServices]-[CheckStatusSepulsaServices]")

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

	// if supplierUV != "" {
	// 	dataSepulsa := models.SchedulerCheckStatusData{
	// 		Supplier : supplierSepulsa,
	// 		Success  : sp,
	// 		Failed   : fp,
	// 		Total    : tp,
	// 	}

	// 	csd = append(csd, dataSepulsa)
	// }

	// if supplierOttoAG != "" {
	// 	dataSepulsa := models.SchedulerCheckStatusData{
	// 		Supplier : supplierSepulsa,
	// 		Success  : sp,
	// 		Failed   : fp,
	// 		Total    : tp,
	// 	}

	// 	csd = append(csd, dataSepulsa)
	// }

	resData := models.SchedulerCheckStatusResp{
		Data:  csd,
		Total: count,
	}

	return resData

}
