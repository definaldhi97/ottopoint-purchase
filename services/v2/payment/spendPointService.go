package payment

import (
	"encoding/json"
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	"ottopoint-purchase/hosts/op_corepoint"
	"ottopoint-purchase/models"
	"ottopoint-purchase/models/dbmodels"
	"ottopoint-purchase/utils"
	"time"
)

func SpendPointService(param models.Params, header models.RequestHeader) (models.SpendingPointVoucher, error) {

	result := models.SpendingPointVoucher{}

	var msgEarning, statusEarning string

	spenPoinReq := op_corepoint.SpendingPointReq{
		AccountID:     param.AccountId,
		TransactionID: utils.Before(param.Comment, "#"),
		Point:         param.Point,
		Comment:       param.Comment,
	}

	// save to scheduler
	schedulerData := dbmodels.TSchedulerRetry{
		// ID
		Code:          constants.CodeSchedulerSpending,
		TransactionID: utils.Before(param.Comment, "#"),
		Count:         0,
		IsDone:        false,
		CreatedAT:     time.Now(),
		// UpdatedAT
	}

	// Spending/Deduct Point
	resSpend, errSpend := op_corepoint.SependingPoint(spenPoinReq, header)

	if errSpend != nil {

		statusEarning = constants.CODE_STATUS_TO
		msgEarning = "Internal Server Error"
		errSaveScheduler := db.DbCon.Create(&schedulerData).Error
		if errSaveScheduler != nil {
			fmt.Println("===== Gagal SaveScheduler ke DB =====")
			fmt.Println(fmt.Sprintf("Error : %v", errSaveScheduler))
			fmt.Println(fmt.Sprintf("===== Phone : %v || RRN : %v =====", param.AccountNumber, param.RRN))
		}

		result.Rc = statusEarning
		result.Rd = msgEarning
		return result, errSpend
	}

	if resSpend.ResponseCode != "00" {
		result.Rc = resSpend.ResponseCode
		result.Rd = resSpend.ResponseDesc
		return result, errSpend
	}

	dataSpend := op_corepoint.DataSpendingPoint{}
	jsonString, _ := json.Marshal(resSpend.Data)
	json.Unmarshal(jsonString, &dataSpend)

	result.Rc = "00"
	result.Rd = "Success"
	result.PointTransferID = dataSpend.PointsTransferID

	return result, nil

}
