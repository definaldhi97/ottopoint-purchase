package callbacks

import (
	"encoding/json"
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	kafka "ottopoint-purchase/hosts/publisher/host"
	sepulsaModels "ottopoint-purchase/hosts/sepulsa/models"
	"ottopoint-purchase/models"
	"ottopoint-purchase/services/v2/Trx"
	"ottopoint-purchase/utils"
	"time"

	"github.com/sirupsen/logrus"
)

// func (t V2_VoucherSepulsaService) CallbackVoucherSepulsa(req sepulsaModels.CallbackTrxReq) models.Response {
func CallbackVoucherSepulsaService(req sepulsaModels.CallbackTrxReq) models.Response {
	// fmt.Println("[ >>>>>>>>>>>>>>>>>> V2 Migrate CallBack Sepulsa Service <<<<<<<<<<<<<<<< ]")
	var res models.Response

	nameservice := "[PackageCallbacks]-[CallbackVoucherSepulsa_V21_Service]"

	logReq := fmt.Sprintf("[TransactionID : %v, CustomerNumber : %v]", req.TransactionID, req.CustomerNumber)

	fmt.Println("Start Delay ", time.Now().Unix())
	time.Sleep(10 * time.Second)

	go func(args sepulsaModels.CallbackTrxReq) {
		// Get Spending By TransactionID and OrderID
		spending, errSpending := db.GetSpendingSepulsa(args.TransactionID, args.OrderID)
		if errSpending != nil {

			logrus.Error(nameservice)
			logrus.Error(fmt.Sprintf("[GetSpendingSepulsa]-[Error : %v]", errSpending))
			logrus.Println(logReq)

		}

		responseCode := models.GetErrorMsg(args.ResponseCode)

		logrus.Info("[HandleCallbackSepulsa] - [ResponseCode] : ", args.ResponseCode)
		logrus.Info("[HandleCallbackSepulsa] - [ResponseDesc] : ", responseCode)

		param := models.Params{
			InstitutionID: spending.Institution,
			NamaVoucher:   spending.Voucher,
			AccountId:     spending.AccountId,
			AccountNumber: spending.AccountNumber,
			RRN:           spending.RRN,
			TrxID:         spending.TransactionId,
			RewardID:      *spending.MRewardID,
			Point:         spending.Point,
		}

		if (responseCode != "Success") && (responseCode != "Pending") {

			resultReversal := Trx.V2_Adding_PointVoucher(param, spending.Point, 1)
			fmt.Println(resultReversal)

			fmt.Println("[ >>>>>>>>>>>>>>>>>>>>>>> Send Publisher <<<<<<<<<<<<<<<<<<<< ]")

			pubreq := models.NotifPubreq{
				Type:           constants.CODE_REVERSAL_POINT,
				NotificationTo: spending.AccountNumber,
				Institution:    spending.Institution,
				ReferenceId:    spending.RRN,
				TransactionId:  spending.CummulativeRef,
				Data: models.DataValue{
					RewardValue: "point",
					Value:       fmt.Sprint(spending.Point),
				},
			}

			bytePub, _ := json.Marshal(pubreq)

			kafkaReq := kafka.PublishReq{
				Topic: constants.TOPIC_PUSHNOTIF_GENERAL,
				Value: bytePub,
			}

			kafkaRes, errKafka := kafka.SendPublishKafka(kafkaReq)
			if errKafka != nil {

				logrus.Error(nameservice)
				logrus.Error(fmt.Sprintf("[SendPublishKafka]-[Error : %v]", errKafka))
				logrus.Println(logReq)

			}
			logrus.Info("[ Response Publisher ] : ", kafkaRes)

		}

		responseSepulsa, _ := json.Marshal(args)

		// Update TSpending
		_, errUpdate := db.UpdateVoucherSepulsa(responseCode, args.ResponseCode, string(responseSepulsa), args.TransactionID, args.OrderID)

		if errUpdate != nil {

			logrus.Error(nameservice)
			logrus.Error(fmt.Sprintf("[UpdateVoucherSepulsa]-[Error : %v]", errUpdate))
			logrus.Println(logReq)
		}

		// Update TSchedulerRetry
		_, err := db.UpdateTSchedulerRetry(spending.RRN)
		if err != nil {

			logrus.Error(nameservice)
			logrus.Error(fmt.Sprintf("[UpdateTSchedulerRetry]-[Error : %v]", err))
			logrus.Println(logReq)
		}

	}(req)

	fmt.Println("End Process ", time.Now().Unix())
	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: nil,
	}

	return res
}
