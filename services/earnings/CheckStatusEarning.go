package earnings

import (
	"encoding/json"
	"errors"
	"fmt"
	"ottopoint-purchase/db"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"

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

	// Get EaringCode from DB
	earning, errEarning := db.GetCheckStatusEarning(referenceId, institution)
	if errEarning != nil || earning.ReferenceId == "" {
		fmt.Println(fmt.Sprintf("[Internal Server Error : %v]", errEarning))
		fmt.Println(fmt.Sprintf("[GetCheckStatusEarning]-[Error : %v]", earning))
		fmt.Println("[Failed to Get Data Earning]-[GetCheckStatusEarning]")

		sugarLogger.Info("[Internal Server Error]")
		sugarLogger.Info("[GetCheckStatusEarning]")
		sugarLogger.Info("[Failed to Get Data Earning]-[GetCheckStatusEarning]")

		res = utils.GetMessageResponse(res, 178, false, errors.New("Earning Rule not found"))

		return res
	}

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
