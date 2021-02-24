package Redeemtion

import (
	"encoding/json"
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	kafka "ottopoint-purchase/hosts/publisher/host"
	sepulsa "ottopoint-purchase/hosts/sepulsa/host"
	sepulsaModels "ottopoint-purchase/hosts/sepulsa/models"
	"ottopoint-purchase/models"
	V21_trx "ottopoint-purchase/services/v2.1/Trx"
	v2_redeemtion "ottopoint-purchase/services/v2/Redeemtion"

	// "ottopoint-purchase/services/v2_migrate"
	"ottopoint-purchase/utils"
	"strconv"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

type V21_VoucherSepulsaService struct {
	General models.GeneralModel
}

func (t V21_VoucherSepulsaService) V21_VoucherSepulsa(req models.VoucherComultaiveReq, param models.Params, header models.RequestHeader) models.Response {
	fmt.Println("[ >>>>>>>>>>>>>>>>>> Migrate V2.1 Voucher Sepulsa Service <<<<<<<<<<<<<<<< ]")

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

	// validasi usage limit voucher
	dtaVocher, _ := db.Get_MReward(param.CampaignID)

	// validasi stock voucher
	if req.Jumlah > dtaVocher.UsageLimit {
		fmt.Println("[ Stock Voucher not Available ]")
		logrus.Info("[SepulsaVoucherService]-[RedeemVoucher")

		res := models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: models.SepulsaRes{
				Code:    "65",
				Msg:     "Voucher not available",
				Success: 0,
				Failed:  req.Jumlah,
				Pending: 0,
			},
		}

		return res
	}

	// validasi limit voucher per user
	countRedeemed, _ := db.GetVoucherRedeemed(param.AccountId, param.RewardID)
	if countRedeemed.Count > dtaVocher.LimitPeruser {
		logrus.Info("[SepulsaVoucherService]-[RedeemVoucher]")
		fmt.Println("[ Payment count limit exceeded ]")

		res = models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: models.SepulsaRes{
				Code:    "66",
				Msg:     "Payment count limit exceeded",
				Success: 0,
				Failed:  req.Jumlah,
				Pending: 0,
			},
		}

		return res
	}

	// process order/trx to sepulsa
	for i := req.Jumlah; i > 0; i-- {

		param.Amount = int64(param.Point)

		param.TrxID = utils.GenTransactionId()

		// spending point and spending usage_limit voucher
		textCommentSpending := param.TrxID + "#" + param.NamaVoucher
		param.Comment = textCommentSpending
		RedeemVouchSP, errRedeemVouchSP := V21_trx.V21_Redeem_PointandVoucher(1, param, header)
		param.PointTransferID = RedeemVouchSP.PointTransferID

		// if RedeemVouchSP.Rc == "10" || RedeemVouchSP.Rd == "Insufficient Point" {

		// 	logrus.Info("[SepulsaVoucherService]-[RedeemVoucher]")
		// 	logrus.Info("[Not enough points]-[Gagal Redeem Voucher]")
		// 	logrus.Info("[Rc] : ", RedeemVouchSP.Rc)
		// 	logrus.Info("[Rd] : ", RedeemVouchSP.Rd)

		// 	// res = utils.GetMessageResponse(res, 500, false, errors.New("Point Tidak Cukup"))
		// 	res = models.Response{
		// 		Meta: utils.ResponseMetaOK(),
		// 		Data: models.UltraVoucherResp{
		// 			Code:    "27",
		// 			Msg:     "Point Tidak Mencukupi",
		// 			Success: 0,
		// 			Failed:  req.Jumlah,
		// 			Pending: 0,
		// 		},
		// 	}

		// 	return res
		// }

		// if RedeemVouchSP.Rc == "208" {

		// 	logrus.Info("[SepulsaVoucherService]-[RedeemVoucher")
		// 	logrus.Error("Error : ", errRedeemVouchSP)
		// 	logrus.Info("[Rc] : ", RedeemVouchSP.Rc)
		// 	logrus.Info("[Rd] : ", RedeemVouchSP.Rd)

		// 	res := models.Response{
		// 		Meta: utils.ResponseMetaOK(),
		// 		Data: models.SepulsaRes{
		// 			Code:    "65",
		// 			Msg:     "Voucher not available",
		// 			Success: 0,
		// 			Failed:  req.Jumlah,
		// 			Pending: 0,
		// 		},
		// 	}

		// 	return res
		// }

		// if RedeemVouchSP.Rc == "209" {

		// 	logrus.Info("[SepulsaVoucherService]-[RedeemVoucher]")
		// 	logrus.Error("Error : ", RedeemVouchSP)
		// 	logrus.Info("[ ResponseCode ] : ", RedeemVouchSP.Rc)
		// 	logrus.Info("[ ResponseDesc ] : ", RedeemVouchSP.Rd)

		// 	res = models.Response{
		// 		Meta: utils.ResponseMetaOK(),
		// 		Data: models.SepulsaRes{
		// 			Code:    "66",
		// 			Msg:     "Payment count limit exceeded",
		// 			Success: 0,
		// 			Failed:  req.Jumlah,
		// 			Pending: 0,
		// 		},
		// 	}

		// 	return res
		// }

		var c string
		for _, vall := range RedeemVouchSP.CouponseVouch {
			c = vall.CouponsCode
		}
		fmt.Println("Value CouponCode : ", c)

		if errRedeemVouchSP != nil || RedeemVouchSP.Rc != "00" {

			logrus.Info("[SepulsaVoucherService]-[RedeemVoucher]")
			logrus.Error("Error : ", errRedeemVouchSP)
			logrus.Info("[Failed Redeem Voucher]-[Gagal Redeem Voucher]")

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

		//////////////////////////////////////-----Before----//////////////////////////////////

		// for i := req.Jumlah; i > 0; i-- {

		// param.TrxID = utils.GenTransactionId()

		// t := i - 1
		// couponID := RedeemVouchSP.CouponseVouch[t].CouponsID
		// couponCode := RedeemVouchSP.CouponseVouch[t].CouponsCode
		// param.CouponID = couponID

		var couponID, couponCode string
		if RedeemVouchSP.Rc == "00" {
			couponID = RedeemVouchSP.CouponseVouch[0].CouponsID
			couponCode = RedeemVouchSP.CouponseVouch[0].CouponsCode
			param.CouponID = couponID
		}

		for _, v := range param.Fields {
			if v == constants.CODE_NOMOR_TELP {
				switch string(req.CustID[0:2]) {
				case "08":
					req.CustID = fmt.Sprintf("0%s", req.CustID[1:])
				case "62":
					req.CustID = fmt.Sprintf("0%s", req.CustID[2:])
				default:
					req.CustID = fmt.Sprintf("0%s", req.CustID)
				}
			}
		}

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

			// Save Error Transaction
			go v2_redeemtion.SaveTransactionSepulsa(param, errTransaction.Error(), reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "01")

			trxIDReversal := utils.GenTransactionId()
			param.TrxID = trxIDReversal

			resultReversal := V21_trx.V21_Adding_PointVoucher(param, param.Point, 1, header)
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

			// // Save Error Transaction
			// go v2_redeemtion.SaveTransactionSepulsa(param, errTransaction.Error(), reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "01")

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
		go v2_redeemtion.SaveTSchedulerRetry(param.RRN, constants.CodeSchedulerSepulsa)
		go v2_redeemtion.SaveDBSepulsa(id, param.InstitutionID, couponID, couponCode, param.AccountNumber, param.AccountId, req.CampaignID)
		go v2_redeemtion.SaveTransactionSepulsa(param, sepulsaRes, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "09")

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

func (t V21_VoucherSepulsaService) V21_CallbackVoucherSepulsa(req sepulsaModels.CallbackTrxReq) models.Response {
	fmt.Println("[ >>>>>>>>>>>>>>>>>> Migrate V2.1 CallBack Sepulsa Service <<<<<<<<<<<<<<<< ]")
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
			TrxID:         utils.GenTransactionId(),
			RewardID:      spending.MRewardID,
			Point:         spending.Point,
		}

		header := models.RequestHeader{
			DeviceID:      "ottopoint-purchase",
			InstitutionID: spending.Institution,
			Geolocation:   "-",
			ChannelID:     "H2H",
			AppsID:        "-",
			Timestamp:     utils.GetTimeFormatYYMMDDHHMMSS(),
			Authorization: "-",
			Signature:     "-",
		}

		if (responseCode != "Success") && (responseCode != "Pending") {

			resultReversal := V21_trx.V21_Adding_PointVoucher(param, spending.Point, 1, header)
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
