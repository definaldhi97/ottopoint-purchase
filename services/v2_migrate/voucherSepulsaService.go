package v2_migrate

import (
	"encoding/json"
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	kafka "ottopoint-purchase/hosts/publisher/host"
	sepulsa "ottopoint-purchase/hosts/sepulsa/host"
	sepulsaModels "ottopoint-purchase/hosts/sepulsa/models"
	"ottopoint-purchase/models"
	"ottopoint-purchase/models/dbmodels"
	"ottopoint-purchase/utils"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/vjeantet/jodaTime"
	"go.uber.org/zap"
)

type VoucherSepulsaMigrateService struct {
	General models.GeneralModel
}

func (t VoucherSepulsaMigrateService) VoucherSepulsa(req models.VoucherComultaiveReq, param models.Params) models.Response {
	fmt.Println("[ >>>>>>>>>>>>>>>>>> Voucher Sepulsa Service <<<<<<<<<<<<<<<< ]")

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

	// total := strconv.Itoa(req.Jumlah)
	param.Amount = int64(param.Point)

	// spending point and spending usage_limit voucher
	textCommentSpending := param.CumReffnum + "#" + param.NamaVoucher
	param.Comment = textCommentSpending
	RedeemVouchSP, errRedeemVouchSP := Redeem_PointandVoucher(req.Jumlah, param)

	if RedeemVouchSP.Rd == "Invalid JWT Token" {
		fmt.Println("Error : ", errRedeemVouchSP)
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

	if RedeemVouchSP.Rd == "not enough points" {
		fmt.Println("Error : ", errRedeemVouchSP)
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

	if RedeemVouchSP.Rd == "Voucher not available" {
		fmt.Println("Error : ", errRedeemVouchSP)
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

	var c string
	for _, vall := range RedeemVouchSP.CouponseVouch {
		c = vall.CouponsCode
	}

	if errRedeemVouchSP != nil || RedeemVouchSP.Rc != "00" || c == "" {
		fmt.Println("Error : ", errRedeemVouchSP)
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
		couponID := RedeemVouchSP.CouponseVouch[t].CouponsID
		couponCode := RedeemVouchSP.CouponseVouch[t].CouponsCode
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
			logrus.Info("[SepulsaService]-[InsertTransaction]")
			logrus.Error("ResponseDesc : ", errTransaction.Error())

			resultReversal := Adding_PointVoucher(param, param.Point, 1)
			fmt.Println(resultReversal)

			fmt.Println("[ >>>>>>>>>>>>>>>>>>>>>>> Send Publisher <<<<<<<<<<<<<<<<<<<< ]")
			pubreq := models.NotifPubreq{
				Type:           constants.CODE_REVERSAL_POINT,
				NotificationTo: param.AccountNumber,
				Institution:    param.InstitutionID,
				ReferenceId:    param.RRN,
				TransactionId:  param.Reffnum,
				Data: models.DataValue{
					RewardValue: "point",
					Value:       strconv.Itoa(param.Point),
				},
			}

			bytePub, _ := json.Marshal(pubreq)

			kafkaReq := kafka.PublishReq{
				Topic: utils.TopicsNotif,
				Value: bytePub,
			}

			kafkaRes, err := kafka.SendPublishKafka(kafkaReq)
			if err != nil {
				logrus.Error("Gagal Send Publisher : ", err)
			}

			logrus.Info("[ Response Publisher ] : ", kafkaRes)

			// Save Error Transaction
			go SaveTransactionSepulsa(param, errTransaction.Error(), reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "01")

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

		id := utils.GenerateTokenUUID()
		go SaveTSchedulerRetry(param.RRN, constants.CodeSchedulerSepulsa)
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

func (t VoucherSepulsaMigrateService) CallbackVoucherSepulsa(req sepulsaModels.CallbackTrxReq) models.Response {
	fmt.Println("[ >>>>>>>>>>>>>>>>>> CallBack Sepulsa Service <<<<<<<<<<<<<<<< ]")
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
			logrus.Error("[ Failed Get SpendingSepulsa ] : ", err.Error())
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
			RewardID:      spending.MRewardID,
			Point:         spending.Point,
		}

		if (responseCode != "Success") && (responseCode != "Pending") {

			resultReversal := Adding_PointVoucher(param, spending.Point, 1)
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

			kafkaRes, err := kafka.SendPublishKafka(kafkaReq)
			if err != nil {
				logrus.Error("Gagal Send Publisher : ", err)
			}
			logrus.Info("[ Response Publisher ] : ", kafkaRes)

		}

		responseSepulsa, _ := json.Marshal(args)

		// Update TSpending
		_, errUpdate := db.UpdateVoucherSepulsa(responseCode, args.ResponseCode, string(responseSepulsa), args.TransactionID, args.OrderID)

		if errUpdate != nil {
			logrus.Error("[UpdateVoucherSepulsa] : ", errUpdate.Error())

		}

		// Update TSchedulerRetry
		_, err = db.UpdateTSchedulerRetry(spending.RRN)
		if err != nil {

			logrus.Error("[SepulsaService]-[FailedUpdateTSchedulerRetry] : ", errUpdate.Error())
		}

	}(req)

	fmt.Println("End Process ", time.Now().Unix())
	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: nil,
	}

	return res
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
		ExpDate:           &ExpireDate,
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
		RedeemAt:          &redeemDate,
		Comment:           param.Comment,
		MRewardID:         param.RewardID,
		ProductCategoryID: param.CategoryID,
		MProductID:        param.ProductID,
	}

	err := db.DbCon.Create(&save).Error
	if err != nil {
		logs.Info(fmt.Sprintf("[Error : %v]", err))
		logs.Info("[Failed Save to DB]")

		name := jodaTime.Format("dd-MM-YYYY", time.Now()) + ".csv"
		go utils.CreateCSVFile(save, name)

		// return err

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

func SaveTSchedulerRetry(trxID, code string) {

	fmt.Println(fmt.Sprintf("[Start-SaveTSchedulerRetry]-[Sepulsa]"))

	schedulerData := dbmodels.TSchedulerRetry{
		Code:          code,
		TransactionID: trxID,
		Count:         0,
		IsDone:        false,
		CreatedAT:     time.Now(),
	}

	errSaveScheduler := db.DbCon.Create(&schedulerData).Error
	if errSaveScheduler != nil {

		fmt.Println("===== Gagal SaveScheduler ke DB =====")
		fmt.Println(fmt.Sprintf("Error : %v", errSaveScheduler))

	}

}
