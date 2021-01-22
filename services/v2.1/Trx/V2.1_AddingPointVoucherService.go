package Trx

import (
	"encoding/json"
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	v2_trx "ottopoint-purchase/services/v2/Trx"

	"ottopoint-purchase/hosts/op_corepoint"
	"ottopoint-purchase/models"
	"ottopoint-purchase/models/dbmodels"
	"ottopoint-purchase/utils"
	"time"

	"github.com/vjeantet/jodaTime"
)

func V21_Adding_PointVoucher(param models.Params, countPoint, countVoucher int, header models.RequestHeader) string {
	fmt.Println("[ >>>>>>>>>>>>>>>>>> Adding / Reversal Point and Voucher <<<<<<<<<<<<<<<< ]")

	var result string

	textComment := param.TrxID + param.InstitutionID + constants.CodeReversal + "#" + "OP009 - Reversal point cause transaction " + param.NamaVoucher + " is failed"
	statusEarning := constants.Success
	msgEarning := constants.MsgSuccess

	// Get expiredDays point
	expPoint, _ := db.ParamData(constants.CODE_CONFIG_COREPOINT_GROUP, constants.CODE_CONFIG_COREPOINT_EXPIRED_POINT_DAYS)

	addingPoinReq := op_corepoint.AddingPointReq{
		AccountID:     param.AccountId,
		TransactionID: param.TrxID,
		ExpiredDays:   expPoint.Value,
		Point:         countPoint,
		Comment:       textComment,
	}

	// save to scheduler
	schedulerData := dbmodels.TSchedulerRetry{
		// ID
		Code:          constants.CodeScheduler,
		TransactionID: param.TrxID,
		Count:         0,
		IsDone:        false,
		CreatedAT:     time.Now(),
		// UpdatedAT
	}

	// Adding/Reversal Point
	// sendReversal, errReversal := host.TransferPoint(param.AccountId, strconv.Itoa(countPoint), textComment)
	addingPoint, errAdding := op_corepoint.AddingPoint(addingPoinReq, header)

	if errAdding != nil {
		statusEarning = constants.TimeOut
		msgEarning = "Internal Server Error"
		errSaveScheduler := db.DbCon.Create(&schedulerData).Error

		if errSaveScheduler != nil {
			fmt.Println("===== Gagal SaveScheduler ke DB =====")
			fmt.Println(fmt.Sprintf("Error : %v", errSaveScheduler))
			fmt.Println(fmt.Sprintf("===== Phone : %v || RRN : %v =====", param.AccountNumber, param.RRN))
			// return
		}
	}

	if addingPoint.ResponseCode != "00" && errAdding == nil {
		statusEarning = addingPoint.ResponseCode
		msgEarning = addingPoint.ResponseDesc
	}

	label := op_corepoint.DataAddingPoint{}
	jsonString, _ := json.Marshal(addingPoint.Data)
	json.Unmarshal(jsonString, &label)

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
		PointsTransferId: label.PointsTransferID,
		// RequestorData   :,
		// ResponderData   :,
		TransType:       constants.CodeReversal,
		AccountId:       param.AccountId,
		ExpiredPoint:    label.ExpiredPoint.Format("2006-01-02"),
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

	// errDeductVouch := db.UpdateUsageLimitVoucher(param.RewardID, latestUsageLimit)
	errDeductVouch := v2_trx.UpdateUsageLimitVoucher(param.RewardID, latestUsageLimit)
	if errDeductVouch != nil {

		result = ">>>>>>> Failed Update Stock Voucher <<<<<<<<"
		return result
	}

	result = ">>>>>>> Adding/Reversal Point Success <<<<<<<<"
	return result

}
