package v2_migrate

import (
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"

	"ottopoint-purchase/hosts/opl/host"
	"ottopoint-purchase/models"
	"ottopoint-purchase/models/dbmodels"
	"ottopoint-purchase/services"
	"ottopoint-purchase/utils"
	"strconv"
	"time"

	"github.com/vjeantet/jodaTime"
)

func Adding_PointVoucher(param models.Params, countPoint, countVoucher int) string {
	fmt.Println("[ >>>>>>>>>>>>>>>>>> Adding / Reversal Point and Voucher <<<<<<<<<<<<<<<< ]")

	var result string

	textComment := param.TrxID + param.InstitutionID + constants.CodeReversal + "#" + "OP009 - Reversal point cause transaction " + param.NamaVoucher + " is failed"
	statusEarning := constants.Success
	msgEarning := constants.MsgSuccess
	// save to scheduler
	schedulerData := dbmodels.TSchedulerRetry{
		// ID
		Code:          constants.CodeScheduler,
		TransactionID: utils.Before(textComment, "#"),
		Count:         0,
		IsDone:        false,
		CreatedAT:     time.Now(),
		// UpdatedAT
	}

	// Adding/Reversal Point
	sendReversal, errReversal := host.TransferPoint(param.AccountId, strconv.Itoa(countPoint), textComment)

	if errReversal != nil || sendReversal.PointsTransferId == "" {
		statusEarning = constants.TimeOut

		fmt.Println(fmt.Sprintf("===== Failed TransferPointOPL to %v || RRN : %v =====", param.AccountNumber, param.RRN))
		for _, val1 := range sendReversal.Form.Children.Customer.Errors {
			if val1 != "" {
				msgEarning = val1
				statusEarning = constants.Failed
			}

		}

		for _, val2 := range sendReversal.Form.Children.Points.Errors {
			if val2 != "" {
				msgEarning = val2
				statusEarning = constants.Failed
			}
		}

		if sendReversal.Message != "" {
			msgEarning = sendReversal.Message
			statusEarning = constants.Failed
		}

		if sendReversal.Error.Message != "" {
			msgEarning = sendReversal.Error.Message
			statusEarning = constants.Failed
		}

		if statusEarning == constants.TimeOut {
			errSaveScheduler := db.DbCon.Create(&schedulerData).Error
			if errSaveScheduler != nil {

				fmt.Println("===== Gagal SaveScheduler ke DB =====")
				fmt.Println(fmt.Sprintf("Error : %v", errSaveScheduler))
				fmt.Println(fmt.Sprintf("===== Phone : %v || RRN : %v =====", param.AccountNumber, param.RRN))

				// return
			}

		}
	}

	expired := services.ExpiredPointService()
	saveReversal := dbmodels.TEarning{
		ID: utils.GenerateTokenUUID(),
		// EarningRule     :,
		// EarningRuleAdd  :,
		PartnerId: param.InstitutionID,
		// ReferenceId     : ,
		TransactionId: param.TrxID,
		// ProductCode     :,
		// ProductName     :,
		AccountNumber: param.AccountNumber,
		// Amount          :,
		Point:   int64(countPoint),
		Commnet: textComment,
		// Remark          :,
		Status:           statusEarning,
		StatusMessage:    msgEarning,
		PointsTransferId: sendReversal.PointsTransferId,
		// RequestorData   :,
		// ResponderData   :,
		TransType:       constants.CodeReversal,
		AccountId:       param.AccountId,
		ExpiredPoint:    expired,
		TransactionTime: time.Now(),
	}

	errSaveReversal := db.DbCon.Create(&saveReversal).Error
	if errSaveReversal != nil {

		fmt.Println(fmt.Sprintf("[Failed Save Reversal to DB]-[Error : %v]", errSaveReversal))
		fmt.Println("[PackageServices]-[SaveEarning]")

		fmt.Println(">>> Save CSV <<<")
		name := jodaTime.Format("dd-MM-YYYY", time.Now()) + ".csv"
		go utils.CreateCSVFile(saveReversal, name)

		result = ">>>>> Create CreateCSVFile <<<<"
		// return result

	}

	// Adding/Reversal Voucher stock
	dtaVocher, _ := db.Get_MReward(param.RewardID)
	latestUsageLimit := dtaVocher.UsageLimit + countVoucher

	errDeductVouch := db.UpdateUsageLimitVoucher(param.RewardID, latestUsageLimit)
	if errDeductVouch != nil {

		result = ">>>>>>> Failed Get Voucher <<<<<<<<"
		return result
	}

	result = ">>>>>>> Adding/Reversal Point Success <<<<<<<<"
	return result

}
