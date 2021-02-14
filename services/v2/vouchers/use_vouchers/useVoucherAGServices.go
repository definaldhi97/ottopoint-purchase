package use_vouchers

import (
	"errors"
	"fmt"
	"ottopoint-purchase/db"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"

	"github.com/sirupsen/logrus"
)

// func UseVoucherAggregator(req models.UseVoucherReq, param models.Params) models.Response {
func UseVoucherAggregatorService(req models.UseVoucherReq, param models.Params) models.Response {
	// fmt.Println("[ >>>>>>>>>>>>>>>>>> V2 Migrate Use Vouhcer Aggregator Migrate Services <<<<<<<<<<<<<<<< ]")

	var res models.Response

	nameservice := "[PackageUserVouchers]-[UseVoucherAggregatorService]"
	logReq := fmt.Sprintf("[CouponID : %v]", req.CouponID)

	logrus.Info(nameservice)

	get, err := db.GetVoucherAg(param.AccountId, req.CouponID)
	if err != nil {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[GetVoucherAg]-[Error : %v]", err))
		logrus.Println(logReq)

		res = utils.GetMessageResponse(res, 422, false, errors.New("Voucher Tidak Ditemukan"))
		return res
	}

	spend, err := db.GetVoucherSpending(get.AccountId, get.CouponID)
	if err != nil {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[GetVoucherSpending]-[Error : %v]", err))
		logrus.Println(logReq)

		res = utils.GetMessageResponse(res, 422, false, errors.New("Terjadi Kesalahan"))
		return res
	}

	// // Update Status Voucher
	// // timeUse := jodaTime.Format("dd-MM-YYYY HH:mm:ss", time.Now())
	// timeUse := time.Now()
	// go db.UpdateVoucher(timeUse, spend.CouponId)

	// codeVoucher := decryptVoucherCode(spend.VoucherCode, spend.CouponId)

	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: models.GetVoucherAgResp{
			Voucher:     param.NamaVoucher,
			VoucherCode: spend.VoucherCode,
			Link:        spend.VoucherLink,
		},
	}
	return res

}
