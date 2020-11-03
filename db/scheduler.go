package db

import (
	"fmt"
	"ottopoint-purchase/models/dbmodels"
)

func GetDataScheduler() ([]dbmodels.TSchedulerRetry, error) {
	res := []dbmodels.TSchedulerRetry{}

	err := DbCon.Raw(`select * from t_scheduler_retry where is_done = false order by id asc`).Scan(&res).Error
	if err != nil {

		fmt.Println(fmt.Sprintf("[Failed to GetDataScheduler]-[Error : %v]", err))
		fmt.Println("[PackageDB]-[GetDataScheduler]")

		return res, err
	}

	return res, nil
}
