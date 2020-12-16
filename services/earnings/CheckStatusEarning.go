package earnings

import (
	"errors"
	"fmt"
	"ottopoint-purchase/db"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"
)

type CheckStatusEarningService struct{}

func (t CheckStatusEarningService) CheckStatusEarningServices(referenceId, institution string) models.Response {
	res := models.Response{}

	fmt.Println("===== CheckStatusEarningServices =====")

	// getId, errId := db.GetIdInstitution(institution)
	// if errId != nil || getId.Name == "" {
	// 	fmt.Println(fmt.Sprintf("[Internal Server Error : %v]", errId))
	// 	fmt.Println(fmt.Sprintf("[GetIdInstitution]-[Error : %v]", getId))
	// 	fmt.Println("[Failed to Get Data Earning]-[GetIdInstitution]")

	// 	sugarLogger.Info("[Internal Server Error]")
	// 	sugarLogger.Info("[GetIdInstitution]")
	// 	sugarLogger.Info("[Failed to Get Institution]-[GetIdInstitution]")

	// 	res = utils.GetMessageResponse(res, 82, false, errors.New("Invalid InstitutionID"))

	// 	return res
	// }

	// Get EaringCode from DB
	earning, errEarning := db.GetCheckStatusEarning(referenceId, institution)
	if errEarning != nil || earning.ReferenceId == "" {

		fmt.Println(fmt.Sprintf("[Internal Server Error : %v]", errEarning))
		fmt.Println(fmt.Sprintf("[GetCheckStatusEarning]-[Error : %v]", earning))
		fmt.Println("[Failed to Get Data Earning]-[GetCheckStatusEarning]")

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
