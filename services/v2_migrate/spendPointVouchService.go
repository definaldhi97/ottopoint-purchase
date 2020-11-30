package v2_migrate

import (
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	"ottopoint-purchase/hosts/op_corepoint"
	"ottopoint-purchase/models"
	"ottopoint-purchase/models/dbmodels"
	"ottopoint-purchase/utils"
	"time"

	"github.com/sirupsen/logrus"
)

func Redeem_PointandVoucher(QtyVoucher int, param models.Params, TrxID string, header models.RequestHeader) (models.SpendingPointVoucher, error) {
	fmt.Println("[ >>>>>>>>>>>>>>>>>> Spending/Deduct Point and Voucher <<<<<<<<<<<<<<<< ]")
	var result models.SpendingPointVoucher
	var msgEarning, statusEarning string
	// validasi usage limit voucher

	logrus.Info("Qty req Voucher : ", QtyVoucher)

	dtaVocher, _ := db.Get_MReward(param.CampaignID)
	logrus.Info("Stock Voucher Available : ", dtaVocher.UsageLimit)

	if QtyVoucher > dtaVocher.UsageLimit {

		logrus.Info("[ Stock Voucher not Available ]")
		result.Rc = constants.RC_VOUCHER_NOT_AVAILABLE
		result.Rd = constants.RD_VOUCHER_NOT_AVAILABLE
		return result, nil
	}

	// deduct/spending point
	// textComment := param.NamaVoucher + "," + "product code : " + param.ProductCodeInternal
	// textComment := param.Reffnum + "#" + param.NamaVoucher
	totalPoint := param.Point * QtyVoucher
	logrus.Info("Comment Spending Point Redeem Voucher : ", param.Comment)

	spenPoinReq := op_corepoint.SpendingPointReq{
		AccountID:     param.AccountId,
		TransactionID: TrxID,
		Point:         totalPoint,
		Comment:       param.Comment,
	}

	// save to scheduler
	schedulerData := dbmodels.TSchedulerRetry{
		// ID
		Code:          constants.CodeSchedulerSpending,
		TransactionID: TrxID,
		Count:         0,
		IsDone:        false,
		CreatedAT:     time.Now(),
		// UpdatedAT
	}

	// Spending/Deduct Point
	resSpend, errSpend := op_corepoint.SpendingPoint(spenPoinReq, header)

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

	// update Usage limit voucher by rewardID
	latestUsageLimitVouch := dtaVocher.UsageLimit - QtyVoucher
	errDeductVouch := db.UpdateUsageLimitVoucher(param.CampaignID, latestUsageLimitVouch)
	if errDeductVouch != nil {
		result.Rc = "500"
		result.Rd = "Internal Server Error"
		return result, errDeductVouch
	}

	// generateCouponsID := utils.GenerateTokenUUID()
	// fmt.Println("Coupons ID : ", generateCouponsID)

	var couponVouc []models.CouponsVoucher
	for i := 0; QtyVoucher > i; i++ {
		a := models.CouponsVoucher{
			CouponsCode: param.CouponCode,
			CouponsID:   utils.GenerateTokenUUID(),
		}
		couponVouc = append(couponVouc, a)
	}

	result.Rc = "00"
	result.Rd = "Success"
	result.CouponseVouch = couponVouc
	return result, nil

}
