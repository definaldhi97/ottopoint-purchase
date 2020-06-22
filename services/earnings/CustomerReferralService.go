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
	"strconv"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

func (t EarningPointServices) CustomerReferralService(req models.EarningReq, institutionID string) models.Response {
	res := models.Response{}

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[CustomerReferralService]",
		zap.String("Earning : ", req.Earning), zap.String("AccountNumber1 : ", req.AccountNumber1),
		zap.String("AccountNumber2 : ", req.AccountNumber2), zap.Int64("Amount : ", req.Amount),
		zap.String("ProductCode : ", req.ProductCode), zap.String("ProductName : ", req.ProductName),
		zap.String("ReferenceId : ", req.ReferenceId), zap.String("Remark : ", req.Remark))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[CustomerReferralService]")
	defer span.Finish()

	fmt.Println("===== CustomerReferralService =====")

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

	// Get CustID OPL from DB AccountNumber1
	dataUser1, errUser1 := db.CheckUser(req.AccountNumber1)
	if errUser1 != nil || dataUser1.CustID == "" {
		fmt.Println(fmt.Sprintf("[Internal Server Error 1 : %v]", errUser1))
		fmt.Println(fmt.Sprintf("[CustomerReferralService]-[Error 1 : %v]", dataUser1))
		fmt.Println("[Failed to Get CustID OPL]-[CheckUser]")

		sugarLogger.Info("[Internal Server Error 1]")
		sugarLogger.Info("[CustomerReferralService]")
		sugarLogger.Info("[Failed to Get CustID OPL]-[CheckUser]")

		res = utils.GetMessageFailedErrorNew(res, constants.RC_ERROR_ACC_NOT_ELIGIBLE, constants.RD_ERROR_ACC_NOT_ELIGIBLE)

		resData, _ := json.Marshal(&res)

		save.StatusCode = "72"
		save.StatusMessage = constants.RD_ERROR_ACC_NOT_ELIGIBLE
		save.ResponderData = string(resData)

		go db.SaveEarning(save)

		return res
	}

	// Get CustID OPL from DB AccountNumber2
	dataUser2, errUser2 := db.CheckUser(req.AccountNumber2)
	if errUser2 != nil || dataUser2.CustID == "" {
		fmt.Println(fmt.Sprintf("[Internal Server Error 2 : %v]", errUser2))
		fmt.Println(fmt.Sprintf("[CustomerReferralService]-[Error 2 : %v]", dataUser2))
		fmt.Println("[Failed to Get CustID OPL]-[CheckUser]")

		sugarLogger.Info("[Internal Server Error 2]")
		sugarLogger.Info("[CustomerReferralService]")
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
		fmt.Println(fmt.Sprintf("[Error : %v]", errEarning))
		fmt.Println(fmt.Sprintf("[CustomerReferralService]-[Error : %v]", earning))
		fmt.Println("[Failed to Get Data Earning]-[GetEarningCode]")

		sugarLogger.Info("[Internal Server Error]")
		sugarLogger.Info("[CustomerReferralService]")
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

	text1 := fmt.Sprintf("Kamu mendapatkan %v poin dari OttoPay, makin sering transaksi makin untung.", int64(earning.PointsAmount))
	totalpoint1 := strconv.Itoa(int(earning.PointsAmount))
	// Send Point
	senPoint1, errPoint1 := opl.TransferPoint(dataUser1.CustID, totalpoint1, text1)
	if errPoint1 != nil || senPoint1.PointsTransferId == "" {

		fmt.Println(fmt.Sprintf("[Internal Server Error 1 : %v]", errPoint1))
		fmt.Println(fmt.Sprintf("[CustomerReferralService]-[Error 1 : %v]", senPoint1))
		fmt.Println("[Failed to send Point]-[TransferPoint]")

		sugarLogger.Info("[Internal Server Error 1]")
		sugarLogger.Info("[CustomerReferralService]")
		sugarLogger.Info("[Failed to send Point]-[TransferPoint]")

		res = utils.GetMessageResponse(res, 80, false, errors.New("Gagal Transfer Point"))

		resData, _ := json.Marshal(&res)

		save.StatusCode = "80"
		save.StatusMessage = "Gagal Transfer Point"
		save.ResponderData = string(resData)

		go db.SaveEarning(save)

		return res
	}

	text2 := fmt.Sprintf("Kamu mendapatkan %v poin dari OttoPay, makin sering transaksi makin untung.", int64(earning.PointsAmount))
	totalpoint2 := strconv.Itoa(int(earning.PointsAmount))
	// Send Point
	senPoint2, errPoint2 := opl.TransferPoint(dataUser2.CustID, totalpoint2, text2)
	if errPoint2 != nil || senPoint2.PointsTransferId == "" {

		fmt.Println(fmt.Sprintf("[Internal Server Error 2 : %v]", errPoint2))
		fmt.Println(fmt.Sprintf("[CustomerReferralService]-[Error 2 : %v]", senPoint2))
		fmt.Println("[Failed to send Point]-[TransferPoint]")

		sugarLogger.Info("[Internal Server Error 2]")
		sugarLogger.Info("[CustomerReferralService]")
		sugarLogger.Info("[Failed to send Point]-[TransferPoint]")

		res = utils.GetMessageResponse(res, 80, false, errors.New("Gagal Transfer Point"))

		resData, _ := json.Marshal(&res)

		save.StatusCode = "80"
		save.StatusMessage = "Gagal Transfer Point"
		save.ResponderData = string(resData)

		go db.SaveEarning(save)

		return res
	}

	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: models.EarningResp{
			ReferenceId: req.ReferenceId,
			Point:       int64(earning.PointsAmount),
		},
	}

	resData, _ := json.Marshal(&res)

	save.StatusCode = "200"
	save.StatusMessage = "Success"
	save.ResponderData = string(resData)
	save.PointTransferId = senPoint1.PointsTransferId + " || " + senPoint2.PointsTransferId
	save.Point = int64(earning.PointsAmount)

	go db.SaveEarning(save)

	return res

}
