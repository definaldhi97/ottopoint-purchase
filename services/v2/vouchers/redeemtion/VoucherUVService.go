package redeemtion

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"ottopoint-purchase/constants"
	kafka "ottopoint-purchase/hosts/publisher/host"
	uv "ottopoint-purchase/hosts/ultra_voucher/host"
	uvmodels "ottopoint-purchase/hosts/ultra_voucher/models"
	"ottopoint-purchase/models"
	"ottopoint-purchase/services"
	"ottopoint-purchase/services/v2/Trx"
	"ottopoint-purchase/utils"

	"github.com/sirupsen/logrus"
)

// func (t V2_VoucherUVService) VoucherUV(req models.VoucherComultaiveReq, param models.Params) models.Response {
func RedeemtionUVServices(req models.VoucherComultaiveReq, param models.Params) models.Response {
	// fmt.Println("[ >>>>>>>>>>>>>>>>>> V2 Migrate Voucher UV Service <<<<<<<<<<<<<<<< ]")

	var res models.Response

	nameservice := "[PackageRedeemtion]-[RedeemtionUVServices]"
	logReq := fmt.Sprintf("[AccountNumber : %v, RewardID : %v]", param.AccountNumber, param.RewardID)

	logrus.Info(nameservice)

	dataorder := services.DataParameterOrder(constants.CODE_CONFIG_UV_GROUP, constants.CODE_CONFIG_UV_NAME, constants.CODE_CONFIG_UV_EMAIL, constants.CODE_CONFIG_UV_PHONE, constants.CODE_CONFIG_UV_EXPIRED)
	param.CumReffnum = utils.GenTransactionId()

	timeExp, _ := strconv.Atoi(dataorder.Expired)
	exp := utils.FormatTimeString(time.Now(), 0, 0, timeExp)
	param.ExpDate = exp

	// total := strconv.Itoa(req.Jumlah)

	// spending point and spending usage_limit voucher
	textCommentSpending := param.CumReffnum + "#" + param.NamaVoucher
	param.Comment = textCommentSpending

	RedeemVouchUV, errRedeemVouchUV := Trx.V2_Redeem_PointandVoucher(req.Jumlah, param)

	param.PointTransferID = RedeemVouchUV.PointTransferID
	logrus.Info("[ Result Spending point / Deduct point ]")
	logrus.Info(RedeemVouchUV)

	if RedeemVouchUV.Rd == "Invalid JWT Token" {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[V2_Redeem_PointandVoucher]-[Response : %v]", RedeemVouchUV))
		logrus.Println("[ ResponseCode ] : ", RedeemVouchUV.Rc)
		logrus.Println("[ ResponseDesc ] : ", RedeemVouchUV.Rd)
		logrus.Println(logReq)

		res = models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: models.UltraVoucherResp{
				Code:    "60",
				Msg:     "Token or Session Expired Please Login Again",
				Success: 0,
				Failed:  req.Jumlah,
				Pending: 0,
			},
		}

		return res
	}

	if RedeemVouchUV.Rd == "not enough points" {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[V2_Redeem_PointandVoucher]-[Response : %v]", RedeemVouchUV))
		logrus.Println("[ ResponseCode ] : ", RedeemVouchUV.Rc)
		logrus.Println("[ ResponseDesc ] : ", RedeemVouchUV.Rd)
		logrus.Println(logReq)

		// res = utils.GetMessageResponse(res, 500, false, errors.New("Point Tidak Cukup"))
		res = models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: models.UltraVoucherResp{
				Code:    "27",
				Msg:     "Point Tidak Mencukupi",
				Success: 0,
				Failed:  req.Jumlah,
				Pending: 0,
			},
		}

		return res
	}

	if RedeemVouchUV.Rc == "208" {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[V2_Redeem_PointandVoucher]-[Response : %v]", RedeemVouchUV))
		logrus.Println("[ ResponseCode ] : ", RedeemVouchUV.Rc)
		logrus.Println("[ ResponseDesc ] : ", RedeemVouchUV.Rd)
		logrus.Println(logReq)

		res = models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: models.UltraVoucherResp{
				Code:    "65",
				Msg:     "Voucher not available",
				Success: 0,
				Failed:  req.Jumlah,
				Pending: 0,
			},
		}

		return res
	}

	if RedeemVouchUV.Rc == "209" {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[V2_Redeem_PointandVoucher]-[Response : %v]", RedeemVouchUV))
		logrus.Println("[ ResponseCode ] : ", RedeemVouchUV.Rc)
		logrus.Println("[ ResponseDesc ] : ", RedeemVouchUV.Rd)
		logrus.Println(logReq)

		res = models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: models.UltraVoucherResp{
				Code:    "66",
				Msg:     "Payment count limit exceeded",
				Success: 0,
				Failed:  req.Jumlah,
				Pending: 0,
			},
		}

		return res
	}

	var c string
	for _, vall := range RedeemVouchUV.CouponseVouch {
		c = vall.CouponsCode
	}
	fmt.Println("Value CouponCode : ", c)

	if errRedeemVouchUV != nil || RedeemVouchUV.Rc != "00" {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[V2_Redeem_PointandVoucher]-[Error : %v]", errRedeemVouchUV))
		logrus.Println("[ ResponseCode ] : ", RedeemVouchUV.Rc)
		logrus.Println("[ ResponseDesc ] : ", RedeemVouchUV.Rd)
		logrus.Println(logReq)

		// res = utils.GetMessageResponse(res, 500, false, errors.New("Gagal! Maaf transaksi Anda tidak dapat dilakukan saat ini. Silahkan dicoba lagi atau hubungi tim kami untuk informasi selengkapnya."))
		res = models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: models.UltraVoucherResp{
				Code:    "01",
				Msg:     "Gagal! Maaf transaksi Anda tidak dapat dilakukan saat ini. Silahkan dicoba lagi atau hubungi tim kami untuk informasi selengkapnya.",
				Success: 0,
				Failed:  req.Jumlah,
				Pending: 0,
			},
		}

		return res
	}

	expired, _ := strconv.Atoi(dataorder.Expired)
	reqOrder := uvmodels.OrderVoucherReq{
		Sku:               param.ProductCode,
		Qty:               req.Jumlah,
		AccountID:         param.AccountId,
		InstitutionRefno:  param.CumReffnum,
		ExpireDateVoucher: expired,
		ReceiverName:      constants.RECEIVER_NAME_UV,
		ReceiverEmail:     dataorder.Email,
		ReceiverPhone:     dataorder.Phone,
	}

	// order to u
	fmt.Println("[ >>>>>>>>>>>>> OrderVoucher UV <<<<<<<<<<< ]")
	order, errOrder := uv.OrderVoucher(reqOrder, param.InstitutionID)

	param.DataSupplier.Rd = order.ResponseDesc
	param.DataSupplier.Rc = order.ResponseCode

	// reffNumberUV
	param.RRN = order.Data.InvoiceUV

	if errOrder != nil || order.ResponseCode == "" {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[OrderVoucher]-[Error : %v]", errOrder))
		logrus.Error(fmt.Sprintf("[Response %v]", order))
		logrus.Println(logReq)

		for i := req.Jumlah; i > 0; i-- {

			// TrxId
			param.TrxID = utils.GenTransactionId()

			fmt.Println(fmt.Sprintf("[Line Save DB : %v]", i))

			t := i - 1
			coupon := RedeemVouchUV.CouponseVouch[t].CouponsID
			param.CouponID = coupon

			// go SaveTransactionUV(param, order, reqOrder, req, "Reedemtion", "09")
			go services.SaveTransactionUV(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "09", timeExp)
		}

		res = models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: models.UltraVoucherResp{
				// Code: "178",
				Code: "68",
				Msg:  "Transaksi Anda sedang dalam proses. Silahkan hubungi customer support kami untuk informasi selengkapnya.",
				// Msg:     "Maaf koneksi timeout. Silahkan dicoba kembali beberapa saat lagi",
				Success: 0,
				Failed:  0,
				Pending: req.Jumlah,
			},
		}
		return res
	}

	if order.ResponseCode == "02" {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[OrderVoucher]-[Response : %v]", order))
		logrus.Println("[ ResponseCode ] : ", order.ResponseCode)
		logrus.Println("[ ResponseDesc ] : ", order.ResponseDesc)
		logrus.Println(logReq)

		totalPoint := param.Point * req.Jumlah
		param.TrxID = utils.GenTransactionId()
		resultReversal := Trx.V2_Adding_PointVoucher(param, totalPoint, req.Jumlah)
		fmt.Println(resultReversal)

		fmt.Println("========== Send Publisher ==========")

		pubreq := models.NotifPubreq{
			Type:           constants.CODE_REVERSAL_POINT,
			NotificationTo: param.AccountNumber,
			Institution:    param.InstitutionID,
			ReferenceId:    param.RRN,
			TransactionId:  param.Reffnum,
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

		kafkaRes, errKafka := kafka.SendPublishKafka(kafkaReq)
		if errKafka != nil {

			logrus.Error(nameservice)
			logrus.Error(fmt.Sprintf("[SendPublishKafka]-[Error : %v]", errKafka))
			logrus.Println(logReq)

		}
		fmt.Println("Response Publisher : ", kafkaRes)

		for i := req.Jumlah; i > 0; i-- {
			fmt.Println(fmt.Sprintf("[Line Save DB : %v]", i))

			// TrxID
			param.TrxID = utils.GenTransactionId()
			t := i - 1
			coupon := RedeemVouchUV.CouponseVouch[t].CouponsID
			param.CouponID = coupon
			go services.SaveTransactionUV(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "01", timeExp)

		}

		fmt.Println("[ >>>>>>>>>>>>> Response Redeemtion Ultra Voucher UV <<<<<<<<<<<<<<<< ]")
		fmt.Println("[ Code ] : ", "176")
		fmt.Println("[ Coummulatif Reff Num ] : ", param.CumReffnum)
		fmt.Println("[ Order ] : ", req.Jumlah)
		fmt.Println("[ Success ] : ", 0)
		fmt.Println("[ Failed ] : ", req.Jumlah)
		fmt.Println("[ Pending ] : ", 0)

		res = models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: models.UltraVoucherResp{
				Code:    "176",
				Msg:     fmt.Sprintf("Voucher yg tersedia %v", order.Data.VouchersAvailable),
				Success: 0,
				Failed:  0,
				Pending: req.Jumlah,
			},
		}

		return res

	}

	if order.ResponseCode != "00" {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[OrderVoucher]-[Response : %v]", order))
		logrus.Println("[ ResponseCode ] : ", order.ResponseCode)
		logrus.Println("[ ResponseDesc ] : ", order.ResponseDesc)
		logrus.Println(logReq)

		// TrxID
		param.TrxID = utils.GenTransactionId()
		totalPoint := param.Point * req.Jumlah
		resultReversal := Trx.V2_Adding_PointVoucher(param, totalPoint, req.Jumlah)
		fmt.Println(resultReversal)

		fmt.Println("========== Send Publisher ==========")

		pubreq := models.NotifPubreq{
			Type:           constants.CODE_REVERSAL_POINT,
			NotificationTo: param.AccountNumber,
			Institution:    param.InstitutionID,
			ReferenceId:    param.RRN,
			TransactionId:  param.Reffnum,
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

		kafkaRes, errKafka := kafka.SendPublishKafka(kafkaReq)
		if errKafka != nil {

			logrus.Error(nameservice)
			logrus.Error(fmt.Sprintf("[SendPublishKafka]-[Error : %v]", errKafka))
			logrus.Println(logReq)

		}
		fmt.Println("Response Publisher : ", kafkaRes)

		for i := req.Jumlah; i > 0; i-- {

			// TrxID
			param.TrxID = utils.GenTransactionId()
			t := i - 1
			coupon := RedeemVouchUV.CouponseVouch[t].CouponsID
			param.CouponID = coupon
			go services.SaveTransactionUV(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "01", timeExp)

		}

		fmt.Println("[ >>>>>>>>>>>>> Response Redeemtion Ultra Voucher UV <<<<<<<<<<<<<<<< ]")
		logrus.Info("[ Code ] : ", "176")
		logrus.Info("[ Coummulatif Reff Num ] : ", param.CumReffnum)
		logrus.Info("[ Order ] : ", req.Jumlah)
		logrus.Info("[ Success ] : ", 0)
		logrus.Info("[ Failed ] : ", 0)
		logrus.Info("[ Pending ] : ", req.Jumlah)

		res = models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: models.UltraVoucherResp{
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
		fmt.Println(fmt.Sprintf("[Line Save DB : %v]", i))

		// TrxID
		param.TrxID = utils.GenTransactionId()
		t := i - 1
		coupon := RedeemVouchUV.CouponseVouch[t].CouponsID
		param.CouponID = coupon
		code := order.Data.VouchersCode[t].Code

		id := utils.GenerateTokenUUID()
		go services.SaveDB(id, param.InstitutionID, param.CouponID, code, param.AccountNumber, param.AccountId, req.CampaignID)
		go services.SaveTransactionUV(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "00", timeExp)

	}

	fmt.Println("[ >>>>>>>>>>>>> Response Redeemtion Ultra Voucher UV <<<<<<<<<<<<<<<< ]")
	logrus.Info("[ Code ] : ", "00")
	logrus.Info("[ Coummulatif Reff Num ] : ", param.CumReffnum)
	logrus.Info("[ Order ] : ", req.Jumlah)
	logrus.Info("[ Success ] : ", req.Jumlah)
	logrus.Info("[ Failed ] : ", 0)
	logrus.Info("[ Pending ] : ", 0)

	fmt.Println("Response UV : ", order)
	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: models.UltraVoucherResp{
			Code:    "00",
			Msg:     "Success",
			Success: req.Jumlah,
			Failed:  0,
			Pending: 0,
		},
	}

	return res
}
