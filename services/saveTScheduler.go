package services

import (
	"fmt"
	"ottopoint-purchase/db"
	"ottopoint-purchase/models/dbmodels"
	"time"
)

func SaveTSchedulerRetry(trxID, code string) {

	fmt.Println(fmt.Sprintf("[Start-SaveTSchedulerRetry][TrxId : %v]", trxID))

	schedulerData := dbmodels.TSchedulerRetry{
		Code:          code,
		TransactionID: trxID,
		Count:         0,
		IsDone:        false,
		CreatedAT:     time.Now(),
	}

	errSaveScheduler := db.DbCon.Create(&schedulerData).Error
	if errSaveScheduler != nil {

		fmt.Println("===== Gagal SaveScheduler ke DB =====")
		fmt.Println(fmt.Sprintf("Error : %v", errSaveScheduler))

	}

}
