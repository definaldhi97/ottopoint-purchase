package v2_migrate

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
	"ottopoint-purchase/utils"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type VoucherUVMigrateService struct {
	General models.GeneralModel
}

func (t VoucherUVMigrateService) VoucherUV(req models.VoucherComultaiveReq, param models.Params) models.Response {
	fmt.Println("[ >>>>>>>>>>>>>>>>>> Voucher UV Service <<<<<<<<<<<<<<<< ]")

	var res models.Response

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[UltraVoucherServices]",
		zap.String("NameVoucher : ", param.NamaVoucher), zap.Int("Jumlah : ", req.Jumlah),
		zap.String("CampaignID : ", req.CampaignID), zap.String("CampaignID : ", req.CampaignID),
		zap.String("CustID2 : ", req.CustID2), zap.String("ProductCode : ", param.ProductCode),
		// zap.Int("Point : ", req.Point),
		zap.String("AccountNumber : ", param.AccountNumber), zap.String("InstitutionID : ", param.InstitutionID))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[UltraVoucherServices]")
	defer span.Finish()

	dataorder := services.DataParameterOrder()
	param.CumReffnum = utils.GenTransactionId()

	timeExp, _ := strconv.Atoi(dataorder.Expired)
	exp := utils.FormatTimeString(time.Now(), 0, 0, timeExp)
	param.ExpDate = exp

	// total := strconv.Itoa(req.Jumlah)

	// spending point and spending usage_limit voucher
	textCommentSpending := param.CumReffnum + "#" + param.NamaVoucher
	param.Comment = textCommentSpending
	RedeemVouchUV, errRedeemVouchUV := Redeem_PointandVoucher(req.Jumlah, param)
	fmt.Println("result Spending point / Deduct point")
	fmt.Println(RedeemVouchUV)

	if RedeemVouchUV.Rd == "Invalid JWT Token" {
		fmt.Println("Error : ", errRedeemVouchUV)
		fmt.Println("[UltraVoucherServices]-[RedeemVoucher]")
		fmt.Println("[Internal Server Error]-[Gagal Redeem Voucher]")

		// sugarLogger.Info("Internal Server Error : ", errredeem)
		sugarLogger.Info("[UltraVoucherServices]-[RedeemVoucher]")
		sugarLogger.Info("[Internal Server Error]-[Gagal Redeem Voucher]")

		// res = utils.GetMessageResponse(res, 422, false, errors.New("Gagal! Maaf transaksi Anda tidak dapat dilakukan saat ini. Silahkan dicoba lagi atau hubungi tim kami untuk informasi selengkapnya."))

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
		fmt.Println("Error : ", errRedeemVouchUV)
		fmt.Println("[UltraVoucherServices]-[RedeemVoucher]")
		fmt.Println("[Not enough points]-[Gagal Redeem Voucher]")

		// sugarLogger.Info("Internal Server Error : ", errredeem)
		sugarLogger.Info("[UltraVoucherServices]-[RedeemVoucher]")
		sugarLogger.Info("[Not enough points]-[Gagal Redeem Voucher]")

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

	if RedeemVouchUV.Rd == "Voucher not available" {
		fmt.Println("Error : ", errRedeemVouchUV)
		fmt.Println("[UltraVoucherServices]-[RedeemVoucher]")
		fmt.Println("[Limit exceeded]-[Gagal Redeem Voucher]")

		// sugarLogger.Info("Internal Server Error : ", errredeem)
		sugarLogger.Info("[UltraVoucherServices]-[RedeemVoucher]")
		sugarLogger.Info("[Limit exceeded]-[Gagal Redeem Voucher]")

		// res = utils.GetMessageResponse(res, 500, false, errors.New("Voucher Sudah Limit"))

		res = models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: models.UltraVoucherResp{
				Code:    "65",
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

	if errRedeemVouchUV != nil || RedeemVouchUV.Rc != "00" || c == "" {
		fmt.Println("Error : ", errRedeemVouchUV)
		fmt.Println("[UltraVoucherServices]-[RedeemVoucher]")
		fmt.Println("[Failed Redeem Voucher]-[Gagal Redeem Voucher]")

		// sugarLogger.Info("Internal Server Error : ", errredeem)
		sugarLogger.Info("[UltraVoucherServices]-[RedeemVoucher]")
		sugarLogger.Info("[Failed Redeem Voucher]-[Gagal Redeem Voucher]")

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

		fmt.Println("Error : ", errOrder)
		fmt.Println("Response OrderVoucher : ", order)
		fmt.Println("[UltraVoucherServices]-[OrderVoucher]")
		fmt.Println("[Failed Order Voucher]-[Gagal Order Voucher]")

		// sugarLogger.Info("Internal Server Error : ", errOrder)
		sugarLogger.Info("[UltraVoucherServices]-[OrderVoucher]")
		sugarLogger.Info("[Failed Order Voucher]-[Gagal Order Voucher]")

		for i := req.Jumlah; i > 0; i-- {

			// TrxId
			param.TrxID = utils.GenTransactionId()

			fmt.Println(fmt.Sprintf("[Line Save DB : %v]", i))

			t := i - 1
			coupon := RedeemVouchUV.CouponseVouch[t].CouponsID
			param.CouponID = coupon

			// go SaveTransactionUV(param, order, reqOrder, req, "Reedemtion", "09")
			go services.SaveTransactionUV(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "09")
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

		fmt.Println("Internal Server Error : ", errOrder)
		fmt.Println("ResponseCode : ", order.ResponseCode)
		fmt.Println("[UltraVoucherServices]-[OrderVoucher]")
		fmt.Println("[Stock not Available]-[Gagal Order Voucher]")

		// sugarLogger.Info("Internal Server Error : ", errOrder)
		sugarLogger.Info("[UltraVoucherServices]-[OrderVoucher]")
		sugarLogger.Info("[Stock not Available]-[Gagal Order Voucher]")

		totalPoint := param.Point * req.Jumlah
		param.TrxID = utils.GenTransactionId()
		resultReversal := Adding_PointVoucher(param, totalPoint, req.Jumlah)
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
			coupon := RedeemVouchUV.CouponseVouch[t].CouponsID
			param.CouponID = coupon
			go services.SaveTransactionUV(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "01")

		}

		fmt.Println("[ >>>>>>>>>>>>> Response Redeemtion Ultra Voucher UV <<<<<<<<<<<<<<<< ]")
		fmt.Println("[ Code ] : ", "176")
		fmt.Println("[ Coummulatif Reff Num ] : ", param.CumReffnum)
		fmt.Println("[ Order ] : ", req.Jumlah)
		fmt.Println("[ Success ] : ", 0)
		fmt.Println("[ Failed ] : ", 0)
		fmt.Println("[ Pending ] : ", req.Jumlah)

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

		fmt.Println("Internal Server Error : ", errOrder)
		fmt.Println("ResponseCode : ", order.ResponseCode)
		fmt.Println("[UltraVoucherServices]-[OrderVoucher]")
		fmt.Println(fmt.Sprintf("[Response %v]", order.ResponseCode))

		// sugarLogger.Info("Internal Server Error : ", errOrder)
		sugarLogger.Info("[UltraVoucherServices]-[OrderVoucher]")
		sugarLogger.Info("[Stock not Available]-[Gagal Order Voucher]")

		// TrxID
		param.TrxID = utils.GenTransactionId()
		totalPoint := param.Point * req.Jumlah
		resultReversal := Adding_PointVoucher(param, totalPoint, req.Jumlah)
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

		kafkaRes, err := kafka.SendPublishKafka(kafkaReq)
		if err != nil {
			fmt.Println("Gagal Send Publisher")
			fmt.Println("Error : ", err)
		}
		fmt.Println("Response Publisher : ", kafkaRes)

		for i := req.Jumlah; i > 0; i-- {

			// TrxID
			param.TrxID = utils.GenTransactionId()
			t := i - 1
			coupon := RedeemVouchUV.CouponseVouch[t].CouponsID
			param.CouponID = coupon
			go services.SaveTransactionUV(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "01")

		}

		fmt.Println("[ >>>>>>>>>>>>> Response Redeemtion Ultra Voucher UV <<<<<<<<<<<<<<<< ]")
		fmt.Println("[ Code ] : ", "0")
		fmt.Println("[ Coummulatif Reff Num ] : ", param.CumReffnum)
		fmt.Println("[ Order ] : ", req.Jumlah)
		fmt.Println("[ Success ] : ", 0)
		fmt.Println("[ Failed ] : ", req.Jumlah)
		fmt.Println("[ Pending ] : ", 0)

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
		go services.SaveTransactionUV(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "00")

	}

	fmt.Println("[ >>>>>>>>>>>>> Response Redeemtion Ultra Voucher UV <<<<<<<<<<<<<<<< ]")
	fmt.Println("[ Code ] : ", "00")
	fmt.Println("[ Coummulatif Reff Num ] : ", param.CumReffnum)
	fmt.Println("[ Order ] : ", req.Jumlah)
	fmt.Println("[ Success ] : ", req.Jumlah)
	fmt.Println("[ Failed ] : ", 0)
	fmt.Println("[ Pending ] : ", 0)

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