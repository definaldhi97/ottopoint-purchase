package earnings

import (
	"errors"
	"fmt"
	"ottopoint-purchase/db"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type GetEarningRuleService struct {
	General models.GeneralModel
}

func (t GetEarningRuleService) NewGetEarningRuleService(productCode string) models.Response {
	res := models.Response{}

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[GetEarningRuleService]", zap.String("ProductCode : ", productCode))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[GetEarningRuleService]")
	defer span.Finish()

	fmt.Println("===== GetEarningRuleService =====")

	// Get EaringCode from DB
	data, err := db.GetEarningCodebyProductCode(productCode)
	if err != nil || data.Code == "" {

		fmt.Println(fmt.Sprintf("[Failed to Get EarningCode]-[Error : %v]", err))
		fmt.Println("[PackageEarnings]-[GetEarningCodebyProductCode]")

		sugarLogger.Info(fmt.Sprintf("[Failed to Get EarningCode]-[Error : %v]", err))
		sugarLogger.Info("[PackageEarnings]-[GetEarningCodebyProductCode]")

		res = utils.GetMessageResponse(res, 178, false, errors.New("Earning Rule not found"))

		return res
	}

	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: models.GetEarningRuleResp{
			Code: data.Code,
		},
	}

	return res

}
