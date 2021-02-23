package redeemtion

import (
	"fmt"
	"ottopoint-purchase/db"
	"ottopoint-purchase/models"
	"ottopoint-purchase/services"
	"sync"

	"github.com/sirupsen/logrus"
)

func UseVoucherOttoAg(req models.VoucherComultaiveReq, redeemComu models.RedeemComuResp, param models.Params, getRespChan chan models.RedeemResponse, ErrRespUseVouc chan error, wg *sync.WaitGroup) {
	fmt.Println("[ >>>>>>>>>>>>>>>>>> Use Voucher Comulative Voucher Otto AG <<<<<<<<<<<<<<<< ]")

	defer wg.Done()
	defer close(getRespChan)
	defer close(ErrRespUseVouc)

	// get CustID
	dataUser, errUser := db.CheckUser(param.AccountNumber)
	if errUser != nil {
		logrus.Info("User Not Eligible, Error : ", errUser)
	} else {
		logrus.Info("[ CampaignID ] :", req.CampaignID)
		logrus.Info("[ Coupon ID ] : ", redeemComu.CouponID)
		logrus.Info("[ Coupon Code ] : ", redeemComu.CouponCode)
		logrus.Info("[ Used ] : ", 1)
		logrus.Info("[ Customer ID ] : ", dataUser.CustID)
		logrus.Info("[ Category voucher ] : ", param.Category)
	}

	// Reedem Use Voucher
	param.Amount = redeemComu.Redeem.Amount
	param.RRN = redeemComu.Redeem.Rrn
	param.CouponID = redeemComu.CouponID
	param.PointTransferID = redeemComu.PointTransferID
	param.Comment = redeemComu.Comment
	resRedeem := services.RedeemUseVoucherComulative(req, param)
	getRespChan <- resRedeem
}
