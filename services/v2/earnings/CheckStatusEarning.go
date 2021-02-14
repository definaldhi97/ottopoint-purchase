package earnings

import (
	"errors"
	"fmt"
	"ottopoint-purchase/db"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"

	"github.com/sirupsen/logrus"
)

type EarningsServices struct{}

func (t EarningsServices) CheckStatusEarningServices(referenceId, institution string) models.Response {
	res := models.Response{}

	nameservice := "[PackageEarnings]-[CheckStatusEarningServices]"

	logrus.Info(nameservice)

	// Get EaringCode from DB
	earning, errEarning := db.GetCheckStatusEarning(referenceId, institution)
	if errEarning != nil || earning.ReferenceId == "" {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[GetCheckStatusEarning]-[Error : %v]", errEarning))
		logrus.Println("ReferenceId : ", referenceId, "Institution : ", institution)

		res = utils.GetMessageResponse(res, 178, false, errors.New("Earning Rule not found"))

		return res
	}

	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: models.EarningResp{
			ReferenceId: referenceId,
			Point:       earning.Point,
		},
	}

	return res

}
