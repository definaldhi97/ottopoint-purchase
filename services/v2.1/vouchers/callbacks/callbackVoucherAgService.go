package callbacks

import (
	"errors"
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	"ottopoint-purchase/models"
	callback "ottopoint-purchase/models/v21/callback"
	"ottopoint-purchase/utils"
	"strings"

	"github.com/sirupsen/logrus"
)

func CallbackVoucherAG_V21_Service(req callback.CallbackVoucherAGReq) models.Response {
	var res models.Response

	nameservice := "[PackageCallbacks]-[CallbackVoucherAG_V21_Service]"

	logReq := fmt.Sprintf("[TransactionID : %v]", req.TransactionId)

	logrus.Info(nameservice)

	// validate TrxID
	errTrxId := db.CheckTrxbyTrxID(req.OrderId)
	logrus.Println(">>> CheckTrxbyTrxID <<<")
	if errTrxId != nil {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[CheckTrxId]-[Error : %v]", errTrxId))
		logrus.Println(logReq)

		res = utils.GetMessageResponse(res, 422, false, errors.New("TrxId Tidak Ditemukan"))

		return res

	}

	logrus.Println(">>> Lanjut <<<")

	reUpdate := db.VoucherTypeDB{}

	// PPOB (1)
	if strings.ToLower(req.VoucherType) == strings.ToLower(constants.VoucherTypePPOB) {

		logrus.Println(">>> PPOB <<<")

		reUpdate = db.VoucherTypeDB{
			VoucherType:  1,
			OrderId:      req.TransactionId,
			ResponseCode: req.Data.ResponseCode,
			ResponseDesc: req.Data.ResponseDesc,
		}

		update := db.UpdateVoucherbyVoucherType(reUpdate, req.OrderId)
		logrus.Info("Response Update : ", update)
	} else {
		res = utils.GetMessageResponse(res, 500, false, errors.New("Internal Server Error"))

		return res

	}

	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: nil,
	}

	return res
}
