package v2_migrate

import (
	"fmt"
	"ottopoint-purchase/db"
	"ottopoint-purchase/models"
	"ottopoint-purchase/services"
	"sync"
)

func UseVoucherOttoAg(req models.VoucherComultaiveReq, redeemComu models.RedeemComuResp, param models.Params, getRespChan chan models.RedeemResponse, ErrRespUseVouc chan error, wg *sync.WaitGroup) {
	fmt.Println("[ >>>>>>>>>>>>>>>>>> Use Voucher Comulative Voucher Otto AG <<<<<<<<<<<<<<<< ]")

	defer wg.Done()
	defer close(getRespChan)
	defer close(ErrRespUseVouc)

	// get CustID
	dataUser, errUser := db.CheckUser(param.AccountNumber)
	if errUser != nil {
		fmt.Println("User Not Eligible, Error : ", errUser)
	} else {
		fmt.Println("[ CampaignID ] :", req.CampaignID)
		fmt.Println("[ Coupon ID ] : ", redeemComu.CouponID)
		fmt.Println("[ Coupon Code ] : ", redeemComu.CouponCode)
		fmt.Println("[ Used ] : ", 1)
		fmt.Println("[ Customer ID ] : ", dataUser.CustID)
		fmt.Println("[ Category voucher ] : ", param.Category)
	}

	// Reedem Use Voucher
	param.Amount = redeemComu.Redeem.Amount
	param.RRN = redeemComu.Redeem.Rrn
	param.CouponID = redeemComu.CouponID
	resRedeem := services.RedeemUseVoucherComulative(req, param)
	getRespChan <- resRedeem
}
