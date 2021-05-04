package payment

import (
	"encoding/json"
	"errors"
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	kafka "ottopoint-purchase/hosts/publisher/host"
	"ottopoint-purchase/models"
	sp "ottopoint-purchase/models/v2/payment"
	"ottopoint-purchase/utils"
	"strconv"

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
	param.ProductName = getData.Voucher

	reversal := AddingPointService(param, header)
	logrus.Info(fmt.Sprintf("[Response Reversal : %v]", reversal))

	fmt.Println("[ >>>>>>>>>>>>>>>>>>>>>>> Send Publisher <<<<<<<<<<<<<<<<<<<< ]")
	pubreq := models.NotifPubreq{
		Type:           constants.CODE_REVERSAL_POINT,
		NotificationTo: param.AccountNumber,
		Institution:    param.InstitutionID,
		ReferenceId:    param.RRN,
		TransactionId:  param.TrxID,
		Data: models.DataValue{
			RewardValue: constants.Point,
			Value:       strconv.Itoa(param.Point),
		},
	}

	bytePub, _ := json.Marshal(pubreq)

	kafkaReq := kafka.PublishReq{
		Topic: utils.TopicsNotif,
		Value: bytePub,
	}

	kafkaRes, errKafka := kafka.SendPublishKafka(kafkaReq)
	if errKafka != nil {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[SendPublishKafka]-[Error : %v]", errKafka))
		logrus.Println(logReq)

	}

	logrus.Info("[ Response Publisher ] : ", kafkaRes)

	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: map[string]interface{}{
			"referenceId": req.ReferenceId,
		},
	}

	return res

}
