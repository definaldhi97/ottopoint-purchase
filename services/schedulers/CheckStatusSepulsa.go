package schedulers

import (
	"fmt"
	sepulsa "ottopoint-purchase/hosts/sepulsa/host"

	"github.com/opentracing/opentracing-go"
)

func (t SchedulerCheckStatusService) CheckStatusSepulsaServices(trxid string) error {
	// res := models.SchedulerCheckStatusData{}

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info(">>> [Start]-[CheckStatusSepulsaServices] <<<")

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[CheckStatusSepulsaServices]")
	defer span.Finish()

	fmt.Println(">>> [Start]-[CheckStatusSepulsaServices] <<<")

	// check status ke sepulsa
	_, errStatus := sepulsa.EwalletDetailTransaction(trxid)
	if errStatus != nil {

		fmt.Println(fmt.Sprintf("[Error from EwalletDetailTransaction]-[Error : %v]", errStatus))
		fmt.Println("[PackageServices]-[EwalletDetailTransaction]")

		sugarLogger.Info(fmt.Sprintf("[Error from EwalletDetailTransaction]-[Error : %v]", errStatus))
		sugarLogger.Info("[PackageServices]-[EwalletDetailTransaction]")

		return errStatus
	}

	return nil

}
