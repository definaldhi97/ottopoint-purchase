package callbacks

import (
	"encoding/json"
	"errors"
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	kafka "ottopoint-purchase/hosts/publisher/host"
	"ottopoint-purchase/models"
	callback "ottopoint-purchase/models/v21/callback"
	"ottopoint-purchase/utils"
	"strings"

	"github.com/sirupsen/logrus"
)

type DataCallbackNotif struct {
	NotifType     string
	InstitutionID string
	ReffNum       string
	TransactionID string
	AccountNumber string
	VoucherCode   string
	ProductName   string
}

func CallbackVoucherAG_V21_Service(req callback.CallbackVoucherAGReq) models.Response {
	var res models.Response

	nameservice := "[PackageCallbacks]-[CallbackVoucherAG_V21_Service]"

	logReq := fmt.Sprintf("[TransactionID : %v]", req.TransactionId)

	logrus.Info(nameservice)

	// validate TrxID
	dataTrx, errTrx := db.CheckTrxbyTrxID(req.OrderId)
	logrus.Println(">>> CheckTrxbyTrxID <<<")
	if errTrx != nil {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[CheckTrxId]-[Error : %v]", errTrx))
		logrus.Println(logReq)

		res = utils.GetMessageResponse(res, 422, false, errors.New("TrxId Tidak Ditemukan"))

		return res

	}

	logrus.Println(">>> Lanjut <<<")

	reUpdate := db.VoucherTypeDB{}

	// PPOB (1)
	if req.VoucherType == constants.VoucherTypePPOB {

		logrus.Println(">>> PPOB <<<")

		dataPPOB := callback.DataVoucherTypePPOB{}

		data1, _ := json.Marshal(&req.Data)

		errPPOB := json.Unmarshal(data1, &dataPPOB)
		fmt.Println("Error Marshall Data PPOB : ", errPPOB)

		reUpdate = db.VoucherTypeDB{
			VoucherType:  1,
			OrderId:      req.TransactionId,
			ResponseCode: dataPPOB.ResponseCode,
			ResponseDesc: dataPPOB.ResponseDesc,
		}

		update := db.UpdateVoucherbyVoucherType(reUpdate, req.OrderId, req)
		logrus.Info("Response Update : ", update)
	} else {

		dataVouchercode := callback.DataVoucherTypeVoucherCode{}

		data2, _ := json.Marshal(&req.Data)

		reUpdate = db.VoucherTypeDB{
			VoucherType:  2,
			OrderId:      req.TransactionId,
			VoucherId:    dataVouchercode.VoucherID,
			VoucherCode:  dataVouchercode.VoucherCode,
			VoucherName:  dataVouchercode.VoucherName,
			IsRedeemed:   dataVouchercode.IsRedeemed,
			RedeemedDate: dataVouchercode.RedeemedDate,
		}

		errVouchercode := json.Unmarshal(data2, &dataVouchercode)
		fmt.Println("Error Marshall Data errVoucherCode : ", errVouchercode)

		update := db.UpdateVoucherbyVoucherType(reUpdate, req.OrderId, req)
		logrus.Info("Response Update : ", update)

		dataVoucher := DataCallbackNotif{
			InstitutionID: req.InstitutionId,
			ReffNum:       req.TransactionId,
			TransactionID: req.OrderId,
			AccountNumber: dataTrx.AccountNumber,
			VoucherCode:   dataVouchercode.VoucherCode,
		}

		if strings.ToLower(dataTrx.ProductType) == strings.ToLower(constants.CategoryPLN) {
			dataVoucher.NotifType = constants.CODE_REDEEM_PLN
		}

		if strings.ToLower(dataTrx.ProductType) == strings.ToLower(constants.CategoryVidio) {
			dataVoucher.NotifType = constants.CODE_REDEEM_VIDIO
		}

		go sendNotifDataVoucher(dataVoucher)

		return res

	}

	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: nil,
	}

	return res
}

func sendNotifDataVoucher(dataVoucher DataCallbackNotif) {

	var topics string

	reqNotif := models.NotifPubreq{
		Type:           dataVoucher.NotifType,
		NotificationTo: dataVoucher.AccountNumber,
		Institution:    dataVoucher.InstitutionID,
		ReferenceId:    dataVoucher.ReffNum,
		TransactionId:  dataVoucher.TransactionID,
	}

	dataNotif := models.DataValue{}

	dataSMS := models.DataValueSMS{}

	dataIssuer, _ := db.GetDataInstitution(dataVoucher.InstitutionID)

	// SMS
	if dataIssuer.NOtificationID == constants.CODE_SMS_NOTIF || dataIssuer.NOtificationID == constants.CODE_SMS_APPS_NOTIF {

		topics = utils.TopicNotifSMS

		dataNotif.RewardValue = dataVoucher.ProductName
		dataNotif.Value = dataVoucher.VoucherCode

	}

	// Notif APK
	if dataIssuer.NOtificationID == constants.CODE_APPS_NOTIF || dataIssuer.NOtificationID == constants.CODE_SMS_APPS_NOTIF {

		topics = utils.TopicsNotif

		dataSMS.ProductName = dataVoucher.ProductName
		dataSMS.Token = dataVoucher.VoucherCode
	}

	bytePub, _ := json.Marshal(reqNotif)

	kafkaReq := kafka.PublishReq{
		Topic: topics,
		Value: bytePub,
	}

	kafkaRes, errKafka := kafka.SendPublishKafka(kafkaReq)
	if errKafka != nil {

		fmt.Println(fmt.Sprintf("[PackageCallbacks]-[sendNotifDataVoucher]"))
		fmt.Println(fmt.Sprintf("[SendPublishKafka]-[Error : %v]", errKafka))

	}

	fmt.Println("Response Publisher : ", kafkaRes)
}
