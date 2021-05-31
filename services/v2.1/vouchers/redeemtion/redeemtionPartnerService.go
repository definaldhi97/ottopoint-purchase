package redeemtion

import (
	"encoding/json"
	"fmt"
	"ottopoint-purchase/constants"
	kafka "ottopoint-purchase/hosts/publisher/host"
	signature "ottopoint-purchase/hosts/signature/host"
	vgmodels "ottopoint-purchase/hosts/voucher_aggregator/models"
	"ottopoint-purchase/models"

	"ottopoint-purchase/services/v2.1/Trx"

	// "ottopoint-purchase/services/v2_migrate"
	"ottopoint-purchase/services"
	"ottopoint-purchase/utils"
	"reflect"
	"strconv"

	vg "ottopoint-purchase/hosts/voucher_aggregator/host"

	"github.com/sirupsen/logrus"
)

var (
	hostPurcahse    = utils.GetEnv("OTTOPOINT_PURCHASE_HOST_CALLBACK_PARTNER", "http://34.101.175.164:8006")
	callbackPartner = utils.GetEnv("OTTOPOINT_PURCHASE_CALLBACK_PARTNER", "/transaction/callback/partner")
)

func RedeemtionOrder_V21_Services(req models.VoucherComultaiveReq, codeScheduler string, param models.Params, head models.RequestHeader) models.Response {

	nameservice := "[PackageRedeemtion]-[RedeemtionOrder_V21_Services]"
	logReq := fmt.Sprintf("[AccountNumber : %v, RewardID : %v]", param.AccountNumber, param.RewardID)

	logrus.Info(nameservice)

	var res models.Response

	dataOrder := services.DataParameterOrder(constants.CODE_CONFIG_AGG_GROUP, constants.CODE_CONFIG_AGG_NAME, constants.CODE_CONFIG_AGG_EMAIL, constants.CODE_CONFIG_AGG_PHONE, constants.CODE_CONFIG_AGG_EXPD)

	timeExp, _ := strconv.Atoi(dataOrder.Expired)

	param.CumReffnum = utils.GenTransactionId()

	param.Amount = int64(param.Point)

	// spending point and spending usage_limit voucher
	textCommentSpending := param.CumReffnum + "#" + param.NamaVoucher
	param.Comment = textCommentSpending

	field := []models.FieldsKey{}

	for line, v := range param.Fields {
		logrus.Info(">> Loop Field <<")

		a := models.FieldsKey{}

		logrus.Info(">> Field : ", v)
		if v == constants.CODE_NOMOR_TELP {

			logrus.Info(">> Field Phone line : ", line)

			switch string(req.CustID[0:2]) {
			case "08":
				a.Value = fmt.Sprintf("0%s", req.CustID[1:])
			case "62":
				a.Value = fmt.Sprintf("0%s", req.CustID[2:])
			default:
				a.Value = fmt.Sprintf("0%s", req.CustID)
			}

			a.Key = constants.CODE_NOMOR_TELP

			field = append(field, a)

		}

		// Wallet
		if v == constants.CODE_NOMOR_KARTU {

			logrus.Info(">> Field Nomor Kartu <<")

			a.Value = req.CustID
			a.Key = constants.CODE_NOMOR_KARTU

			field = append(field, a)
		}

		// PLN, Game (FF)
		if v == constants.CODE_ID_Pelanggan {

			logrus.Info(">> Field ID Pelanggan <<")

			a.Value = req.CustID
			a.Key = constants.CODE_ID_Pelanggan

			field = append(field, a)
		}

		// Game (ML)
		if v == constants.CODE_ID_Server {

			logrus.Info(">> Field ID Server <<")

			a.Value = req.CustID2
			a.Key = constants.CODE_ID_Server

			field = append(field, a)
		}

	}

	fmt.Println("Field : ", field)
	fmt.Println("Param Field : ", param.Fields)

	// return res

	// Potong Point & Poting Stock Voucher
	RedeemVouchAG, errRedeemVouchAG := Trx.V21_Redeem_PointandVoucher(req.Jumlah, param, head)

	param.PointTransferID = RedeemVouchAG.PointTransferID

	if RedeemVouchAG.Rc == "10" || RedeemVouchAG.Rd == "Insufficient Point" {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[V21_Redeem_PointandVoucher]-[Error : %v]", errRedeemVouchAG))
		logrus.Println("[Rc] : ", RedeemVouchAG.Rc)
		logrus.Println("[Rd] : ", RedeemVouchAG.Rd)
		logrus.Println(logReq)

		res = models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: vgmodels.ResponseVoucherAg{
				Code:    "60",
				Msg:     "Token or Session Expired Please Login Again",
				Success: 0,
				Failed:  req.Jumlah,
				Pending: 0,
			},
		}

		return res
	}

	if RedeemVouchAG.Rc == "208" {

		logrus.Info("[VoucherAgService]-[RedeemVoucher]")
		logrus.Info("[Rc] : ", RedeemVouchAG.Rc)
		logrus.Info("[Rd] : ", RedeemVouchAG.Rd)

		// res = utils.GetMessageResponse(res, 500, false, errors.New("Voucher Sudah Limit"))

		res = models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: vgmodels.ResponseVoucherAg{
				Code:    "65",
				Msg:     "Voucher not available",
				Success: 0,
				Failed:  req.Jumlah,
				Pending: 0,
			},
		}

		return res
	}

	if RedeemVouchAG.Rc == "209" {

		logrus.Info("[VoucherAgService]-[RedeemVoucher]")
		logrus.Error("Error : ", errRedeemVouchAG)
		logrus.Info("[ ResponseCode ] : ", RedeemVouchAG.Rc)
		logrus.Info("[ ResponseDesc ] : ", RedeemVouchAG.Rd)

		res = models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: vgmodels.ResponseVoucherAg{
				Code:    "66",
				Msg:     "Payment count limit exceeded",
				Success: 0,
				Failed:  req.Jumlah,
				Pending: 0,
			},
		}

		return res
	}

	if errRedeemVouchAG != nil || RedeemVouchAG.Rc != "00" {
		logrus.Info("[VoucherAgService]-[RedeemVoucher]")
		logrus.Info("[Rc] : ", RedeemVouchAG.Rc)
		logrus.Info("[Rd] : ", RedeemVouchAG.Rd)

		// res = utils.GetMessageResponse(res, 500, false, errors.New("Gagal! Maaf transaksi Anda tidak dapat dilakukan saat ini. Silahkan dicoba lagi atau hubungi tim kami untuk informasi selengkapnya."))
		res = models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: vgmodels.ResponseVoucherAg{
				Code:    "01",
				Msg:     "Gagal! Maaf transaksi Anda tidak dapat dilakukan saat ini. Silahkan dicoba lagi atau hubungi tim kami untuk informasi selengkapnya.",
				Success: 0,
				Failed:  req.Jumlah,
				Pending: 0,
			},
		}

		return res
	}

	var failed, success, pending int

	for i := req.Jumlah; i > 0; i-- {

		param.TrxID = utils.GenTransactionId()

		nama := "OTTOPOINT"
		reqOrder := vgmodels.RequestOrderVoucherAgV11{
			ProductCode:    param.ProductCodeInternal,
			Qty:            1,
			FieldValue:     field,
			OrderID:        param.TrxID,
			CustomerName:   nama,
			CustomerEmail:  dataOrder.Email,
			CustomerPhone:  dataOrder.Phone,
			DeliveryMethod: 1,
			RedeemCallback: hostPurcahse + callbackPartner,
		}

		// if param.SupplierID == constants.CODE_VENDOR_UV {
		// 	reqOrder.CustomerPhone = param.AccountNumber
		// }

		fmt.Println("Start - OrderVoucherAggregator")
		logrus.Info("[VoucherAgService]-[OrderVoucher]")

		// Generate Signature
		head.ChannelID = "H2H"
		sign, err := signature.Signature(reqOrder, head)
		if err != nil {
			res = models.Response{
				Meta: utils.ResponseMetaOK(),
				Data: utils.GetMessageFailedErrorNew(
					res,
					constants.RC_ERROR_INVALID_SIGNATURE,
					constants.RD_ERROR_INVALID_SIGNATURE,
				),
			}
			return res
		}

		// Get Signature from interface{}
		s := reflect.ValueOf(sign.Data)
		for _, k := range s.MapKeys() {
			head.Signature = fmt.Sprintf("%s", s.MapIndex(k))
		}

		order, errorder := vg.OrderVoucherV11(reqOrder, head)

		logrus.Info("Response Order : ", order)

		param.DataSupplier.Rd = order.ResponseDesc
		param.DataSupplier.Rc = order.ResponseCode

		if errorder != nil || order.ResponseCode == "" {
			// Reversal Start Here
			logrus.Info("[VoucherAgServices]-[OrderVoucher]")
			logrus.Error("[Failed Order Voucher]-[Gagal Order Voucher] : ", errorder.Error())

			// for _, val := range RedeemVouchAG.CouponseVouch {
			// 	param.CouponID = val.CouponsID
			// }

			t := i - 1
			param.CouponID = RedeemVouchAG.CouponseVouch[t].CouponsID

			go services.SaveTSchedulerRetry(param.TrxID, codeScheduler)
			go services.SaveTransactionVoucherAgMigrate(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, constants.Pending, timeExp)

			// res = models.Response{
			// 	Meta: utils.ResponseMetaOK(),
			// 	Data: models.UltraVoucherResp{
			// 		Code:    "68",
			// 		Msg:     "Transaksi Anda sedang dalam proses. Silahkan hubungi customer support kami untuk informasi selengkapnya.",
			// 		Success: 0,
			// 		Failed:  0,
			// 		Pending: req.Jumlah,
			// 	},
			// }

			// return res

			pending++

			continue
		}

		// Handle Stock Not Available
		if order.ResponseCode == "04" {

			logrus.Info("[VoucherAgServices]-[OrderVoucher]")
			logrus.Info("[Stock not Available]-[Gagal Order Voucher]")
			logrus.Info("[Response Code ] : ", order.ResponseCode)
			logrus.Info("[Response Desc] : ", order.ResponseDesc)

			// for _, val := range RedeemVouchAG.CouponseVouch {
			// 	param.CouponID = val.CouponsID
			// }

			t := i - 1
			param.CouponID = RedeemVouchAG.CouponseVouch[t].CouponsID

			go services.SaveTransactionVoucherAgMigrate(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, constants.Failed, timeExp)

			// totalPoint := param.Point
			// resultReversal := Trx.V21_Adding_PointVoucher(param, totalPoint, req.Jumlah, head)
			// fmt.Println(resultReversal)

			// fmt.Println("[ >>>>>>>>>>>>>>>>>>>>>>> Send Publisher <<<<<<<<<<<<<<<<<<<< ]")

			// pubreq := models.NotifPubreq{
			// 	Type:           constants.CODE_REVERSAL_POINT,
			// 	NotificationTo: param.AccountNumber,
			// 	Institution:    param.InstitutionID,
			// 	ReferenceId:    param.RRN,
			// 	TransactionId:  param.TrxID,
			// 	Data: models.DataValue{
			// 		RewardValue: "point",
			// 		Value:       strconv.Itoa(totalPoint),
			// 	},
			// }

			// bytePub, _ := json.Marshal(pubreq)
			// kafkaReq := kafka.PublishReq{
			// 	Topic: utils.TopicsNotif,
			// 	Value: bytePub,
			// }

			// kafkaRes, err := kafka.SendPublishKafka(kafkaReq)
			// if err != nil {
			// 	fmt.Println("Gagal Send Publisher")
			// 	fmt.Println("Error : ", err)
			// }

			// fmt.Println("Response Publisher : ", kafkaRes)

			// res = models.Response{
			// 	Meta: utils.ResponseMetaOK(),
			// 	Data: vgmodels.ResponseVoucherAg{
			// 		Code:    "01",
			// 		Msg:     "Gagal! Maaf transaksi Anda tidak dapat dilakukan saat ini. Silahkan dicoba lagi atau hubungi tim kami untuk informasi selengkapnya.",
			// 		Success: 0,
			// 		Failed:  req.Jumlah,
			// 		Pending: 0,
			// 	},
			// }

			// return res

			failed++
			continue

		}

		// Handle Pending Status
		if order.ResponseCode == "09" {

			logrus.Info("[VoucherAgServices]-[OrderVoucher]")
			logrus.Info("[VoucherAggregator pending]-[OrderVoucher]")
			logrus.Info("[Response Code ] : ", order.ResponseCode)
			logrus.Info("[Response Desc] : ", order.ResponseDesc)

			// for _, val := range RedeemVouchAG.CouponseVouch {
			// 	param.CouponID = val.CouponsID
			// }

			t := i - 1
			param.CouponID = RedeemVouchAG.CouponseVouch[t].CouponsID

			go services.SaveTSchedulerRetry(param.TrxID, codeScheduler)
			go services.SaveTransactionVoucherAgMigrate(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, constants.Pending, timeExp)

			// res = models.Response{
			// 	Meta: utils.ResponseMetaOK(),
			// 	Data: models.UltraVoucherResp{
			// 		Code:    "09",
			// 		Msg:     "Transaksi Anda sedang dalam proses. Silahkan hubungi customer support kami untuk informasi selengkapnya.",
			// 		Success: 0,
			// 		Failed:  0,
			// 		Pending: req.Jumlah,
			// 	},
			// }

			// return res

			pending++
			continue
		}

		// Handle General Error
		if order.ResponseCode != "00" {

			// Reversal Start Here
			fmt.Println("Internal Server Error : ", errorder)
			fmt.Println("ResponseCode : ", order.ResponseCode)
			fmt.Println("[VoucherAgServices]-[FailedOrder]")
			fmt.Println(fmt.Sprintf("[Response %v]", order.ResponseCode))

			// for _, val := range RedeemVouchAG.CouponseVouch {
			// 	param.CouponID = val.CouponsID
			// }

			t := i - 1
			param.CouponID = RedeemVouchAG.CouponseVouch[t].CouponsID

			go services.SaveTransactionVoucherAgMigrate(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, constants.Failed, timeExp)

			// res = models.Response{
			// 	Meta: utils.ResponseMetaOK(),
			// 	Data: vgmodels.ResponseVoucherAg{
			// 		Code:    "01",
			// 		Msg:     "Gagal! Maaf transaksi Anda tidak dapat dilakukan saat ini. Silahkan dicoba lagi atau hubungi tim kami untuk informasi selengkapnya.",
			// 		Success: 0,
			// 		Failed:  req.Jumlah,
			// 		Pending: 0,
			// 	},
			// }

			// return res

			failed++
			continue

		}

		fmt.Println(fmt.Sprintf("[Line : %v]", i))

		// TrxId
		// param.TrxID = utils.GenTransactionId()

		fmt.Println(fmt.Sprintf("[Line Save DB : %v]", i))

		// for _, val := range RedeemVouchAG.CouponseVouch {
		// 	param.CouponID = val.CouponsID
		// }

		t := i - 1
		param.CouponID = RedeemVouchAG.CouponseVouch[t].CouponsID

		if param.SupplierID == constants.CODE_VENDOR_AGREGATOR {

			// Check Order Status
			statusOrder, errStatus := vg.CheckStatusOrder(vgmodels.RequestCheckOrderStatus{
				OrderID:       param.TrxID,
				RecordPerPage: "1",
				CurrentPage:   "1",
			}, head)

			if errStatus != nil {

				// Handle Error Here
				fmt.Println("[VoucherAgServices]-[FailedCheckOrderStatus]")
				fmt.Println("Internal Server Error : ", errStatus)
				fmt.Println("ResponseCode : ", statusOrder.ResponseCode)
				fmt.Println(fmt.Sprintf("[Response %v]", statusOrder.ResponseCode))

			}

			param.RRN = statusOrder.Data.TransactionID

			for i := req.Jumlah; i > 0; i-- {
				fmt.Println(fmt.Sprintf("[Line Save DB : %v]", i))

				t := i - 1
				coupon := RedeemVouchAG.CouponseVouch[t].CouponsID
				param.CouponID = coupon
				// code := RedeemVouchAG.CouponseVouch[t].CouponsCode

				// voucherID := statusOrder.Data.Vouchers[t].VoucherID
				voucherCode := statusOrder.Data.Vouchers[t].VoucherCode
				expDate := statusOrder.Data.Vouchers[t].ExpiredDate
				voucherLink := statusOrder.Data.Vouchers[t].Link

				a := []rune(param.CouponID)
				key32 := string(a[0:32])
				key := []byte(key32)
				chiperText := []byte(voucherCode)
				plaintText, err := utils.EncryptAES(chiperText, key)
				if err != nil {
					res = utils.GetMessageFailedErrorNew(res, constants.RC_FAILED_DECRYPT_VOUCHER, constants.RD_FAILED_DECRYPT_VOUCHER)
					return res
				}

				// Use Voucher ID as a Transaction ID
				param.TrxID = utils.GenTransactionId()
				param.ExpDate = expDate
				param.CouponCode = fmt.Sprintf("%s", plaintText)
				param.VoucherLink = voucherLink

				id := utils.GenerateTokenUUID()
				go services.SaveDBVoucherAgMigrate(id, param.InstitutionID, param.CouponID, voucherCode, param.AccountNumber, param.AccountId, req.CampaignID)

			}

		}

		go services.SaveTransactionVoucherAgMigrate(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, constants.Success, timeExp)

		success++

	}

	if failed != 0 {

		fmt.Println(fmt.Sprintf("[ Reversal CumReff : %v ]", param.CumReffnum))

		go func(param models.Params, failed int) {

			totalPoint := param.Point * failed
			param.TrxID = utils.GenTransactionId()
			resultReversal := Trx.V21_Adding_PointVoucher(param, totalPoint, failed, head)
			fmt.Println(resultReversal)

			fmt.Println("[ >>>>>>>>>>>>>>>>>>>>>>> Send Publisher <<<<<<<<<<<<<<<<<<<< ]")

			pubreq := models.NotifPubreq{
				Type:           constants.CODE_REVERSAL_POINT,
				NotificationTo: param.AccountNumber,
				Institution:    param.InstitutionID,
				ReferenceId:    param.RRN,
				TransactionId:  param.CumReffnum,
				Data: models.DataValue{
					RewardValue: constants.Point,
					Value:       strconv.Itoa(totalPoint),
				},
			}

			bytePub, _ := json.Marshal(pubreq)

			kafkaReq := kafka.PublishReq{
				Topic: utils.TopicsNotif,
				Value: bytePub,
			}

			kafkaRes, err := kafka.SendPublishKafka(kafkaReq)
			if err != nil {
				fmt.Println("Gagal Send Publisher")
				fmt.Println("Error : ", err)
			}

			fmt.Println("Response Publisher : ", kafkaRes)
		}(param, failed)

	}

	var code, msg string

	if success != 0 {
		code = "00"
		msg = "Transaksi Berhasil"
	} else if failed != 0 {
		code = "01"
		msg = "Transaksi Gagal"
	} else if pending != 0 {
		code = "09"
		msg = "Transaksi Sedang Dalam Proses"
	}

	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: vgmodels.ResponseVoucherAg{
			Code:    code,
			Msg:     msg,
			Success: success,
			Failed:  failed,
			Pending: pending,
		},
	}

	return res

}

func msgRedeemtionCummulative(success, pending, failed int) (string, string) {

	var Code_RC_Comulative, Message_Comulative string
	// Sukses
	if (success != 0) && (pending == 0) && (failed == 0) {
		Code_RC_Comulative = "00"
		Message_Comulative = "Transaksi Berhasil"

		return Code_RC_Comulative, Message_Comulative
	}

	// Sukses & Gagal
	if (success != 0) && (pending == 0) && (failed != 0) {
		Code_RC_Comulative = "174"
		Message_Comulative = fmt.Sprintf("%v Voucher Anda berhasil dirukar namun %v voucher tidak berhasil. Poin yang tidak digunakan akan dikembalikan ke saldo Anda", success, failed)

		return Code_RC_Comulative, Message_Comulative
	}

	// Sukses & Pending
	if (success != 0) && (pending != 0) && (failed == 0) {
		Code_RC_Comulative = "175"
		Message_Comulative = fmt.Sprintf("%v Voucher Anda berhasil ditukar & %v Transaksi Anda sedang dalam proses", success, pending)

		return Code_RC_Comulative, Message_Comulative
	}

	// Sukses & Pending & Gagal
	if (success != 0) && (pending != 0) && (failed != 0) {
		Code_RC_Comulative = "33"
		Message_Comulative = fmt.Sprintf("%v Vucher Anda berhasil ditukar namun %v Voucher pending dan %v voucher tidak berhasil. Harap hubungi customer support untuk informasi lebih lanjut.", success, pending, failed)
		// Message_Comulative = fmt.Sprintf("%v Voucher Anda berhasil dirukar namun %v voucher tidak berhasil. Poin yang tidak digunakan akan dikembalikan ke saldo Anda", success, failed)

		return Code_RC_Comulative, Message_Comulative
	}

	// Pending
	if (success == 0) && (pending != 0) && (failed == 0) {
		Code_RC_Comulative = "56"
		Message_Comulative = fmt.Sprintf("%v Transaksi Anda sedang dalam proses. Silahkan hubungi tim kami untuk informasi selengkapnya.", pending)

		return Code_RC_Comulative, Message_Comulative
	}

	// Pending & Gagal
	if (success == 0) && (pending != 0) && (failed != 0) {
		Code_RC_Comulative = "57"
		Message_Comulative = fmt.Sprintf("%v Transaksi Anda sedang dalam proses & %v Transaksi Anda Gagal.Poin yang tidak digunakan akan dikembalikan ke saldo Anda", success, failed)

		return Code_RC_Comulative, Message_Comulative
	}

	// Gagal
	if (success == 0) && (pending == 0) && (failed != 0) {
		Code_RC_Comulative = "01"
		Message_Comulative = "Transaksi Gagal"

		return Code_RC_Comulative, Message_Comulative
	}

	return Code_RC_Comulative, Message_Comulative
}
