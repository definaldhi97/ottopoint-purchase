package callbacks

import (
	"fmt"
	"ottopoint-purchase/db"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"

	"github.com/sirupsen/logrus"
)

func CallbackVoucherAgg(req models.CallbackRequestVoucherAg) models.Response {

	var res models.Response

	nameservice := "[PackageCallbacks]-[CallbackVoucherAgg]"

	logReq := fmt.Sprintf("[TransactionID : %v, VoucherID : %v]", req.TransactionID, req.Data.VoucherID)

	logrus.Info(nameservice)

	// Get TSpending
	tspending, errSpending := db.GetVoucherAgSpending(req.Data.VoucherCode, req.TransactionID)
	if errSpending != nil {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[GetVoucherAgSpending]-[Error : %v]", errSpending))
		logrus.Println(logReq)

		res = utils.GetMessageResponse(res, 422, false, errSpending)

		return res
	}

	// Update TSpending
	_, err := db.UpdateVoucherAg(req.Data.RedeemedDate, req.Data.UsedDate, tspending.ID)
	if err != nil {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[GetVoucherAgSpending]-[Error : %v]", err))
		logrus.Println(logReq)

	}

	go db.UpdateTSchedulerVoucherAG(req.Data.OrderID)

	res = models.Response{
		Meta: utils.ResponseMetaOK(),
	}

	return res
}
