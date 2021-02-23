package redeemtion

import (
	"encoding/json"
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	kafka "ottopoint-purchase/hosts/publisher/host"
	sepulsa "ottopoint-purchase/hosts/sepulsa/host"
	sepulsaModels "ottopoint-purchase/hosts/sepulsa/models"
	"ottopoint-purchase/models"
	v2_redeemtion "ottopoint-purchase/services/v2/vouchers/redeemtion"

	V21_trx "ottopoint-purchase/services/v2.1/Trx"

	// "ottopoint-purchase/services/v2_migrate"
	"ottopoint-purchase/utils"
	"strconv"

	"github.com/sirupsen/logrus"
)

// func (t V21_VoucherSepulsaService) V21_VoucherSepulsa(req models.VoucherComultaiveReq, param models.Params, header models.RequestHeader) models.Response {
func RedeemtionSepulsa_V21_Service(req models.VoucherComultaiveReq, param models.Params, header models.RequestHeader) models.Response {
	// fmt.Println("[ >>>>>>>>>>>>>>>>>> Migrate V2.1 Voucher Sepulsa Service <<<<<<<<<<<<<<<< ]")

	nameservice := "[PackageRedeemtionService]-[RedeemtionSepulsa_V21_Service]"
	logReq := fmt.Sprintf("[AccountNumber : %v, RewardID : %v]", param.AccountNumber, param.RewardID)

	logrus.Info(nameservice)

	var res models.Response

	param.CumReffnum = utils.GenTransactionId()

	// validasi usage limit voucher
	dtaVocher, _ := db.Get_MReward(param.CampaignID)

	// validasi stock voucher
	if req.Jumlah > dtaVocher.UsageLimit {

		logrus.Error(nameservice)
		logrus.Error("[ Stock Voucher not Available ]")
		logrus.Println(logReq)

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

		logrus.Error(nameservice)
		logrus.Error("[ Payment count limit exceeded ]")
		logrus.Println(logReq)

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

		var c string
		for _, vall := range RedeemVouchSP.CouponseVouch {
			c = vall.CouponsCode
		}
		fmt.Println("Value CouponCode : ", c)

		if errRedeemVouchSP != nil || RedeemVouchSP.Rc != "00" {

			logrus.Error(nameservice)
			logrus.Error("[ Payment count limit exceeded ]")
			logrus.Error(fmt.Sprintf("[V21_Redeem_PointandVoucher]-[Error : %v]", errRedeemVouchSP))
			logrus.Println(logReq)

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

		var couponID, couponCode string
		if RedeemVouchSP.Rc == "00" {
			couponID = RedeemVouchSP.CouponseVouch[0].CouponsID
			couponCode = RedeemVouchSP.CouponseVouch[0].CouponsCode
			param.CouponID = couponID
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

			logrus.Error(nameservice)
			logrus.Error(fmt.Sprintf("[EwalletInsertTransaction]-[Error : %v]", errTransaction))
			logrus.Println(logReq)

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

			kafkaRes, errKafka := kafka.SendPublishKafka(kafkaReq)
			if errKafka != nil {

				logrus.Error(nameservice)
				logrus.Error(fmt.Sprintf("[SendPublishKafka]-[Error : %v]", errKafka))
				logrus.Println(logReq)

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
