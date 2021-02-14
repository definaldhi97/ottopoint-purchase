package callbacks

import (
	"fmt"
	"time"

	"ottopoint-purchase/db"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"

	"github.com/sirupsen/logrus"
)

func CallbackVoucherUV(req models.UseVoucherUVReq, param models.Params, campaignID string) models.Response {
	var res models.Response

	nameservice := "[PackageCallbacks]-[CallbackVoucherUV]"
	logReq := fmt.Sprintf("[AccountId : %v, VoucherCode : %v]", req.AccountId, req.VoucherCode)

	logrus.Info(nameservice)

	// timeUse := jodaTime.Format("dd-MM-YYYY HH:mm:ss", time.Now())
	timeUse := time.Now()

	_, errUpdate := db.UpdateVoucher(timeUse, param.CouponID)
	if errUpdate != nil {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[UpdateVoucher]-[Error : %v]", errUpdate))
		logrus.Println(logReq)

	}

	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: models.UseVoucherUVResp{
			Voucher: param.NamaVoucher,
		},
	}
	return res

}
