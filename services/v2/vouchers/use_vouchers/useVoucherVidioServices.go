package use_vouchers

import (
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"
	"time"

	"github.com/sirupsen/logrus"
)

// func  UseVoucherVidio(couponId string) models.Response {
func UseVoucherVidioServices(couponId string) models.Response {
	// fmt.Println("[ >>>>>>>>>>>>>>>>>> V2 Migrate Use Vouhcer Vidio Migrate Services <<<<<<<<<<<<<<<< ]")

	resp := models.Response{}

	nameservice := "[PackageUserVouchers]-[UseVoucherVidioServices]"
	logReq := fmt.Sprintf("[CouponID : %v]", couponId)

	// get voucher
	getVouc, errGetVouc := db.GetUseVoucher(couponId)
	if errGetVouc != nil || getVouc.AccountNumber == "" {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[GetUseVoucher]-[Error : %v]", errGetVouc))
		logrus.Println(logReq)

		resp = utils.GetMessageFailedErrorNew(resp, constants.RC_VOUCHER_NOTFOUND, constants.RD_VOUCHER_NOTFOUND)
		return resp
	}

	// update transaction redeem into use
	// timeUse := jodaTime.Format("dd-MM-YYYY HH:mm:ss", time.Now())
	timeUse := time.Now()
	_, errUpdate := db.UpdateVoucher(timeUse, getVouc.CouponId)
	if errUpdate != nil {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[UpdateVoucher]-[Error : %v]", errUpdate))
		logrus.Println(logReq)

	}

	respVouch := models.RespUseVoucher{}
	respVouch.Code = getVouc.ProductCode
	respVouch.CouponID = getVouc.CouponId
	respVouch.Used = true
	respVouch.CampaignID = *getVouc.MRewardID
	respVouch.CustomerID = getVouc.AccountId

	resp = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: respVouch,
	}

	return resp
}
