package earnings

import (
	"errors"
	"fmt"
	"ottopoint-purchase/db"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"

	"github.com/sirupsen/logrus"
)

func (t EarningsServices) NewGetEarningRuleService(productCode string) models.Response {
	res := models.Response{}

	nameservice := "[PackageEarnings]-[NewGetEarningRuleService]"

	logrus.Info(nameservice)

	// Get EaringCode from DB
	data, err := db.GetEarningCodebyProductCode(productCode)
	if err != nil || data.Code == "" {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[GetEarningCodebyProductCode]-[Error : %v]", err))
		logrus.Println("Request : ", productCode)

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
