package v2_migrate

import (
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	opl "ottopoint-purchase/hosts/opl/host"
	"ottopoint-purchase/models"
	"strconv"
	"strings"
)

func Redeem_PointandVoucher(QtyVoucher int, param models.Params) (models.SpendingPointVoucher, error) {
	fmt.Println("[ >>>>>>>>>>>>>>>>>> Redeem Point and Voucher <<<<<<<<<<<<<<<< ]")
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
	textComment := param.NamaVoucher + "," + "product code : " + param.ProductCodeInternal
	fmt.Println("Comment Spending Point Redeem Voucher : ", textComment)
	replCostPoint := strings.ReplaceAll(strconv.Itoa(param.Point), ",", ".")
	fmt.Println("Cost Point Voucher before : ", strconv.Itoa(param.Point))
	fmt.Println("Cost point Voucher : ", replCostPoint)
	resSpend, errSpend := opl.SpendPoint(param.AccountId, replCostPoint, textComment)

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
		fmt.Println("Result Error Spending Point :")
		fmt.Println(msgEarning)
		fmt.Println(statusEarning)
		//

		result.Rc = "500"
		result.Rd = "Internal Server Error"
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

	result.Rc = "00"
	result.Rd = "Success"
	result.CouponsCode = param.CouponCode
	return result, nil

}
