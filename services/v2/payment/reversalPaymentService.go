package payment

import (
	"errors"
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	"ottopoint-purchase/models"
	sp "ottopoint-purchase/models/v2/payment"
	"ottopoint-purchase/utils"

	"github.com/sirupsen/logrus"
)

func ReversalPaymentService(req sp.ReversalPaymentReq, param models.Params, header models.RequestHeader) models.Response {
	res := models.Response{}

	nameservice := "[PackagePayment]-[ReversalPaymentService]"
	logReq := fmt.Sprintf("[ReferenceId : %v]", req.ReferenceId)

	logrus.Info(nameservice)

	// validate reffId
	check, errCheck := db.CheckReffIdSplitBillReversal(req.ReferenceId)
	if errCheck == nil || check.ReferenceId != "" {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[CheckReffIdSplitBillReversal]-[Reference ID is Found : %v]", check.ReferenceId))
		logrus.Println(logReq)

		res = utils.GetMessageResponse(res, 172, false, errors.New("Duplicate Reference ID"))

		return res
	}

	getData, errData := db.CheckReffIdSplitBill(req.ReferenceId)
	if errData != nil || getData.RRN == "" {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[CheckReffIdSplitBill]-[Reference ID is Found : %v]", check.ReferenceId))
		logrus.Println(logReq)

		res = utils.GetMessageResponse(res, 172, false, errors.New("ReffId Not Found"))

		return res
	}

	textComment := param.TrxID + param.InstitutionID + constants.CodeReversal + "#" + "OP009 - Reversal point cause transaction " + getData.Voucher + " is failed"

	param.AccountId = getData.AccountId
	param.Point = getData.Point
	param.AccountNumber = getData.AccountNumber
	param.Comment = textComment
	param.RRN = req.ReferenceId
	param.TrxID = utils.GenTransactionId()

	reversal := AddingPointService(param, header)
	logrus.Info(fmt.Sprintf("[Response Reversal : %v]", reversal))

	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: map[string]interface{}{
			"referenceId": req.ReferenceId,
		},
	}

	return res

}
