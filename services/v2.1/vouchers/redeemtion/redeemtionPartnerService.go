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

func RedeemtionOrder_V21_Services(req models.VoucherComultaiveReq, param models.Params, head models.RequestHeader) models.Response {

	nameservice := "[PackageRedeemtion]-[RedeemtionOrder_V21_Services]"
	logReq := fmt.Sprintf("[AccountNumber : %v, RewardID : %v]", param.AccountNumber, param.RewardID)

	logrus.Info(nameservice)

	var res models.Response

	dataOrder := services.DataParameterOrder(constants.CODE_CONFIG_AGG_GROUP, constants.CODE_CONFIG_AGG_NAME, constants.CODE_CONFIG_AGG_EMAIL, constants.CODE_CONFIG_AGG_PHONE, constants.CODE_CONFIG_AGG_EXPD)

	timeExp, _ := strconv.Atoi(dataOrder.Expired)

	param.CumReffnum = utils.GenTransactionId()
	param.TrxID = utils.GenTransactionId()

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

	nama := "OTTOPOINT"
	reqOrder := vgmodels.RequestOrderVoucherAgV11{
		ProductCode:    param.ProductCodeInternal,
		Qty:            req.Jumlah,
		FieldValue:     field,
		OrderID:        param.CumReffnum,
		CustomerName:   nama,
		CustomerEmail:  dataOrder.Email,
		CustomerPhone:  dataOrder.Phone,
		DeliveryMethod: 1,
		RedeemCallback: hostPurcahse + callbackPartner,
	}

	if param.SupplierID == constants.CODE_VENDOR_UV {
		reqOrder.CustomerPhone = param.AccountNumber
	}

	fmt.Println("Start - OrderVoucherAggregator")
	logrus.Info("[VoucherAgService]-[OrderVoucher]")

	// Generate Signature
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

		for i := req.Jumlah; i > 0; i-- {
			// TrxId
			param.TrxID = utils.GenTransactionId()

			fmt.Println(fmt.Sprintf("[Line Save DB : %v]", i))

			t := i - 1
			coupon := RedeemVouchAG.CouponseVouch[t].CouponsID
			param.CouponID = coupon
			// go services.SaveTSchedulerRetry(param.TrxID, constants.CodeSchedulerVoucherAG)
			go services.SaveTransactionVoucherAgMigrate(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, constants.Pending, timeExp)
		}
		res = models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: models.UltraVoucherResp{
				Code:    "68",
				Msg:     "Transaksi Anda sedang dalam proses. Silahkan hubungi customer support kami untuk informasi selengkapnya.",
				Success: 0,
				Failed:  0,
				Pending: req.Jumlah,
			},
		}

		return res
	}

	// Handle Stock Not Available
	if order.ResponseCode == "04" {

		logrus.Info("[VoucherAgServices]-[OrderVoucher]")
		logrus.Info("[Stock not Available]-[Gagal Order Voucher]")
		logrus.Info("[Response Code ] : ", order.ResponseCode)
		logrus.Info("[Response Desc] : ", order.ResponseDesc)

		totalPoint := param.Point * req.Jumlah
		param.TrxID = utils.GenTransactionId()
		resultReversal := Trx.V21_Adding_PointVoucher(param, totalPoint, req.Jumlah, head)
		fmt.Println(resultReversal)

		fmt.Println("[ >>>>>>>>>>>>>>>>>>>>>>> Send Publisher <<<<<<<<<<<<<<<<<<<< ]")

		pubreq := models.NotifPubreq{
			Type:           constants.CODE_REVERSAL_POINT,
			NotificationTo: param.AccountNumber,
			Institution:    param.InstitutionID,
			ReferenceId:    param.RRN,
			TransactionId:  param.CumReffnum,
			Data: models.DataValue{
				RewardValue: "point",
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

		for i := req.Jumlah; i > 0; i-- {

			// TrxID
			param.TrxID = utils.GenTransactionId()

			fmt.Println(fmt.Sprintf("[Line Save DB : %v]", i))

			t := i - 1
			coupon := RedeemVouchAG.CouponseVouch[t].CouponsID
			param.CouponID = coupon

			go services.SaveTransactionVoucherAgMigrate(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, constants.Failed, timeExp)
		}

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

	// Handle Pending Status
	if order.ResponseCode == "09" {

		logrus.Info("[VoucherAgServices]-[OrderVoucher]")
		logrus.Info("[VoucherAggregator pending]-[OrderVoucher]")
		logrus.Info("[Response Code ] : ", order.ResponseCode)
		logrus.Info("[Response Desc] : ", order.ResponseDesc)

		for i := req.Jumlah; i > 0; i-- {
			fmt.Println(fmt.Sprintf("[Line : %v]", i))

			// TrxId
			param.TrxID = utils.GenTransactionId()

			fmt.Println(fmt.Sprintf("[Line Save DB : %v]", i))

			t := i - 1
			coupon := RedeemVouchAG.CouponseVouch[t].CouponsID
			param.CouponID = coupon
			// go services.SaveTSchedulerRetry(param.TrxID, constants.CodeSchedulerVoucherAG)
			go services.SaveTransactionVoucherAgMigrate(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, constants.Pending, timeExp)

		}

		res = models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: models.UltraVoucherResp{
				Code:    "09",
				Msg:     "Transaksi Anda sedang dalam proses. Silahkan hubungi customer support kami untuk informasi selengkapnya.",
				Success: 0,
				Failed:  0,
				Pending: req.Jumlah,
			},
		}

		return res
	}

	// Handle General Error
	if order.ResponseCode != "00" {

		// Reversal Start Here
		fmt.Println("Internal Server Error : ", errorder)
		fmt.Println("ResponseCode : ", order.ResponseCode)
		fmt.Println("[VoucherAgServices]-[FailedOrder]")
		fmt.Println(fmt.Sprintf("[Response %v]", order.ResponseCode))

		totalPoint := param.Point * req.Jumlah
		param.TrxID = utils.GenTransactionId()
		resultReversal := Trx.V21_Adding_PointVoucher(param, totalPoint, req.Jumlah, head)
		fmt.Println(resultReversal)

		fmt.Println("[ >>>>>>>>>>>>>>>>>>>>>>> Send Publisher <<<<<<<<<<<<<<<<<<<< ]")

		pubreq := models.NotifPubreq{
			Type:           constants.CODE_REVERSAL_POINT,
			NotificationTo: param.AccountNumber,
			Institution:    param.InstitutionID,
			ReferenceId:    param.RRN,
			TransactionId:  param.CumReffnum,
			Data: models.DataValue{
				RewardValue: "point",
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

		for i := req.Jumlah; i > 0; i-- {

			fmt.Println(fmt.Sprintf("[Line Save DB : %v]", i))

			// TrxID
			param.TrxID = utils.GenTransactionId()

			t := i - 1
			coupon := RedeemVouchAG.CouponseVouch[t].CouponsID
			param.CouponID = coupon

			go services.SaveTransactionVoucherAgMigrate(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, constants.Failed, timeExp)

		}

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

	for i := req.Jumlah; i > 0; i-- {
		fmt.Println(fmt.Sprintf("[Line : %v]", i))

		// TrxId
		param.TrxID = utils.GenTransactionId()

		fmt.Println(fmt.Sprintf("[Line Save DB : %v]", i))

		t := i - 1
		coupon := RedeemVouchAG.CouponseVouch[t].CouponsID
		param.CouponID = coupon

		go services.SaveTransactionVoucherAgMigrate(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, constants.Success, timeExp)

	}

	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: vgmodels.ResponseVoucherAg{
			Code:    "00",
			Msg:     "Success",
			Success: req.Jumlah,
			Failed:  0,
			Pending: 0,
		},
	}

	return res

}
