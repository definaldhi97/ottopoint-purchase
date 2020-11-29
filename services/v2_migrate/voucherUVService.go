package v2_migrate

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	kafka "ottopoint-purchase/hosts/publisher/host"
	uv "ottopoint-purchase/hosts/ultra_voucher/host"
	uvmodels "ottopoint-purchase/hosts/ultra_voucher/models"
	"ottopoint-purchase/models"
	"ottopoint-purchase/services"
	"ottopoint-purchase/utils"

	"github.com/astaxie/beego/logs"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/vjeantet/jodaTime"
	"go.uber.org/zap"
)

type VoucherUVMigrateService struct {
	General models.GeneralModel
}

func (t VoucherUVMigrateService) VoucherUV(req models.VoucherComultaiveReq, param models.Params, header models.RequestHeader) models.Response {
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
	textCommentSpending := param.CumReffnum + param.InstitutionID + constants.CodeSchedulerSpending + "#" + param.NamaVoucher
	param.Comment = textCommentSpending
	RedeemVouchUV, errRedeemVouchUV := Redeem_PointandVoucher(req.Jumlah, param, param.CumReffnum, header)
	logrus.Info("Response Spending point / Deduct point")
	logrus.Info(RedeemVouchUV)

	if RedeemVouchUV.Rc == "10" || RedeemVouchUV.Rd == "Insufficient Point" {

		logrus.Info("[UltraVoucherServices]-[RedeemVoucher]")
		logrus.Info("[Not enough points]-[Gagal Redeem Voucher]")
		logrus.Info("[Rc] : ", RedeemVouchUV.Rc)
		logrus.Info("[Rd] : ", RedeemVouchUV.Rd)

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

	if RedeemVouchUV.Rc == "208" || RedeemVouchUV.Rd == "Voucher not available" {

		logrus.Info("[UltraVoucherServices]-[RedeemVoucher]")
		logrus.Info("[Voucher not available]-[Gagal Redeem Voucher]")
		logrus.Info("[Rc] : ", RedeemVouchUV.Rc)
		logrus.Info("[Rd] : ", RedeemVouchUV.Rd)

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
		logrus.Info("[UltraVoucherServices]-[RedeemVoucher]")
		logrus.Info("[Rc] : ", RedeemVouchUV.Rc)
		logrus.Info("[Rd] : ", RedeemVouchUV.Rd)

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
	fmt.Println("[ >>>>>>>>>>>>>>>>>>>>>>>>>>> OrderVoucher UV <<<<<<<<<<<<<<<<<<<<< ]")
	order, errOrder := uv.OrderVoucher(reqOrder, param.InstitutionID)

	param.DataSupplier.Rd = order.ResponseDesc
	param.DataSupplier.Rc = order.ResponseCode

	// reffNumberUV
	param.RRN = order.Data.InvoiceUV

	if errOrder != nil || order.ResponseCode == "" {

		logrus.Info("[UltraVoucherServices]-[OrderVoucher]")
		logrus.Info("[Failed Order Voucher]-[Gagal Order Voucher]")
		logrus.Error("Error oreder UV : ", errOrder)
		logrus.Info("Response OrderVoucher : ", order)

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

		logrus.Info("[UltraVoucherServices]-[OrderVoucher]")
		logrus.Info("[Stock not Available]-[Gagal Order Voucher]")
		logrus.Info("[ Response OrderVoucher ] : ", order)
		logrus.Info("[ ResponseCode ] : ", order.ResponseCode)

		totalPoint := param.Point * req.Jumlah
		param.TrxID = utils.GenTransactionId()
		resultReversal := Adding_PointVoucher(param, totalPoint, req.Jumlah, param.CumReffnum, header)
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
			logrus.Error("Gagal Send Publisher : ", err)
		}
		logrus.Info("[ Response Publisher ] : ", kafkaRes)

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
		logrus.Info("[ Code ] : ", "176")
		logrus.Info("[ Coummulatif Reff Num ] : ", param.CumReffnum)
		logrus.Info("[ Order ] : ", req.Jumlah)
		logrus.Info("[ Success ] : ", 0)
		logrus.Info("[ Failed ] : ", 0)
		logrus.Info("[ Pending ] : ", req.Jumlah)

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

		logrus.Info("[UltraVoucherServices]-[OrderVoucher]")
		logrus.Info("[Stock not Available]-[Gagal Order Voucher]")

		logrus.Info("[ ResponseCode ] : ", order.ResponseCode)
		logrus.Info("[ ResponseDecs ] : ", order.ResponseDesc)

		// sugarLogger.Info("Internal Server Error : ", errOrder)

		// TrxID
		param.TrxID = utils.GenTransactionId()
		totalPoint := param.Point * req.Jumlah
		resultReversal := Adding_PointVoucher(param, totalPoint, req.Jumlah, param.CumReffnum, header)
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

			logrus.Error("Gagal Send Publisher : ", err)
		}
		logrus.Info("[ Response Publisher ] : ", kafkaRes)

		for i := req.Jumlah; i > 0; i-- {

			// TrxID
			param.TrxID = utils.GenTransactionId()
			t := i - 1
			coupon := RedeemVouchUV.CouponseVouch[t].CouponsID
			param.CouponID = coupon
			go services.SaveTransactionUV(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "01")

		}

		fmt.Println("[ >>>>>>>>>>>>> Response Redeemtion Ultra Voucher UV <<<<<<<<<<<<<<<< ]")
		logrus.Info("[ Code ] : ", "01")
		logrus.Info("[ Coummulatif Reff Num ] : ", param.CumReffnum)
		logrus.Info("[ Order ] : ", req.Jumlah)
		logrus.Info("[ Success ] : ", 0)
		logrus.Info("[ Failed ] : ", req.Jumlah)
		logrus.Info("[ Pending ] : ", 0)

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
	logrus.Info("[ Code ] : ", "00")
	logrus.Info("[ Coummulatif Reff Num ] : ", param.CumReffnum)
	logrus.Info("[ Order ] : ", req.Jumlah)
	logrus.Info("[ Success ] : ", req.Jumlah)
	logrus.Info("[ Failed ] : ", 0)
	logrus.Info("[ Pending ] : ", 0)

	logrus.Info("Response UV : ", order)
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

func (t VoucherUVMigrateService) CallbackVoucherUV(req models.UseVoucherUVReq, param models.Params, campaignID string) models.Response {
	var res models.Response

	fmt.Println("[ >>>>>>>>>>>>>>>>>>>>> Callbakc Voucher UV Service <<<<<<<<<<<<<<<<<< ]")

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[UseVoucherUV]",
		zap.String("AccountNumber : ", param.AccountNumber), zap.String("InstitutionID : ", param.InstitutionID),
		zap.String("category : ", param.Category), zap.String("campaignId : ", campaignID),
		zap.String("AccountID : ", param.AccountId), zap.String("AccountNumber : ", param.AccountNumber),
		zap.String("VoucherCode : ", req.VoucherCode))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[UseVoucherUV]")
	defer span.Finish()

	logrus.Info("Campaign : ", campaignID)
	logrus.Info("CouponID : ", param.CouponID)
	logrus.Info("ProductCode : ", param.CouponCode)
	logrus.Info("AccountID : ", param.AccountId)

	// // Use Voucher to Openloyalty
	// use, err2 := opl.CouponVoucherCustomer(campaignID, param.CouponID, param.CouponCode, param.AccountId, 1)

	// var useErr string
	// for _, value := range use.Coupons {
	// 	useErr = value.CouponID
	// }

	// if err2 != nil || useErr == "" {

	// 	logs.Info(fmt.Sprintf("[Error : %v]", err2))
	// 	logs.Info(fmt.Sprintf("[Response : %v]", use))
	// 	logs.Info("[Error from OPL]-[CouponVoucherCustomer]")

	// 	// go SaveTransactionUV(param, useUV, reqUV, req, "Used", "01", "")

	// 	res = utils.GetMessageResponse(res, 400, false, errors.New("Gagal Use Voucher, Harap coba lagi"))
	// 	return res
	// }

	timeUse := jodaTime.Format("dd-MM-YYYY HH:mm:ss", time.Now())

	_, errUpdate := db.UpdateVoucher(timeUse, param.CouponID)
	if errUpdate != nil {

		logs.Info(fmt.Sprintf("[Error : %v]", errUpdate))
		logs.Info("[Gagal Update Voucher]")
		logs.Info("[UseVoucherUV]-[Package-Services]")

	}

	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: models.UseVoucherUVResp{
			Voucher: param.NamaVoucher,
		},
	}
	return res

}
