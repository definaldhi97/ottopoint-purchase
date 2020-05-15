package services

import (
	redismodels "ottopoint-purchase/hosts/redis_token/models"
	"ottopoint-purchase/models"

	// "github.com/vjeantet/jodaTime"
	"github.com/opentracing/opentracing-go"
)

type EarningRuleServices struct {
	General models.GeneralModel
}

func (t EarningRuleServices) NewEarningRuleServices(req models.EarningRuleReq, dataToken redismodels.TokenResp, header models.RequestHeader) models.Response {
	res := models.Response{}

	// sugarLogger := t.General.OttoZaplog
	// sugarLogger.Info("[NewEarningRuleServices]",
	// 	zap.String("Code : ", req.Code), zap.String("institution", header.InstitutionID))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[NewEarningRuleServices]")
	defer span.Finish()

	// getRule, _ := db.CheckRule(req.Code)

	// if getRule.AllTimeActive == true {
	// 	if getRule.StartAt
	// }

	// jodaTime.Format("dd-MM-YYYY", getRule.),

	return res

}
