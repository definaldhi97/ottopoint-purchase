package earnings

import (
	"encoding/json"
	"errors"
	"fmt"
	"ottopoint-purchase/db"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"
	"strconv"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type CheckStatusEarningService struct {
	General models.GeneralModel
}

func (t CheckStatusEarningService) CheckStatusEarningServices(referenceId, institution string) models.Response {
	res := models.Response{}

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[CheckStatusEarningServices]",
		zap.String("ReferenceId : ", referenceId), zap.String("Institution : ", institution))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[CheckStatusEarningServices]")
	defer span.Finish()

	fmt.Println("===== CheckStatusEarningServices =====")

	var statusMessage string
	var statusCode int

	// Get EaringCode from DB
	earning, errEarning := db.GetCheckStatusEarning(referenceId, institution)

	code, _ := strconv.Atoi(earning.StatusCode)

	statusMessage = earning.StatusMessage
	statusCode = code

	if errEarning != nil || earning.ReferenceId == "" {

		fmt.Println(fmt.Sprintf("[Error : %v]", errEarning))
		fmt.Println(fmt.Sprintf("[GetCheckStatusEarning]-[Error : %v]", earning))
		fmt.Println("[Failed to Get Data Earning]-[GetCheckStatusEarning]")

		sugarLogger.Info("[Internal Server Error]")
		sugarLogger.Info("[GetCheckStatusEarning]")
		sugarLogger.Info("[Failed to Get Data Earning]-[GetCheckStatusEarning]")

		earningLog, errLog := db.GetErrorLogStatusEarning(referenceId, institution)
		if errLog != nil || earningLog.ReferenceId == "" {

			fmt.Println(fmt.Sprintf("[Error : %v]", errLog))
			fmt.Println(fmt.Sprintf("[GetErrorLogStatusEarning]-[Error : %v]", earningLog))
			fmt.Println("[Failed to Get Data Earning]-[GetErrorLogStatusEarning]")

			sugarLogger.Info("[Internal Server Error]")
			sugarLogger.Info("[GetErrorLogStatusEarning]")
			sugarLogger.Info("[Failed to Get Data Earning]-[GetErrorLogStatusEarning]")

			res = utils.GetMessageResponse(res, 153, false, errors.New("Data Not Found"))

			return res
		}

		code, _ := strconv.Atoi(earningLog.StatusCode)

		statusMessage = earningLog.StatusMessage
		statusCode = code
	}

	if earning.StatusCode == "200" {
		resEarning := models.ResponseEarning{}

		resp := []byte(earning.ResponderData)
		errRespEarning := json.Unmarshal(resp, &resEarning)
		fmt.Println("Error Unmarshal : ", errRespEarning)

		res = models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: resEarning,
		}

		return res
	}

	res = models.Response{
		Meta: models.MetaData{
			Code:    statusCode,
			Status:  false,
			Message: statusMessage,
		},
	}

	return res

}
