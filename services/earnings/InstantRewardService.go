package earnings

import (
	"encoding/json"
	"errors"
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	opl "ottopoint-purchase/hosts/opl/host"
	"ottopoint-purchase/models"
	"ottopoint-purchase/models/dbmodels"
	"ottopoint-purchase/utils"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

func (t EarningPointServices) InstantRewardService(req models.EarningReq, institutionID string) models.Response {
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

	reqData, _ := json.Marshal(&req)

	save := dbmodels.TEarning{
		// ID             : ,
		EarningRule:    req.Earning,
		MInstitutionId: institutionID,
		ReferenceId:    req.ReferenceId,
		Transactionid:  utils.GenTransactionId(),
		ProductCode:    req.ProductCode,
		ProductName:    req.ProductName,
		AccountNumber1: req.AccountNumber1,
		AccountNumber2: req.AccountNumber2,
		Amount:         req.Amount,
		// Point:          int64(earning.PointsAmount),
		Remark: req.Remark,
		// StatusCode     : ,
		// StatusMessage  : ,
		// PointTransferId: senPoint.PointsTransferId,
		RequestorData: string(reqData),
		// ResponderData  : ,
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

		resData, _ := json.Marshal(&res)

		save.StatusCode = "72"
		save.StatusMessage = constants.RD_ERROR_ACC_NOT_ELIGIBLE
		save.ResponderData = string(resData)

		go db.SaveEarning(save)

		return res
	}

	// Get EaringCode from DB
	earning, errEarning := db.GetEarningCode(req.Earning)
	if errEarning != nil || earning.Code == "" {
		fmt.Println(fmt.Sprintf("[Internal Server Error : %v]", errEarning))
		fmt.Println(fmt.Sprintf("[InstantRewardService]-[Error : %v]", earning))
		fmt.Println("[Failed to Get Data Earning]-[GetEarningCode]")

		sugarLogger.Info("[Internal Server Error]")
		sugarLogger.Info("[InstantRewardService]")
		sugarLogger.Info("[Failed to Get Data Earning]-[GetEarningCode]")

		res = utils.GetMessageResponse(res, 178, false, errors.New("Earning Rule not found"))

		resData, _ := json.Marshal(&res)

		save.StatusCode = "178"
		save.StatusMessage = "Earning Rule not found"
		save.ResponderData = string(resData)

		go db.SaveEarning(save)

		return res
	}

	errValidate := utils.ValidateTimeActive(earning.Active, earning.AllTimeActive, earning.StartAt, earning.EndAt)

	if errValidate == false {

		res = utils.GetMessageResponse(res, 183, false, errors.New("Earning rule is not active"))

		resData, _ := json.Marshal(&res)

		save.StatusCode = "183"
		save.StatusMessage = "Earning rule is not active"
		save.ResponderData = string(resData)

		go db.SaveEarning(save)

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

		res = utils.GetMessageResponse(res, 182, false, errors.New("Failed give Voucher"))

		resData, _ := json.Marshal(&res)

		save.StatusCode = "182"
		save.StatusMessage = "Failed give Voucher"
		save.ResponderData = string(resData)

		go db.SaveEarning(save)

		return res
	}

	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: models.EarningResp{
			ReferenceId: req.ReferenceId,
		},
	}

	resData, _ := json.Marshal(&res)

	save.StatusCode = "200"
	save.StatusMessage = "Success"
	save.ResponderData = string(resData)
	// save.PointTransferId = senPoint.PointsTransferId
	// save.Point = int64(earning.PointsAmount)

	go db.SaveEarning(save)

	return res

}
