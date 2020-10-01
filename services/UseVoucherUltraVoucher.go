package services

import (
	"encoding/json"
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	"ottopoint-purchase/hosts/opl/host"
	opl "ottopoint-purchase/hosts/opl/host"
	kafka "ottopoint-purchase/hosts/publisher/host"
	uv "ottopoint-purchase/hosts/ultra_voucher/host"
	uvmodels "ottopoint-purchase/hosts/ultra_voucher/models"
	"ottopoint-purchase/models"
	"ottopoint-purchase/models/dbmodels"
	"ottopoint-purchase/utils"
	"strconv"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/vjeantet/jodaTime"
	"go.uber.org/zap"
)

type UseVoucherUltraVoucher struct {
	General models.GeneralModel
}

func (t UseVoucherUltraVoucher) UltraVoucherServices(req models.VoucherComultaiveReq, param models.Params) models.Response {
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

	// var err bool
	// success := 0
	// failed := 0
	// couponOPL := []models.CouponsRedeem{}

	// for i := req.Jumlah; i >= 1; i-- {

	dataorder := DataParameterOrder()

	// param.Reffnum = utils.GenTransactionId()

	param.CumReffnum = utils.GenTransactionId()

	timeExp, _ := strconv.Atoi(dataorder.Expired)

	exp := utils.FormatTimeString(time.Now(), 0, 0, timeExp)

	param.ExpDate = exp

	total := strconv.Itoa(req.Jumlah)

	param.Amount = int64(param.Point)

	// redeem to opl (potong point)
	redeem, errredeem := host.RedeemVoucherCumulative(req.CampaignID, param.AccountId, total, "0")

	if redeem.Message == "Invalid JWT Token" {
		fmt.Println("Error : ", errredeem)
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

	if redeem.Error == "Not enough points" {
		fmt.Println("Error : ", errredeem)
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

	if redeem.Error == "Limit exceeded" {
		fmt.Println("Error : ", errredeem)
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
	for _, vall := range redeem.Coupons {
		c = vall.Code
	}

	if errredeem != nil || redeem.Error != "" || c == "" {
		fmt.Println("Error : ", errredeem)
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

	nama := "OTTOPOINT"
	expired, _ := strconv.Atoi(dataorder.Expired)
	reqOrder := uvmodels.OrderVoucherReq{
		Sku:               param.ProductCode,
		Qty:               req.Jumlah,
		AccountID:         param.AccountId,
		InstitutionRefno:  param.CumReffnum,
		ExpireDateVoucher: expired,
		ReceiverName:      nama,
		ReceiverEmail:     dataorder.Email,
		ReceiverPhone:     dataorder.Phone,
	}

	// order to u
	fmt.Println(">>> OrderVoucher <<<")
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
			coupon := redeem.Coupons[t].Id
			param.CouponID = coupon

			// go SaveTransactionUV(param, order, reqOrder, req, "Reedemtion", "09")
			go SaveTransactionUV(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "09")
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

		for i := req.Jumlah; i > 0; i-- {
			fmt.Println(fmt.Sprintf("[Line : %v]", i))

			t := i - 1
			CouponID := redeem.Coupons[t].Id
			couponCode := redeem.Coupons[t].Code

			fmt.Println(fmt.Sprintf("[Reversal Voucher %v]", param.NamaVoucher))
			use, err2 := opl.CouponVoucherCustomer(req.CampaignID, CouponID, couponCode, param.AccountId, 1)

			var useErr string
			for _, value := range use.Coupons {
				useErr = value.CouponID
			}

			if err2 != nil || useErr == "" {
				// res = utils.GetMessageResponse(res, 400, false, errors.New("Gagal Redeem Voucher, Harap coba lagi"))
				// return res

				fmt.Println("[UltraVoucherServices]-[CouponVoucherCustomer]")
				fmt.Println(fmt.Sprintf("[UltraVoucherServices]-[Error : %v]", err2))
				sugarLogger.Info("[UltraVoucherServices]-[CouponVoucherCustomer]")
			}
		}

		// TrxID
		param.TrxID = utils.GenTransactionId()

		text := param.TrxID + param.InstitutionID + constants.CodeReversal + "#" + "OP009 - Reversal point cause transaction " + param.NamaVoucher + " is failed"
		// Text := "OP009 - " + "Reversal point cause transaction " + param.NamaVoucher + " is failed"
		point := param.Point * req.Jumlah
		totalPoint := strconv.Itoa(point)
		param.TrxID = utils.GenTransactionId()
		sendReversal, errReversal := host.TransferPoint(param.AccountId, totalPoint, text)
		if errReversal != nil {
			fmt.Println("Internal Server Error : ", errReversal)
			fmt.Println("[UltraVoucherServices]-[TransferPoint]")
			fmt.Println("[Failed Order Voucher]-[Gagal Reversal Point]")
		}

		expired := ExpiredPointService()

		saveReversal := dbmodels.TEarning{
			ID: utils.GenerateTokenUUID(),
			// EarningRule     :,
			// EarningRuleAdd  :,
			PartnerId: param.InstitutionID,
			// ReferenceId     : ,
			TransactionId: param.TrxID,
			// ProductCode     :,
			// ProductName     :,
			AccountNumber: param.AccountNumber,
			// Amount          :,
			Point: int64(point),
			// Remark          :,
			Status:           constants.Success,
			StatusMessage:    "Success",
			PointsTransferId: sendReversal.PointsTransferId,
			// RequestorData   :,
			// ResponderData   :,
			TransType:       constants.CodeReversal,
			AccountId:       param.AccountId,
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
			NotificationTo: param.AccountNumber,
			Institution:    param.InstitutionID,
			ReferenceId:    param.RRN,
			TransactionId:  param.Reffnum,
			Data: models.DataValue{
				RewardValue: "point",
				Value:       totalPoint,
			},
		}
		// pubreq := models.NotifPubreq{
		// 	Type:           constants.CODE_REVERSAL_POINT,
		// 	NotificationTo: param.AccountNumber,
		// 	Institution:    param.InstitutionID,
		// 	ReferenceId:    param.RRN,
		// 	TransactionId:  param.Reffnum,
		// 	Data: models.DataValue{
		// 		RewardValue: param.NamaVoucher,
		// 		Value:       strconv.Itoa(point),
		// 	},
		// }

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
			coupon := redeem.Coupons[t].Id
			param.CouponID = coupon

			go SaveTransactionUV(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "01")
		}

		// res = utils.GetMessageResponse(res, 145, false, errors.New(fmt.Sprintf("Voucher yg tersedia %v", order.Data.VouchersAvailable)))
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
		// res.Data = "Stok Tidak Tersedia"

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

		for i := req.Jumlah; i > 0; i-- {
			fmt.Println(fmt.Sprintf("[Line : %v]", i))

			t := i - 1
			CouponID := redeem.Coupons[t].Id
			couponCode := redeem.Coupons[t].Code

			fmt.Println(fmt.Sprintf("[Reversal Voucher %v]", param.NamaVoucher))
			use, err2 := opl.CouponVoucherCustomer(req.CampaignID, CouponID, couponCode, param.AccountId, 1)
			var useErr string
			for _, value := range use.Coupons {
				useErr = value.CouponID
			}

			if err2 != nil || useErr == "" {

				fmt.Println("[UltraVoucherServices]-[CouponVoucherCustomer]")
				fmt.Println(fmt.Sprintf("[UltraVoucherServices]-[Error : %v]", err2))
				sugarLogger.Info("[UltraVoucherServices]-[CouponVoucherCustomer]")
			}
		}

		// TrxID
		param.TrxID = utils.GenTransactionId()

		text := param.TrxID + param.InstitutionID + constants.CodeReversal + "#" + "OP009 - Reversal point cause transaction " + param.NamaVoucher + " is failed"

		// Text := "OP009 - " + "Reversal point cause transaction " + param.NamaVoucher + " is failed"
		point := param.Point * req.Jumlah
		totalPoint := strconv.Itoa(point)
		reversal, errReversal := host.TransferPoint(param.AccountId, totalPoint, text)
		if errReversal != nil || reversal.PointsTransferId == "" {
			fmt.Println("Internal Server Error : ", errReversal)
			fmt.Println("[UltraVoucherServices]-[TransferPoint]")
			fmt.Println("[Failed Order Voucher]-[Gagal Reversal Point]")
		}

		expired := ExpiredPointService()

		saveReversal := dbmodels.TEarning{
			ID: utils.GenerateTokenUUID(),
			// EarningRule     :,
			// EarningRuleAdd  :,
			PartnerId: param.InstitutionID,
			// ReferenceId     : ,
			TransactionId: param.TrxID,
			// ProductCode     :,
			// ProductName     :,
			AccountNumber: param.AccountNumber,
			// Amount          :,
			Point: int64(point),
			// Remark          :,
			Status:           constants.Success,
			StatusMessage:    "Success",
			PointsTransferId: reversal.PointsTransferId,
			// RequestorData   :,
			// ResponderData   :,
			TransType:       constants.CodeReversal,
			AccountId:       param.AccountId,
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
			NotificationTo: param.AccountNumber,
			Institution:    param.InstitutionID,
			ReferenceId:    param.RRN,
			TransactionId:  param.Reffnum,
			Data: models.DataValue{
				RewardValue: "point",
				Value:       totalPoint,
			},
		}
		// pubreq := models.NotifPubreq{
		// 	Type:           constants.CODE_REVERSAL_POINT,
		// 	NotificationTo: param.AccountNumber,
		// 	Institution:    param.InstitutionID,
		// 	ReferenceId:    param.RRN,
		// 	TransactionId:  param.Reffnum,
		// 	Data: models.DataValue{
		// 		RewardValue: param.NamaVoucher,
		// 		Value:       strconv.Itoa(point),
		// 	},
		// }

		bytePub, _ := json.Marshal(pubreq)

		kafkaReq := kafka.PublishReq{
			Topic: "ottopoint-notification-reversal",
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
			coupon := redeem.Coupons[t].Id
			param.CouponID = coupon

			go SaveTransactionUV(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "01")
		}

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

	for i := req.Jumlah; i > 0; i-- {

		fmt.Println(fmt.Sprintf("[Line Save DB : %v]", i))

		// TrxID
		param.TrxID = utils.GenTransactionId()

		t := i - 1
		coupon := redeem.Coupons[t].Id
		code := order.Data.VouchersCode[t].Code

		param.CouponID = coupon

		id := utils.GenerateTokenUUID()
		go SaveDB(id, param.InstitutionID, coupon, code, param.AccountNumber, param.AccountId, req.CampaignID)

		go SaveTransactionUV(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "00")
	}

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

func SaveDB(id, institution, coupon, vouchercode, phone, custIdOPL, campaignID string) {
	fmt.Println("[SaveDB]-[UltraVoucherServices]")
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
		fmt.Println("[Failed Save to DB ]", err)
		fmt.Println("[Package-Services]-[UltraVoucherServices]")
		// return err
	}
}

func DataParameterOrder() models.ParamUV {
	res := models.ParamUV{}

	nama := "" // nama
	email := "UV_EMAIL_ORDER"
	phone := "UV_PHONE_ORDER"
	expired := "UV_EXPIRED_VOUCHER"

	datanama, errnama := db.ParamData(nama)
	if errnama != nil {
		fmt.Println("[Error get data from Db m_paramaters]")
		fmt.Println("Error : ", errnama)
		fmt.Println("Code :", nama)
	}

	dataemail, erremail := db.ParamData(email)
	if erremail != nil {
		fmt.Println("[Error get data from Db m_paramaters]")
		fmt.Println("Error : ", erremail)
		fmt.Println("Code :", email)
	}

	dataphone, errphone := db.ParamData(phone)
	if errphone != nil {
		fmt.Println("[Error get data from Db m_paramaters]")
		fmt.Println("Error : ", errphone)
		fmt.Println("Code :", phone)
	}

	dataexpired, errexpired := db.ParamData(expired)
	if errexpired != nil {
		fmt.Println("[Error get data from Db m_paramaters]")
		fmt.Println("Error : ", errexpired)
		fmt.Println("Code :", expired)
	}

	res = models.ParamUV{
		Nama:    datanama.Value,
		Email:   dataemail.Value,
		Phone:   dataphone.Value,
		Expired: dataexpired.Value,
	}

	return res
}

func chekStatus(institutionId, reffNo string) (uvmodels.OrderVoucherResp, error) {
	res := uvmodels.OrderVoucherResp{}
	var err error
	var no int

	for i := 3; i > 0; i-- {
		no++

		fmt.Println(fmt.Sprintf("[Percobaan ke : %v]", no))
		time.Sleep(5 * time.Second)

		res, err = uv.CheckStatusOrder(institutionId, reffNo)
		fmt.Println(fmt.Sprintf("[Response ke %v : %v]", no, res))

		if res.ResponseCode != "" {
			return res, nil
		}

	}

	return res, err

}
