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

func UpdateSchedulerStatus(status bool, count int, trxId string) error {
	res := dbmodels.TSchedulerRetry{}

	err := DbCon.Raw(`update t_scheduler_retry set is_done = ?, count = ? where coupon_id = ?`, status, count, trxId).Scan(&res).Error
	if err != nil {

		fmt.Println(fmt.Sprintf("[Failed to UpdateSchedulerStatus]-[Error : %v]", err))
		fmt.Println("[PackageDB]-[UpdateSchedulerStatus]")

		return err
	}

	return nil
}
