package v2_migrate

import (
	"fmt"
	"ottopoint-purchase/db"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/vjeantet/jodaTime"
	"go.uber.org/zap"
)

type CallabckUVServices struct {
	General models.GeneralModel
}

func (t CallabckUVServices) CallbackVoucherUV(req models.UseVoucherUVReq, param models.Params, campaignID string) models.Response {
	var res models.Response

	fmt.Println("[ >>>>>>>>>>>>>>>>>>>>> Callbakc Voucher UV Service <<<<<<<<<<<<<<<<<< ]")

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[UseVoucherUV]",
		zap.String("AccountNumber : ", param.AccountNumber), zap.String("InstitutionID : ", param.InstitutionID),
		zap.String("category : ", param.Category), zap.String("campaignId : ", campaignID),
		zap.String("AccountID : ", param.AccountId), zap.String("AccountNumber : ", param.AccountNumber),
		zap.String("VoucherCode : ", req.VoucherCode))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[UseVoucherUV]")
	defer span.Finish()

	logrus.Info("Campaign : ", campaignID)
	logrus.Info("CouponID : ", param.CouponID)
	logrus.Info("ProductCode : ", param.CouponCode)
	logrus.Info("AccountID : ", param.AccountId)

	// // Use Voucher to Openloyalty
	// use, err2 := opl.CouponVoucherCustomer(campaignID, param.CouponID, param.CouponCode, param.AccountId, 1)

	// var useErr string
	// for _, value := range use.Coupons {
	// 	useErr = value.CouponID
	// }

	// if err2 != nil || useErr == "" {

	// 	logs.Info(fmt.Sprintf("[Error : %v]", err2))
	// 	logs.Info(fmt.Sprintf("[Response : %v]", use))
	// 	logs.Info("[Error from OPL]-[CouponVoucherCustomer]")

	// 	// go SaveTransactionUV(param, useUV, reqUV, req, "Used", "01", "")

	// 	res = utils.GetMessageResponse(res, 400, false, errors.New("Gagal Use Voucher, Harap coba lagi"))
	// 	return res
	// }

	timeUse := jodaTime.Format("dd-MM-YYYY HH:mm:ss", time.Now())

	_, errUpdate := db.UpdateVoucher(timeUse, param.CouponID)
	if errUpdate != nil {

		logs.Info(fmt.Sprintf("[Error : %v]", errUpdate))
		logs.Info("[Gagal Update Voucher]")
		logs.Info("[UseVoucherUV]-[Package-Services]")

	}

	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: models.UseVoucherUVResp{
			Voucher: param.NamaVoucher,
		},
	}
	return res

}
