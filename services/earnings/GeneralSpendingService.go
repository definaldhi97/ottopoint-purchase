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

type EarningPointServices struct {
	General models.GeneralModel
}

func (t EarningPointServices) GeneralSpendingService(req models.EarningReq, institutionID string) models.Response {
	res := models.Response{}

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[GeneralSpendingService]",
		zap.String("Earning : ", req.Earning), zap.String("AccountNumber1 : ", req.AccountNumber1),
		zap.String("AccountNumber2 : ", req.AccountNumber2), zap.Int64("Amount : ", req.Amount),
		zap.String("ProductCode : ", req.ProductCode), zap.String("ProductName : ", req.ProductName),
		zap.String("ReferenceId : ", req.ReferenceId), zap.String("Remark : ", req.Remark))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[GeneralSpendingService]")
	defer span.Finish()

	fmt.Println("===== GeneralSpendingService =====")

	reqData, _ := json.Marshal(&req)

	save := dbmodels.TEarning{
		// ID             : ,
		EarningRule:    req.Earning,
		PartnerId: institutionID,
		ReferenceId:    req.ReferenceId,
		TransactionId:  utils.GenTransactionId(),
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
		fmt.Println(fmt.Sprintf("[GeneralSpendingService]-[Error : %v]", dataUser))
		fmt.Println("[Failed to Get CustID OPL]-[CheckUser]")

		sugarLogger.Info("[Internal Server Error]")
		sugarLogger.Info("[GeneralSpendingService]")
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
		fmt.Println(fmt.Sprintf("[GeneralSpendingService]-[Error : %v]", earning))
		fmt.Println("[Failed to Get Data Earning]-[GetEarningCode]")

		sugarLogger.Info("[Internal Server Error]")
		sugarLogger.Info("[GeneralSpendingService]")
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

	var listsku []string
	// Unmarshal string to []string
	excludedSku := []byte(earning.ExcludedSkus)
	errSku := json.Unmarshal(excludedSku, &listsku)
	fmt.Println("Error Unmarshal : ", errSku)

	// Validate Excluded SKU
	for _, value := range listsku {
		if req.ProductCode == value {
			fmt.Println("=== Product Excluded ===")

			res = utils.GetMessageResponse(res, 179, false, errors.New("Earning Failed, Excluded SKU"))

			resData, _ := json.Marshal(&res)

			save.StatusCode = "179"
			save.StatusMessage = "Earning Failed, Excluded SKU"
			save.ResponderData = string(resData)

			go db.SaveEarning(save)

			return res
		}
	}

	// validate Min Amount
	if req.Amount < int64(earning.MinOrderValue) {
		fmt.Println("=== Amount belum cukup ===")

		res = utils.GetMessageResponse(res, 180, false, errors.New("Earning Failed, Min Amount"))

		resData, _ := json.Marshal(&res)

		save.StatusCode = "180"
		save.StatusMessage = "Earning Failed, Min Amount"
		save.ResponderData = string(resData)

		go db.SaveEarning(save)

		return res

	}

	var listskuId []string
	var multiply, point int64
	// Unmarshal string to []string
	skuId := []byte(earning.SkuIds)
	errskuId := json.Unmarshal(skuId, &listskuId)
	fmt.Println("Error Unmarshal : ", errskuId)

	// validate MultiplyPoint
	for _, val := range listskuId {
		if req.ProductCode == val {
			multiply = int64(earning.Multiplier)
		}
	}

	point = int64(earning.PointValue)
	if multiply != 0 {
		point = point * multiply
	}

	text := fmt.Sprintf("Kamu mendapatkan %v poin dari OttoPay, makin sering transaksi makin untung.", int64(point))
	totalpoint := strconv.Itoa(int(point))
	// Send Point
	senPoint, errPoint := opl.TransferPoint(dataUser.CustID, totalpoint, text)
	if errPoint != nil || senPoint.PointsTransferId == "" {

		fmt.Println(fmt.Sprintf("[Internal Server Error : %v]", errPoint))
		fmt.Println(fmt.Sprintf("[GeneralSpendingService]-[Error : %v]", senPoint))
		fmt.Println("[Failed to send Point]-[TransferPoint]")

		sugarLogger.Info("[Internal Server Error]")
		sugarLogger.Info("[GeneralSpendingService]")
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
			Point:       point,
		},
	}

	resData, _ := json.Marshal(&res)

	save.StatusCode = "200"
	save.StatusMessage = "Success"
	save.PointTransferId = senPoint.PointsTransferId
	save.Point = point
	save.ResponderData = string(resData)

	go db.SaveEarning(save)

	return res

}
