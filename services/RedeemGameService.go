package services

import (
	ottomart "ottopoint-purchase/hosts/ottomart/models"
	"ottopoint-purchase/models"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type RedeemGameService struct {
	General models.GeneralModel
}

func (t RedeemGameService) RedeemGame(req models.RedeemVoucherRequest, dataToken ottomart.ResponseToken) models.Response {
	var res models.Response

	resMeta := models.MetaData{
		Code:    200,
		Status:  true,
		Message: "Succesful",
	}

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[RedeemGameService-RedeemGame]", zap.String("Category", req.CampaignID), zap.String("CampaignID", req.CampaignID), zap.String("CustID", req.CustID), zap.String("ProductCode", req.ProductCode))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[RedeemGameService-RedeemGame]")
	defer span.Finish()

	res = models.Response{
		Meta: resMeta,
	}
	return res
}
