package earnings

import (
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	opl "ottopoint-purchase/hosts/opl/host"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"
	"strconv"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

func (t EarningPointServices) CustomeEventRuleService(req models.EarningReq) models.Response {
	res := models.Response{}

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[CustomeEventRuleService]",
		zap.String("Earning : ", req.Earning), zap.String("AccountNumber1 : ", req.AccountNumber1),
		zap.String("AccountNumber2 : ", req.AccountNumber2), zap.Int64("Amount : ", req.Amount),
		zap.String("ProductCode : ", req.ProductCode), zap.String("ProductName : ", req.ProductName),
		zap.String("ReferenceId : ", req.ReferenceId), zap.String("Remark : ", req.Remark))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[CustomeEventRuleService]")
	defer span.Finish()

	fmt.Println("===== CustomeEventRuleService =====")

	// Get CustID OPL from DB
	dataUser, errUser := db.CheckUser(req.AccountNumber1)
	if errUser != nil || dataUser.CustID == "" {
		fmt.Println(fmt.Sprintf("[Internal Server Error : %v]", errUser))
		fmt.Println(fmt.Sprintf("[CustomeEventRuleService]-[Error : %v]", dataUser))
		fmt.Println("[Failed to Get CustID OPL]-[CheckUser]")

		sugarLogger.Info("[Internal Server Error]")
		sugarLogger.Info("[CustomeEventRuleService]")
		sugarLogger.Info("[Failed to Get CustID OPL]-[CheckUser]")

		res = utils.GetMessageFailedErrorNew(res, constants.RC_ERROR_ACC_NOT_ELIGIBLE, constants.RD_ERROR_ACC_NOT_ELIGIBLE)
		return res
	}

	// Get EaringCode from DB
	earning, errEarning := db.GetEarningCode(req.Earning)
	if errEarning != nil || earning.Code == "" {
		fmt.Println(fmt.Sprintf("[Internal Server Error : %v]", errEarning))
		fmt.Println(fmt.Sprintf("[CustomeEventRuleService]-[Error : %v]", earning))
		fmt.Println("[Failed to Get Data Earning]-[GetEarningCode]")

		sugarLogger.Info("[Internal Server Error]")
		sugarLogger.Info("[CustomeEventRuleService]")
		sugarLogger.Info("[Failed to Get Data Earning]-[GetEarningCode]")

		// response belum ada
		return res
	}

	text := fmt.Sprintf("Kamu mendapatkan %v poin dari OttoPay, makin sering transaksi makin untung.", int64(earning.PointsAmount))
	totalpoint := strconv.Itoa(int(earning.PointsAmount))
	// Send Point
	senPoint, errPoint := opl.TransferPoint(dataUser.CustID, totalpoint, text)
	if errPoint != nil || senPoint.PointsTransferId == "" {

		fmt.Println(fmt.Sprintf("[Internal Server Error : %v]", errPoint))
		fmt.Println(fmt.Sprintf("[CustomeEventRuleService]-[Error : %v]", senPoint))
		fmt.Println("[Failed to send Point]-[TransferPoint]")

		sugarLogger.Info("[Internal Server Error]")
		sugarLogger.Info("[CustomeEventRuleService]")
		sugarLogger.Info("[Failed to send Point]-[TransferPoint]")

		// response belum ada
		return res
	}

	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: models.EarningResp{
			ReferenceId: req.ReferenceId,
			Point:       int64(earning.PointsAmount),
		},
	}

	return res

}
