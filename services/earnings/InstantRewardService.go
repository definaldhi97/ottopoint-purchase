package earnings

import (
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	opl "ottopoint-purchase/hosts/opl/host"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

func (t EarningPointServices) InstantRewardService(req models.EarningReq) models.Response {
	res := models.Response{}

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[InstantRewardService]",
		zap.String("Earning : ", req.Earning), zap.String("AccountNumber1 : ", req.AccountNumber1),
		zap.String("AccountNumber2 : ", req.AccountNumber2), zap.Int64("Amount : ", req.Amount),
		zap.String("ProductCode : ", req.ProductCode), zap.String("ProductName : ", req.ProductName),
		zap.String("ReferenceId : ", req.ReferenceId), zap.String("Remark : ", req.Remark))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[InstantRewardService]")
	defer span.Finish()

	fmt.Println("===== InstantRewardService =====")

	// Get EaringCode from DB
	earning, errEarning := db.GetEarningCode(req.Earning)
	if errEarning != nil || earning.Code == "" {
		fmt.Println(fmt.Sprintf("[Internal Server Error : %v]", errEarning))
		fmt.Println(fmt.Sprintf("[InstantRewardService]-[Error : %v]", earning))
		fmt.Println("[Failed to Get Data Earning]-[GetEarningCode]")

		sugarLogger.Info("[Internal Server Error]")
		sugarLogger.Info("[InstantRewardService]")
		sugarLogger.Info("[Failed to Get Data Earning]-[GetEarningCode]")

		// response belum ada
		return res
	}

	validateActive, errValidate := utils.ValidateTimeActive(earning.Active, earning.AllTimeActive, earning.StartAt, earning.EndAt)

	if errValidate == false {

		// response belum ada
		res = utils.GetMessageFailedErrorNew(res, constants.RC_ERROR_ACC_NOT_ELIGIBLE, validateActive)
		return res
	}

	// Get CustID OPL from DB
	dataUser, errUser := db.CheckUser(req.AccountNumber1)
	if errUser != nil || dataUser.CustID == "" {
		fmt.Println(fmt.Sprintf("[Internal Server Error : %v]", errUser))
		fmt.Println(fmt.Sprintf("[InstantRewardService]-[Error : %v]", dataUser))
		fmt.Println("[Failed to Get CustID OPL]-[CheckUser]")

		sugarLogger.Info("[Internal Server Error]")
		sugarLogger.Info("[InstantRewardService]")
		sugarLogger.Info("[Failed to Get CustID OPL]-[CheckUser]")

		res = utils.GetMessageFailedErrorNew(res, constants.RC_ERROR_ACC_NOT_ELIGIBLE, constants.RD_ERROR_ACC_NOT_ELIGIBLE)
		return res
	}

	// send Voucher
	sendVoucher, errVoucher := opl.RedeemVoucherCumulative(earning.RewardCampaignID, dataUser.CustID, "1", "1")

	var voucherCode string
	for _, value := range sendVoucher.Coupons {
		voucherCode = value.Code
	}

	if errVoucher != nil || voucherCode == "" {
		fmt.Println(fmt.Sprintf("[Internal Server Error : %v]", errEarning))
		fmt.Println(fmt.Sprintf("[InstantRewardService]-[Error : %v]", sendVoucher))
		fmt.Println("[Failed to Send Voucher]-[GetEarningCode]")

		sugarLogger.Info("[Internal Server Error]")
		sugarLogger.Info("[InstantRewardService]")
		sugarLogger.Info("[Failed to Send Voucher]-[GetEarningCode]")

		// response belum ada
		return res
	}

	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: models.EarningResp{
			ReferenceId: req.ReferenceId,
		},
	}

	return res

}
