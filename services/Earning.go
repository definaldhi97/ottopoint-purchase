package services

import (
	"errors"
	"ottopoint-purchase/constants"
	opl "ottopoint-purchase/hosts/opl/host"
	redismodels "ottopoint-purchase/hosts/redis_token/models"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"

	"github.com/astaxie/beego/logs"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type EarningServices struct {
	General models.GeneralModel
}

func (t EarningServices) EarningPoint(req models.RulePointReq, dataToken redismodels.TokenResp, header models.RequestHeader) models.Response {
	var res models.Response

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[EarningPoint-Services]",
		zap.String("rule", req.EventName), zap.Int("amount", req.Amount), zap.String("rc", req.RC), zap.String("institution", header.InstitutionID))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[RedeemVoucher]")
	defer span.Finish()

	switch header.InstitutionID {
	case constants.OTTOPAY:
		res = OttopayEarning(req, dataToken)
	}

	return res
}

func OttopayEarning(req models.RulePointReq, dataToken redismodels.TokenResp) models.Response {
	res := models.Response{}

	dataPoint, errPoint := opl.RulePoint(req.EventName, dataToken.Data)
	if errPoint != nil {
		// sugarLogger.Info("[Error-dataPoint :", errPoint)
		logs.Info("[Error-dataPoint :", errPoint)
		res = utils.GetMessageResponse(res, 422, false, errors.New("Gagal Dapat Point"))
	}

	res = models.Response{
		Data: dataPoint.Point,
		Meta: utils.ResponseMetaOK(),
	}

	return res
}
