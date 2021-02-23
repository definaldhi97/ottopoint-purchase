package services

import (
	"encoding/json"
	"fmt"
	"log"
	"ottopoint-purchase/constants"
	db "ottopoint-purchase/db"
	auth "ottopoint-purchase/hosts/auth/host"
	"ottopoint-purchase/hosts/opl/host"
	opl "ottopoint-purchase/hosts/opl/host"
	kafka "ottopoint-purchase/hosts/publisher/host"
	sepulsa "ottopoint-purchase/hosts/sepulsa/host"
	sepulsaModels "ottopoint-purchase/hosts/sepulsa/models"
	"ottopoint-purchase/models"
	"ottopoint-purchase/models/dbmodels"
	"ottopoint-purchase/utils"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/opentracing/opentracing-go"
	"github.com/vjeantet/jodaTime"
	"go.uber.org/zap"
)

type UseSepulsaService struct {
	General models.GeneralModel
}

func (t UseSepulsaService) SepulsaServices(req models.VoucherComultaiveReq, param models.Params) models.Response {
	var res models.Response

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[SepulsaServices]",
		zap.String("NameVoucher : ", param.NamaVoucher), zap.Int("Jumlah : ", req.Jumlah),
		zap.String("CampaignID : ", req.CampaignID), zap.String("CampaignID : ", req.CampaignID),
		zap.String("CustID2 : ", req.CustID2), zap.String("ProductCode : ", param.ProductCode),
		zap.String("AccountNumber : ", param.AccountNumber), zap.String("InstitutionID : ", param.InstitutionID))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[SepulsaServices]")
	defer span.Finish()

	param.CumReffnum = utils.GenTransactionId()

	total := strconv.Itoa(req.Jumlah)
	param.Amount = int64(param.Point)

	// redeem to opl (potong point)
	redeem, errredeem := host.RedeemVoucherCumulative(req.CampaignID, param.AccountId, total, "0")
	log.Printf("Redeem Coupons: %v", redeem.Coupons)
	if redeem.Message == "Invalid JWT Token" {
		fmt.Println("Error : ", errredeem)
		fmt.Println("[SepulsaVoucherService]-[RedeemVoucher]")
		fmt.Println("[Internal Server Error]-[Gagal Redeem Voucher]")

		sugarLogger.Info("[SepulsaVoucherService]-[RedeemVoucher]")
		sugarLogger.Info("[Internal Server Error]-[Gagal Redeem Voucher]")

		res := models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: models.SepulsaRes{
				Code:    "65",
				Msg:     "Token or session Expired Please Login Again",
				Success: 0,
				Failed:  req.Jumlah,
				Pending: 0,
			},
		}

		return res
	}

	if redeem.Error == "Not enough points" {
		fmt.Println("Error : ", errredeem)
		fmt.Println("[SepulsaVoucherService]-[RedeemVoucher]")
		fmt.Println("[Not enough points]-[Gagal Redeem Voucher]")

		sugarLogger.Info("[SepulsaVoucherService]-[RedeemVoucher]")
		sugarLogger.Info("[Not enough points]-[Gagal Redeem Voucher]")

		res := models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: models.SepulsaRes{
				Code:    "27",
				Msg:     "Point Tidak Mencukupi",
				Success: 0,
				Failed:  req.Jumlah,
				Pending: 0,
			},
		}

		return res
	}

	if redeem.Error == "Limit exceeded" {
		fmt.Println("Error : ", errredeem)
		fmt.Println("[SepulsaVoucherService]-[RedeemVoucher]")
		fmt.Println("[Limit exceed]-[Gagal Redeem Voucher]")

		sugarLogger.Info("[SepulsaVoucherService]-[RedeemVoucher]")
		sugarLogger.Info("[Limit exceed]-[Gagal Redeem Voucher]")

		res := models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: models.SepulsaRes{
				Code:    "65",
				Msg:     "Payment count limit exceed",
				Success: 0,
				Failed:  req.Jumlah,
				Pending: 0,
			},
		}

		return res
	}

	var coupon string
	for _, val := range redeem.Coupons {
		coupon = val.Code
	}

	if errredeem != nil || redeem.Error != "" || coupon == "" {
		fmt.Println("Error : ", errredeem)
		fmt.Println("[SepulsaVoucherService]-[RedeemVoucher]")
		fmt.Println("[Failed Redeem Voucher]-[Gagal Redeem Voucher]")

		sugarLogger.Info("[SepulsaVoucherService]-[RedeemVoucher]")
		sugarLogger.Info("[Failed Redeem Voucher]-[Gagal Redeem Voucher]")

		res = models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: models.SepulsaRes{
				Code:    "01",
				Msg:     "Gagal! Maaf transaksi Anda tidak dapat dilakukan saat ini. Silahkan dicoba lagi atau hubungi tim kami untuk informasi selengkapnya",
				Success: 0,
				Failed:  req.Jumlah,
				Pending: 0,
			},
		}

		return res
	}

	for i := req.Jumlah; i > 0; i-- {

		param.TrxID = utils.GenTransactionId()

		t := i - 1
		couponID := redeem.Coupons[t].Id
		couponCode := redeem.Coupons[t].Code
		param.CouponID = couponID

		productID, _ := strconv.Atoi(param.ProductCode)
		reqOrder := sepulsaModels.EwalletInsertTrxReq{
			CustomerNumber: req.CustID,
			OrderID:        param.TrxID,
			ProductID:      productID,
		}

		// Create Transaction Ewallet
		sepulsaRes, errTransaction := sepulsa.EwalletInsertTransaction(reqOrder)
		if errTransaction != nil {
			fmt.Println("[SepulsaService]-[InsertTransaction]")
			fmt.Println("ResponseDesc : ", errTransaction.Error())

			sugarLogger.Info("[SepulsaService]-[InsertTransaction]")
			sugarLogger.Info(fmt.Sprintf("[SepulsaService]-[FailedInsertTransaction]-[%v", errTransaction.Error()))

			// Start Reversal
			text := param.TrxID + param.InstitutionID + constants.CodeReversal + "#" + "OP009 - Reversal point couse transaction " + param.NamaVoucher + " is failed"

			schedulerData := dbmodels.TSchedulerRetry{
				Code:          constants.CodeScheduler,
				TransactionID: utils.Before(text, "#"),
				Count:         0,
				IsDone:        false,
				CreatedAT:     time.Now(),
			}

			sendReversal, errReversal := host.TransferPoint(param.AccountId, strconv.Itoa(param.Point), text)
			statusEarning := constants.Success
			msgEarning := constants.MsgSuccess

			if errReversal != nil || sendReversal.PointsTransferId == "" {

				statusEarning = constants.TimeOut

				fmt.Println(fmt.Sprintf("===== Failed TransferPointOPL to %v || TrxID : %v =====", param.AccountNumber, param.TrxID))

				for _, val1 := range sendReversal.Form.Children.Customer.Errors {
					if val1 != "" {
						msgEarning = val1
						statusEarning = constants.Failed
					}
				}

				for _, val2 := range sendReversal.Form.Children.Points.Errors {
					if val2 != "" {
						msgEarning = val2
						statusEarning = constants.Failed
					}
				}

				if sendReversal.Message != "" {
					msgEarning = sendReversal.Message
					statusEarning = constants.Failed
				}

				if sendReversal.Error.Message != "" {
					msgEarning = sendReversal.Error.Message
					statusEarning = constants.Failed
				}

				if statusEarning == constants.TimeOut {
					errSaveScheduler := db.DbCon.Create(&schedulerData).Error
					if errSaveScheduler != nil {

						fmt.Println("===== Gagal SaveScheduler ke DB =====")
						fmt.Println(fmt.Sprintf("Error : %v", errSaveScheduler))
						fmt.Println(fmt.Sprintf("===== Phone : %v || RRN : %v =====", param.AccountNumber, param.RRN))

						// return
					}

				}

			}

			expired := ExpiredPointService()

			saveReversal := dbmodels.TEarning{
				ID:               utils.GenerateTokenUUID(),
				PartnerId:        param.InstitutionID,
				TransactionId:    param.TrxID,
				AccountNumber:    param.AccountNumber,
				Point:            int64(param.Point),
				Status:           statusEarning,
				StatusMessage:    msgEarning,
				PointsTransferId: sendReversal.PointsTransferId,
				TransType:        constants.CodeReversal,
				AccountId:        param.AccountId,
				ExpiredPoint:     expired,
				TransactionTime:  time.Now(),
			}

			errSaveReversal := db.DbCon.Create(&saveReversal).Error
			if errSaveReversal != nil {

				fmt.Println(fmt.Sprintf("[Failed Save Reversal to DB]-[Error : %v]", errSaveReversal))
				fmt.Println("[PackageServices]-[SaveEarning]")

				fmt.Println(">>> Save CSV <<<")
				name := jodaTime.Format("dd-MM-YYYY", time.Now()) + ".csv"
				go utils.CreateCSVFile(saveReversal, name)

			}

			fmt.Println("========== Send Publisher ==========")

			pubreq := models.NotifPubreq{
				Type:           constants.CODE_REVERSAL_POINT,
				NotificationTo: param.AccountNumber,
				Institution:    param.InstitutionID,
				ReferenceId:    param.TrxID,
				TransactionId:  param.CumReffnum,
				Data: models.DataValue{
					RewardValue: "point",
					Value:       fmt.Sprint(param.Point),
				},
			}

			bytePub, _ := json.Marshal(pubreq)

			kafkaReq := kafka.PublishReq{
				Topic: constants.TOPIC_PUSHNOTIF_GENERAL,
				Value: bytePub,
			}

			kafkaRes, err := kafka.SendPublishKafka(kafkaReq)
			if err != nil {
				fmt.Println("Gagal Send Publisher")
				fmt.Println("Error : ", err)
			}

			// Save Error Transaction
			go SaveTransactionSepulsa(param, errTransaction.Error(), reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "01")

			fmt.Println("Response Publisher : ", kafkaRes)
			res = models.Response{
				Meta: utils.ResponseMetaOK(),
				Data: models.SepulsaRes{
					Code:    "01",
					Msg:     "Gagal! Maaf transaksi Anda tidak dapat dilakukan saat ini. Silahkan dicoba lagi atau hubungi tim kami untuk informasi selengkapnya.",
					Success: 0,
					Failed:  req.Jumlah,
					Pending: 0,
				},
			}
			return res

		}

		param.DataSupplier.Rd = sepulsaRes.Status
		param.DataSupplier.Rc = sepulsaRes.ResponseCode
		param.RRN = sepulsaRes.TransactionID

		// Use Voucher to OpenLoyalty
		// go CouponVoucherCustomer(req.CampaignID, couponID, couponCode, param.AccountId, 1)
		use, err2 := opl.CouponVoucherCustomer(req.CampaignID, couponID, couponCode, param.AccountId, 1)
		if err2 != nil {

			logrus.Info(fmt.Sprintf("[Error : %v]", err2))
			logrus.Info(fmt.Sprintf("[Response : %v]", use))
			logrus.Info("[Error from OPL]-[CouponVoucherCustomer]")

			sugarLogger.Info("[SepulsaService]-[CouponVoucherCustomer]")
			sugarLogger.Info(fmt.Sprintf("[SepulsaService]-[FailedCouponVoucherCustomer]-[%v", err2.Error()))

			// Save Scheduler Retry With Code SC003
			text := param.RRN + param.InstitutionID + constants.CodeReversal + "#" + "OP009 - Reversal point couse transaction " + param.NamaVoucher + " is failed"

			schedulerData := dbmodels.TSchedulerRetry{
				Code:          constants.CodeSchedulerSepulsa,
				TransactionID: utils.Before(text, "#"),
				Count:         0,
				IsDone:        false,
				CreatedAT:     time.Now(),
			}

			errSaveScheduler := db.DbCon.Create(&schedulerData).Error
			if errSaveScheduler != nil {

				fmt.Println("===== Gagal SaveScheduler ke DB =====")
				fmt.Println(fmt.Sprintf("Error : %v", errSaveScheduler))
				fmt.Println(fmt.Sprintf("===== Phone : %v || RRN : %v =====", param.AccountNumber, param.RRN))

			}

		}

		id := utils.GenerateTokenUUID()
		go SaveTSchedulerRetry(param)
		go SaveDBSepulsa(id, param.InstitutionID, couponID, couponCode, param.AccountNumber, param.AccountId, req.CampaignID)
		go SaveTransactionSepulsa(param, sepulsaRes, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "09")

	}

	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: models.SepulsaRes{
			Code:    "00",
			Msg:     fmt.Sprintf("Selamat Penukaran %s Kamu Berhasil, Silahkan Cek Saldo Kamu!", param.NamaVoucher),
			Success: req.Jumlah,
			Failed:  0,
			Pending: 0,
		},
	}

	return res
}

func (t UseSepulsaService) HandleCallbackRequest(req sepulsaModels.CallbackTrxReq) models.Response {
	var res models.Response

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[SepulsaService]",
		zap.String("TransactionID : ", req.TransactionID), zap.String("OrderID : ", req.OrderID),
		zap.String("Status : ", req.Status), zap.String("Desc : ", req.ResponseCode),
	)

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[SepulsaService]")
	defer span.Finish()

	fmt.Println("Start Delay ", time.Now().Unix())
	time.Sleep(10 * time.Second)
	go func(args sepulsaModels.CallbackTrxReq) {

		// Get Spending By TransactionID and OrderID
		spending, err := db.GetSpendingSepulsa(args.TransactionID, args.OrderID)
		if err != nil {
			fmt.Println("[GetSpendingSepulsa] : ", err.Error())

			sugarLogger.Info("[SepulsaService]-[GetSpendingSepulsa]")
			sugarLogger.Info(fmt.Sprintf("[SepulsaService]-[FailedGetSpendingSepulsa]-[%v", err.Error()))
		}

		responseCode := models.GetErrorMsg(args.ResponseCode)

		fmt.Println("[HandleCallbackSepulsa] - [ResponseCode] : ", args.ResponseCode)
		fmt.Println("[HandleCallbackSepulsa] - [ResponseDesc] : ", responseCode)

		go t.clearCacheBalance(spending.AccountNumber)

		if (responseCode != "Success") &&
			(responseCode != "Pending") {

			text := args.OrderID + spending.Institution + constants.CodeScheduler + "#" + "OP09 - Reversal point cause transaction " + spending.Voucher + " is failed"

			schedulerData := dbmodels.TSchedulerRetry{
				Code:          constants.CodeScheduler,
				TransactionID: utils.Before(text, "#"),
				Count:         0,
				IsDone:        false,
				CreatedAT:     time.Now(),
			}

			fmt.Sprintln("[HandleCallbackSepulsa]-[CreateScheduler] : ", schedulerData)
			sugarLogger.Info(fmt.Sprintf("[HandleCallbackSepulsa]-[CreateScheduler] : %v", schedulerData))

			// Start Reversal Point
			point := strconv.Itoa(spending.Point)
			sendReversal, errReversal := host.TransferPoint(spending.AccountId, point, text)

			statusEarning := constants.Success
			msgEarning := constants.MsgSuccess

			if errReversal != nil || sendReversal.PointsTransferId == "" {

				statusEarning = constants.TimeOut

				fmt.Println(fmt.Sprintf("===== Failed TransferPointOPL to %v || RRN : %v =====", spending.AccountNumber, spending.RRN))

				for _, val1 := range sendReversal.Form.Children.Customer.Errors {
					if val1 != "" {
						msgEarning = val1
						statusEarning = constants.Failed
					}
				}

				for _, val2 := range sendReversal.Form.Children.Points.Errors {
					if val2 != "" {
						msgEarning = val2
						statusEarning = constants.Failed
					}
				}

				if sendReversal.Message != "" {
					msgEarning = sendReversal.Message
					statusEarning = constants.Failed
				}

				if sendReversal.Error.Message != "" {
					msgEarning = sendReversal.Error.Message
					statusEarning = constants.Failed
				}

				if statusEarning == constants.TimeOut {
					errSaveScheduler := db.DbCon.Create(&schedulerData).Error
					if errSaveScheduler != nil {

						fmt.Println("===== Gagal SaveScheduler ke DB =====")
						fmt.Println(fmt.Sprintf("Error : %v", errSaveScheduler))
						fmt.Println(fmt.Sprintf("===== Phone : %v || RRN : %v =====", spending.AccountNumber, spending.RRN))

						// return
					}

				}

			}

			expired := ExpiredPointService()

			saveReversal := dbmodels.TEarning{
				ID: utils.GenerateTokenUUID(),
				// EarningRule     :,
				// EarningRuleAdd  :,
				PartnerId: spending.Institution,
				// ReferenceId     : ,
				TransactionId: utils.GenTransactionId(),
				// ProductCode     :,
				// ProductName     :,
				AccountNumber: spending.AccountNumber,
				// Amount          :,
				Point: int64(spending.Point),
				// Remark          :,
				Status:           statusEarning,
				StatusMessage:    msgEarning,
				PointsTransferId: sendReversal.PointsTransferId,
				// RequestorData   :,
				// ResponderData   :,
				TransType:       constants.CodeReversal,
				AccountId:       spending.AccountId,
				ExpiredPoint:    expired,
				TransactionTime: time.Now(),
			}

			errSaveReversal := db.DbCon.Create(&saveReversal).Error
			if errSaveReversal != nil {

				fmt.Println(fmt.Sprintf("[Failed Save Reversal to DB]-[Error : %v]", errSaveReversal))
				fmt.Println("[PackageServices]-[SaveEarning]")

				fmt.Println(">>> Save CSV <<<")
				name := jodaTime.Format("dd-MM-YYYY", time.Now()) + ".csv"
				go utils.CreateCSVFile(saveReversal, name)

			}

			fmt.Println("========== Send Publisher ==========")

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

			kafkaRes, err := kafka.SendPublishKafka(kafkaReq)
			if err != nil {
				fmt.Println("Gagal Send Publisher")
				fmt.Println("Error : ", err)
			}

			fmt.Println("Response Publisher : ", kafkaRes)

		}

		responseSepulsa, _ := json.Marshal(args)

		// Update TSpending
		res, errUpdate := db.UpdateVoucherSepulsa(responseCode, args.ResponseCode, string(responseSepulsa), args.TransactionID, args.OrderID)
		if errUpdate != nil {
			fmt.Println("[UpdateVoucherSepulsa] : ", errUpdate.Error())

			sugarLogger.Info("[SepulsaService]-[UpdateVoucherSepulsa]")
			sugarLogger.Info(fmt.Sprintf("[SepulsaService]-[FailedUpdateVoucherSepulsa]-[%v", errUpdate.Error()))
		}

		fmt.Sprintln("[SuccessUpdateVoucherSepulsa] : ", res)
		sugarLogger.Info("[SepulsaService]-[SuccessUpdateVoucherSepulsa]")

		transactionID := spending.RRN + spending.Institution + constants.CodeReversal + "#" + "OP009 - Reversal point couse transaction " + spending.Voucher + " is failed"

		// Update TSchedulerRetry
		_, err = db.UpdateTSchedulerRetry(utils.Before(transactionID, "#"))
		if err != nil {
			fmt.Println("[UpdateTSchedulerRetry] : ", err.Error())

			sugarLogger.Info("[SepulsaService]-[UpdateTSchedulerRetry]")
			sugarLogger.Info(fmt.Sprintf("[SepulsaService]-[FailedUpdateTSchedulerRetry]-[%v", err.Error()))
		}

	}(req)

	fmt.Println("End Process ", time.Now().Unix())
	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: nil,
	}

	return res
}

func (t UseSepulsaService) clearCacheBalance(phone string) {
	fmt.Println(">>>>>>> Clear Cache Get Balance <<<<<<")
	clearCacheBalance, err := auth.ClearCacheBalance(phone)
	if err != nil {
		fmt.Println("Clear Cache Balance Error : ", err)
		return
	}
	if clearCacheBalance.ResponseCode != "00" {
		fmt.Println("Message : ", clearCacheBalance.Messages)
		fmt.Println("Response Code : ", clearCacheBalance.ResponseCode)
		return
	}
	fmt.Println("Clear Cache Get Balance: ", clearCacheBalance.Messages)
	return

}

func SaveTSchedulerRetry(param models.Params) {

	fmt.Println(fmt.Sprintf("[Start-SaveTSchedulerRetry]-[Sepulsa]"))

	text := param.RRN + param.InstitutionID + constants.CodeReversal + "#" + "OP009 - Reversal point couse transaction " + param.NamaVoucher + " is failed"
	schedulerData := dbmodels.TSchedulerRetry{
		Code:          constants.CodeSchedulerSepulsa,
		TransactionID: utils.Before(text, "#"),
		Count:         0,
		IsDone:        false,
		CreatedAT:     time.Now(),
	}

	errSaveScheduler := db.DbCon.Create(&schedulerData).Error
	if errSaveScheduler != nil {

		fmt.Println("===== Gagal SaveScheduler ke DB =====")
		fmt.Println(fmt.Sprintf("Error : %v", errSaveScheduler))
		fmt.Println(fmt.Sprintf("===== Phone : %v || RRN : %v =====", param.AccountNumber, param.RRN))

	}

}

func SaveDBSepulsa(id, institution, coupon, vouchercode, phone, custIdOPL, campaignID string) {
	fmt.Println("[SaveDB]-[SepulsaVoucherService]")
	save := dbmodels.UserMyVocuher{
		ID:            id,
		InstitutionID: institution,
		CouponID:      coupon,
		VoucherCode:   vouchercode,
		Phone:         phone,
		AccountId:     custIdOPL,
		CampaignID:    campaignID,
	}

	err := db.DbCon.Create(&save).Error
	if err != nil {
		fmt.Println("[Failed Save to DB]", err)
		fmt.Println("[Package-Service]-[SepulsaService]")
	}
}

func SaveTransactionSepulsa(param models.Params, res interface{}, reqdata interface{}, reqOP models.VoucherComultaiveReq, transType, status string) {

	fmt.Println(fmt.Sprintf("[Start-SaveDB]-[Sepulsa]-[%v]", transType))

	var ExpireDate time.Time
	var redeemDate time.Time

	var saveStatus string
	switch status {
	case "00":
		saveStatus = constants.Success
	case "09":
		saveStatus = constants.Pending
	case "01":
		saveStatus = constants.Failed
	}

	if transType == constants.CODE_TRANSTYPE_REDEMPTION {
		ExpireDate = utils.ExpireDateVoucherAGt(constants.EXPDATE_VOUCHER)
		redeemDate = time.Now()
	}

	reqSepulsa, _ := json.Marshal(&reqdata)
	responseSepulsa, _ := json.Marshal(&res)
	reqdataOP, _ := json.Marshal(&reqOP)

	// timeRedeem := jodaTime.Format("dd-MM-YYYY HH:mm:ss", time.Now())

	save := dbmodels.TSpending{
		ID:                utils.GenerateTokenUUID(),
		AccountNumber:     param.AccountNumber,
		Voucher:           param.NamaVoucher,
		MerchantID:        param.MerchantID,
		CustID:            reqOP.CustID,
		RRN:               param.RRN,
		TransactionId:     param.TrxID,
		ProductCode:       param.ProductCode,
		Amount:            int64(param.Amount),
		TransType:         transType,
		IsUsed:            true,
		UsedAt:            &redeemDate,
		ProductType:       param.ProductType,
		Status:            saveStatus,
		ExpDate:           utils.DefaultNulTime(ExpireDate),
		Institution:       param.InstitutionID,
		CummulativeRef:    param.CumReffnum,
		DateTime:          utils.GetTimeFormatYYMMDDHHMMSS(),
		Point:             param.Point,
		ResponderRc:       param.DataSupplier.Rc,
		ResponderRd:       param.DataSupplier.Rd,
		RequestorData:     string(reqSepulsa),
		ResponderData:     string(responseSepulsa),
		RequestorOPData:   string(reqdataOP),
		SupplierID:        param.SupplierID,
		CouponId:          param.CouponID,
		CampaignId:        param.CampaignID,
		AccountId:         param.AccountId,
		RedeemAt:          utils.DefaultNulTime(redeemDate),
		Comment:           param.Comment,
		MRewardID:         param.RewardID,
		ProductCategoryID: param.CategoryID,
		MProductID:        param.ProductID,
		PointsTransferID:  param.PointTransferID,
	}

	err := db.DbCon.Create(&save).Error
	if err != nil {
		logrus.Info(fmt.Sprintf("[Error : %v]", err))
		logrus.Info("[Failed Save to DB]")

		name := jodaTime.Format("dd-MM-YYYY", time.Now()) + ".csv"
		go utils.CreateCSVFile(save, name)

		// return err

	}

}

func CouponVoucherCustomer(campaignID, couponID, couponCode, accountID string, useVoucher int) {
	use, err2 := opl.CouponVoucherCustomer(campaignID, couponID, couponCode, accountID, 1)
	if err2 != nil {

		logrus.Info(fmt.Sprintf("[Error : %v]", err2))
		logrus.Info(fmt.Sprintf("[Response : %v]", use))
		logrus.Info("[Error from OPL]-[CouponVoucherCustomer]")

	}
}
