package services

import (
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	opl "ottopoint-purchase/hosts/opl/host"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/opentracing/opentracing-go"
	"github.com/vjeantet/jodaTime"
	"go.uber.org/zap"
)

type UseVoucherVidioService struct {
	General models.GeneralModel
}

func (t UseVoucherVidioService) UseVoucherVidio(couponId string) models.Response {
	fmt.Println("[ Use Voucher Vidio Service ]")

	resp := models.Response{Meta: utils.ResponseMetaOK()}

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[ViewVoucher-Services]",
		zap.String("Coupon Id : ", couponId))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[ViewVoucher]")
	defer span.Finish()

	// get voucher
	getVouc, errGetVouc := db.GetUseVoucher(couponId)
	if errGetVouc != nil || getVouc.AccountNumber == "" {
		logs.Info("Internal server error")
		logs.Info("[Failed get data voucher from DB]")

		// sugarLogger.Info("Internal Server Error : ", errGet)
		sugarLogger.Info("[GetVoucher-Servcies]-[GetVoucher]")
		sugarLogger.Info("[Failed get data from DB]")
		resp = utils.GetMessageFailedErrorNew(resp, constants.RC_VOUCHER_NOTFOUND, constants.RD_VOUCHER_NOTFOUND)
		return resp
	}

	// get CustId Users
	dtaUsers, errUesr := db.GetUser(getVouc.AccountNumber)
	if errUesr != nil {
		fmt.Println("Internal Server Error")
		fmt.Println("Failed get User from db : ", errUesr)
		resp = utils.GetMessageFailedErrorNew(resp, 500, "Internal Server Error")
		return resp
	}

	// use voucher OPL
	dtaUseVouc, err2 := opl.CouponVoucherCustomer(getVouc.CampaignId, getVouc.CouponId, getVouc.ProductCode, dtaUsers.CustID, 1)
	if err2 != nil {
		fmt.Println("Internal Server Error")
		fmt.Println("Error access OPL use Voucher : ", err2)
		resp = utils.GetMessageFailedErrorNew(resp, 500, "Internal Server Error")
		return resp
	}

	// update transaction redeem into use
	timeUse := jodaTime.Format("dd-MM-YYYY HH:mm:ss", time.Now())
	_, errUpdate := db.UpdateVoucher(timeUse, getVouc.CouponId)
	if errUpdate != nil {
		logs.Info(fmt.Sprintf("[Error : %v]", errUpdate))
		logs.Info("[Gagal Update Voucher]")
		logs.Info("[UseVoucherUV]-[Package-Services]")
	}

	fmt.Println("Success Use Voucher ....")
	fmt.Println(dtaUseVouc)

	respVouch := models.RespUseVoucher{}
	respVouch.Code = dtaUseVouc.Coupons[0].Code
	respVouch.CouponID = dtaUseVouc.Coupons[0].CouponID
	respVouch.Used = dtaUseVouc.Coupons[0].Used
	respVouch.CampaignID = dtaUseVouc.Coupons[0].CampaignID
	respVouch.CustomerID = dtaUseVouc.Coupons[0].CustomerID

	resp = models.Response{
		Data: respVouch,
	}

	return resp

}
