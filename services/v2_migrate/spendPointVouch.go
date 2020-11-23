package v2_migrate

import (
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	opl "ottopoint-purchase/hosts/opl/host"
	"ottopoint-purchase/models"
	"ottopoint-purchase/models/dbmodels"
	"ottopoint-purchase/utils"
	"strconv"
	"strings"
	"time"
)

func Redeem_PointandVoucher(QtyVoucher int, param models.Params) (models.SpendingPointVoucher, error) {
	fmt.Println("[ >>>>>>>>>>>>>>>>>> Spending/Deduct Point and Voucher <<<<<<<<<<<<<<<< ]")
	var result models.SpendingPointVoucher
	var msgEarning, statusEarning string
	// validasi usage limit voucher

	fmt.Println("Qty req Voucher : ", QtyVoucher)
	fmt.Println("Stock Voucher Available : ", param.UsageLimitVoucher)
	if QtyVoucher > param.UsageLimitVoucher {
		fmt.Println("[ Stock Voucher not Available ]")
		result.Rc = constants.RC_VOUCHER_NOT_AVAILABLE
		result.Rd = constants.RD_VOUCHER_NOT_AVAILABLE
		return result, nil
	}
	// deduct/spending point
	// textComment := param.NamaVoucher + "," + "product code : " + param.ProductCodeInternal
	// textComment := param.Reffnum + "#" + param.NamaVoucher
	totalPoint := param.Point * QtyVoucher
	fmt.Println("Comment Spending Point Redeem Voucher : ", param.Comment)
	replCostPoint := strings.ReplaceAll(strconv.Itoa(totalPoint), ",", ".")
	fmt.Println("Cost Point Voucher before : ", strconv.Itoa(totalPoint))
	fmt.Println("Cost point Voucher : ", replCostPoint)

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

	resSpend, errSpend := opl.SpendPoint(param.AccountId, replCostPoint, param.Comment)

	if errSpend != nil || resSpend.PointsTransferId == "" {

		statusEarning = constants.CODE_STATUS_TO
		msgEarning = "Internal Server Error"

		for _, val1 := range resSpend.Form.Children.Customer.Errors {
			if val1 != "" {
				msgEarning = val1
				statusEarning = constants.CODE_FAILED
			}
		}

		for _, val2 := range resSpend.Form.Children.Points.Errors {
			if val2 != "" {
				msgEarning = val2
				statusEarning = constants.CODE_FAILED
			}
		}

		if resSpend.Message != "" {
			msgEarning = resSpend.Message
			statusEarning = constants.CODE_FAILED
		}

		// check scheduler
		if statusEarning == constants.TimeOut {
			errSaveScheduler := db.DbCon.Create(&schedulerData).Error
			if errSaveScheduler != nil {

				fmt.Println("===== Gagal SaveScheduler ke DB =====")
				fmt.Println(fmt.Sprintf("Error : %v", errSaveScheduler))
				fmt.Println(fmt.Sprintf("===== Phone : %v || RRN : %v =====", param.AccountNumber, param.RRN))

				// return
			}
		}

		result.Rc = "500"
		result.Rd = msgEarning
		return result, errSpend
	}

	// update Usage limit voucher by rewardID
	latestUsageLimitVouch := param.UsageLimitVoucher - QtyVoucher
	errDeductVouch := db.UpdateUsageLimitVoucher(param.CampaignID, latestUsageLimitVouch)
	if errDeductVouch != nil {
		result.Rc = "500"
		result.Rd = "Internal Server Error"
		return result, errDeductVouch
	}

	generateCouponsID := utils.GenerateTokenUUID()
	fmt.Println("Coupons ID : ", generateCouponsID)

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
