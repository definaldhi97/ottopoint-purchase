package callbacks

import (
	"encoding/json"
	"fmt"
	"ottopoint-purchase/db"
	sepulsaModels "ottopoint-purchase/hosts/sepulsa/models"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"
	"time"

	"github.com/sirupsen/logrus"
)

// func V21_CallbackVoucherSepulsa(req sepulsaModels.CallbackTrxReq) models.Response {
// 	fmt.Println("[ >>>>>>>>>>>>>>>>>> Migrate V2.1 CallBack Sepulsa Service <<<<<<<<<<<<<<<< ]")
func CallbackVoucherSepulsa_V21_Service(req sepulsaModels.CallbackTrxReq) models.Response {
	var res models.Response

	nameservice := "[PackageCallbacks]-[CallbackVoucherSepulsa_V21_Service]"

	logReq := fmt.Sprintf("[TransactionID : %v, CustomerNumber : %v]", req.TransactionID, req.CustomerNumber)

	logrus.Info(nameservice)

	logrus.Println("Start Delay ", time.Now().Unix())
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

		if (responseCode != "Success") && (responseCode != "Pending") {

			logrus.
				WithField("OrderID", req.OrderID).
				WithField("TransactionID", spending.TransactionId).
				WithField("CustomerNumber", req.CustomerNumber).
				WithField("Status", req.Status).
				WithField("ResponseCode", responseCode).
				Warn("Failed Order Sepulsa")

		}

		responseSepulsa, _ := json.Marshal(args)

		// Update TSpending
		_, errUpdate := db.UpdateVoucherSepulsa(responseCode, args.ResponseCode, string(responseSepulsa), args.TransactionID, args.OrderID)

		if errUpdate != nil {

			logrus.Error(nameservice)
			logrus.Error(fmt.Sprintf("[UpdateVoucherSepulsa]-[Error : %v]", errUpdate))
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
