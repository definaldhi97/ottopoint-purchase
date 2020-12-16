package services

import (
	"encoding/json"
	"fmt"
	"ottopoint-purchase/constants"
	db "ottopoint-purchase/db"
	"ottopoint-purchase/models"
	"ottopoint-purchase/models/dbmodels"
	"ottopoint-purchase/utils"
	"reflect"
	"strconv"
	"time"

	"ottopoint-purchase/hosts/opl/host"
	opl "ottopoint-purchase/hosts/opl/host"
	kafka "ottopoint-purchase/hosts/publisher/host"
	signature "ottopoint-purchase/hosts/signature/host"
	vg "ottopoint-purchase/hosts/voucher_aggregator/host"
	vgmodels "ottopoint-purchase/hosts/voucher_aggregator/models"

	ODU "ottodigital.id/library/utils"

	"github.com/astaxie/beego/logs"
	"github.com/opentracing/opentracing-go"
	"github.com/vjeantet/jodaTime"
	"go.uber.org/zap"
)

var (
	callbackHost              = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_CALLBACK_VOUCHERAG", "http://34.101.248.102:8600")
	callbackOttoPointPurchase = ODU.GetEnv("OTTOPOINT_PURCHASE_CALLBACK_VOUCHERAG", "/transaction/v2/callback/redeem/voucherag")
)

type VoucherAgServices struct {
	General models.GeneralModel
}

func (t VoucherAgServices) RedeemVoucher(req models.VoucherComultaiveReq, param models.Params, head models.RequestHeader) models.Response {

	var res models.Response

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[VoucherComulative-Services]",
		zap.String("NameVoucher : ", param.NamaVoucher), zap.Int("Jumlah : ", req.Jumlah),
		zap.String("CampaignID : ", req.CampaignID), zap.String("CustID : ", req.CustID),
		zap.String("CustID2 : ", req.CustID2), zap.String("ProductCode : ", param.ProductCode),
		// zap.Int("Point : ", reiq.Pont),
		zap.String("AccountNumber : ", param.AccountNumber), zap.String("InstitutionID : ", param.InstitutionID))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[RedeemVoucher]")
	defer span.Finish()

	dataOrder := DataParameterOrderVoucherAg()

	param.CumReffnum = utils.GenTransactionId()

	total := strconv.Itoa(req.Jumlah)

	param.Amount = int64(param.Point)

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
			Data: vgmodels.ResponseVoucherAg{
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
			Data: vgmodels.ResponseVoucherAg{
				Code:    "65",
				Msg:     "Payment count limit exceeded",
				Success: 0,
				Failed:  req.Jumlah,
				Pending: 0,
			},
		}

		return res
	}

	if errredeem != nil || redeem.Error != "" {
		fmt.Println("Error : ", errredeem)
		fmt.Println("[VoucherAgService]-[RedeemVoucher]")
		fmt.Println("[Failed Redeem Voucher]-[Gagal Redeem Voucher]")

		// sugarLogger.Info("Internal Server Error : ", errredeem)
		sugarLogger.Info("[VoucherAgService]-[RedeemVoucher]")
		sugarLogger.Info("[Failed Redeem Voucher]-[Gagal Redeem Voucher]")

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
	reqOrder := vgmodels.RequestOrderVoucherAg{
		ProductCode:    param.ProductCode,
		Qty:            req.Jumlah,
		OrderID:        param.CumReffnum,
		CustomerName:   nama,
		CustomerEmail:  dataOrder.Email,
		CustomerPhone:  dataOrder.Phone,
		DeliveryMethod: 1,
		RedeemCallback: callbackHost + callbackOttoPointPurchase,
	}

	fmt.Println("Start - OrderVoucherAggregator")
	sugarLogger.Info("[VoucherAgService]-[OrderVoucher]")

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

	order, errorder := vg.OrderVoucher(reqOrder, head)

	param.DataSupplier.Rd = order.ResponseDesc
	param.DataSupplier.Rc = order.ResponseCode

	if errorder != nil || order.ResponseCode == "" {

		// Reversal Start Here
		fmt.Println("Error : ", errorder)
		fmt.Println("Response OrderVoucher : ", order)
		fmt.Println("[VoucherAgServices]-[OrderVoucher]")
		fmt.Println("[Failed Order Voucher]-[Gagal Order Voucher]")

		sugarLogger.Info("[VoucherAgServices]-[OrderVoucher]")
		sugarLogger.Info("[Failed Order Voucher]-[Gagal Order Voucher]")

		for i := req.Jumlah; i > 0; i-- {

			// TrxId
			param.TrxID = utils.GenTransactionId()

			fmt.Println(fmt.Sprintf("[Line Save DB : %v]", i))

			t := i - 1
			coupon := redeem.Coupons[t].Id
			param.CouponID = coupon

			go SaveTransactionVoucherAg(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "09")
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

		// Start Reversal Here
		fmt.Println("Internal Server Error : ", errorder)
		fmt.Println("ResponseCode : ", order.ResponseCode)
		fmt.Println("[VoucherAgServices]-[OrderVoucher]")
		fmt.Println("[Stock not Available]-[Gagal Order Voucher]")

		sugarLogger.Info("[VoucherAgServices]-[OrderVoucher]")
		sugarLogger.Info("[Stock not Available]-[Gagal Order Voucher]")

		for i := req.Jumlah; i > 0; i-- {
			fmt.Println(fmt.Sprintf("[Line : %v]", i))

			t := i - 1
			coupon := redeem.Coupons[t].Id
			couponCode := redeem.Coupons[t].Code
			use, err2 := opl.CouponVoucherCustomer(req.CampaignID, coupon, couponCode, param.AccountId, 1)

			var useErr string
			for _, value := range use.Coupons {
				useErr = value.CouponID
			}

			if err2 != nil || useErr == "" {
				fmt.Println("[VoucherAgService]-[CouponVoucherCustomer]")
				fmt.Println(fmt.Sprintf("[VoucherAgService]-[Error : %v]", err2))
				sugarLogger.Info("[VoucherAgService]-[CouponVoucherCustomer]")
			}

		}

		// TrxID
		point := param.Point * req.Jumlah
		totalPoint := strconv.Itoa(point)

		param.TrxID = utils.GenTransactionId()
		text := param.TrxID + param.InstitutionID + constants.CodeReversal + "#" + "OP009 - Reversal point cause transaction " + param.NamaVoucher + " is failed"

		schedulerData := dbmodels.TSchedulerRetry{
			Code:          constants.CodeScheduler,
			TransactionID: utils.Before(text, "#"),
			Count:         0,
			IsDone:        false,
			CreatedAT:     time.Now(),
		}

		sendReversal, errReversal := host.TransferPoint(param.AccountId, totalPoint, text)

		statusEarning := constants.Success
		msgEarning := constants.MsgSuccess

		if errReversal != nil || sendReversal.PointsTransferId == "" {

			statusEarning = constants.TimeOut

			fmt.Println(fmt.Sprintf("===== Failed TransferPointOPL to %v || RRN : %v =====", param.AccountNumber, param.RRN))

			statusEarning = constants.TimeOut

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
			Point:            int64(point),
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
			ReferenceId:    param.RRN,
			TransactionId:  param.Reffnum,
			Data: models.DataValue{
				RewardValue: "point",
				Value:       totalPoint,
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
			coupon := redeem.Coupons[t].Id
			param.CouponID = coupon

			go SaveTransactionVoucherAg(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "01")
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

		fmt.Println("Error : ", errorder)
		fmt.Println("Response OrderVoucher : ", order)
		fmt.Println("[VoucherAggregator]-[OrderVoucher]")
		fmt.Println("[Pending Order]-[Gagal Order Voucher]")

		// sugarLogger.Info("Internal Server Error : ", errOrder)
		sugarLogger.Info("[VoucherAggregator]-[OrderVoucher]")
		sugarLogger.Info("[Pending Order]-[Gagal Order Voucher]")

		for i := req.Jumlah; i > 0; i-- {
			fmt.Println(fmt.Sprintf("[Line : %v]", i))

			// TrxId
			param.TrxID = utils.GenTransactionId()

			fmt.Println(fmt.Sprintf("[Line Save DB : %v]", i))

			t := i - 1
			coupon := redeem.Coupons[t].Id
			param.CouponID = coupon

			go SaveTransactionVoucherAg(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "09")

		}

		// Save To Scheduler
		text := param.CumReffnum + param.InstitutionID + constants.CodeReversal + "#" + "OP009 - Reversal point cause transaction " + param.NamaVoucher + " is failed"
		schedulerData := dbmodels.TSchedulerRetry{
			Code:          constants.CodeSchedulerVoucherAG,
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

			// return
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

	// Handle General Error
	if order.ResponseCode != "00" {

		// Reversal Start Here
		fmt.Println("Internal Server Error : ", errorder)
		fmt.Println("ResponseCode : ", order.ResponseCode)
		fmt.Println("[VoucherAgServices]-[FailedOrder]")
		fmt.Println(fmt.Sprintf("[Response %v]", order.ResponseCode))

		sugarLogger.Info("[VoucherAgServices]-[FailedOrder]")
		sugarLogger.Info("[VoucherAgServices]-[Gagal Order Voucher]")

		for i := req.Jumlah; i > 0; i-- {
			fmt.Println(fmt.Sprintf("[Line : %v]", i))

			t := i - 1
			coupon := redeem.Coupons[t].Id
			couponCode := redeem.Coupons[t].Code

			use, err2 := opl.CouponVoucherCustomer(req.CampaignID, coupon, couponCode, param.AccountId, 1)
			var useErr string
			for _, value := range use.Coupons {
				useErr = value.CouponID
			}

			if err2 != nil || useErr == "" {

				fmt.Println("[VoucherAgServices]-[CouponVoucherCustomer]")
				fmt.Println(fmt.Sprintf("[VoucherAgServices]-[Error : %v]", err2))
				sugarLogger.Info("[VoucherAgServices]-[CouponVoucherCustomer]")
			}

		}

		param.TrxID = utils.GenTransactionId()

		text := param.TrxID + param.InstitutionID + constants.CodeReversal + "#" + "OP009 - Reversal point cause transaction " + param.NamaVoucher + " is failed"

		point := param.Point * req.Jumlah
		totalPoint := strconv.Itoa(point)

		schedulerData := dbmodels.TSchedulerRetry{
			Code:          constants.CodeScheduler,
			TransactionID: utils.Before(text, "#"),
			Count:         0,
			IsDone:        false,
			CreatedAT:     time.Now(),
		}

		reversal, errReversal := host.TransferPoint(param.AccountId, totalPoint, text)

		statusEarning := constants.Success
		msgEarning := constants.MsgSuccess

		if errReversal != nil || reversal.PointsTransferId == "" {

			statusEarning = constants.TimeOut
			statusEarning = constants.TimeOut

			for _, val1 := range reversal.Form.Children.Customer.Errors {
				if val1 != "" {
					msgEarning = val1
					statusEarning = constants.Failed
				}
			}

			for _, val2 := range reversal.Form.Children.Points.Errors {
				if val2 != "" {
					msgEarning = val2
					statusEarning = constants.Failed
				}
			}

			if reversal.Message != "" {
				msgEarning = reversal.Message
				statusEarning = constants.Failed
			}

			if reversal.Error.Message != "" {
				msgEarning = reversal.Error.Message
				statusEarning = constants.Failed
			}

			if statusEarning == constants.TimeOut {
				errSaveScheduler := db.DbCon.Create(&schedulerData).Error
				if errSaveScheduler != nil {

					fmt.Println("===== Gagal SaveScheduler ke DB =====")
					fmt.Println(fmt.Sprintf("Error : %v", errSaveScheduler))
					fmt.Println(fmt.Sprintf("===== Phone : %v || RRN : %v =====", param.AccountNumber, param.RRN))

				}

			}

		}

		expired := ExpiredPointService()

		saveReversal := dbmodels.TEarning{
			ID:               utils.GenerateTokenUUID(),
			PartnerId:        param.InstitutionID,
			TransactionId:    param.TrxID,
			AccountNumber:    param.AccountNumber,
			Point:            int64(point),
			Status:           statusEarning,
			StatusMessage:    msgEarning,
			PointsTransferId: reversal.PointsTransferId,
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
			ReferenceId:    param.RRN,
			TransactionId:  param.Reffnum,
			Data: models.DataValue{
				RewardValue: "point",
				Value:       totalPoint,
			},
		}

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

			go SaveTransactionVoucherAg(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "01")

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

	// Save To Scheduler
	text := param.CumReffnum + param.InstitutionID + constants.CodeReversal + "#" + "OP009 - Reversal point cause transaction " + param.NamaVoucher + " is failed"
	schedulerData := dbmodels.TSchedulerRetry{
		Code:          constants.CodeSchedulerVoucherAG,
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

	// Check Order Status
	statusOrder, errStatus := vg.CheckStatusOrder(vgmodels.RequestCheckOrderStatus{
		OrderID:       param.CumReffnum,
		RecordPerPage: fmt.Sprintf("%d", req.Jumlah),
		CurrentPage:   "1",
	}, head)
	if errStatus != nil {

		// Handle Error Here
		fmt.Println("Internal Server Error : ", errorder)
		fmt.Println("ResponseCode : ", order.ResponseCode)
		fmt.Println("[VoucherAgServices]-[FailedCheckOrderStatus]")
		fmt.Println(fmt.Sprintf("[Response %v]", order.ResponseCode))

		sugarLogger.Info("[VoucherAgServices]-[FailedCheckOrderStatus]")
		sugarLogger.Info("[Failed Check Order Status]-[Gagal Order Voucher]")

	}

	param.RRN = statusOrder.Data.TransactionID

	for i := req.Jumlah; i > 0; i-- {

		fmt.Println(fmt.Sprintf("[Line Save DB : %v]", i))

		t := i - 1
		coupon := redeem.Coupons[t].Id
		code := redeem.Coupons[t].Code
		param.CouponID = coupon

		couponID := statusOrder.Data.Vouchers[t].VoucherID
		couponCode := statusOrder.Data.Vouchers[t].VoucherCode
		expDate := statusOrder.Data.Vouchers[t].ExpiredDate
		voucherLink := statusOrder.Data.Vouchers[t].Link

		a := []rune(coupon)
		key32 := string(a[0:32])
		key := []byte(key32)
		chiperText := []byte(couponCode)
		plaintText, err := utils.EncryptAES(chiperText, key)
		if err != nil {
			res = utils.GetMessageFailedErrorNew(res, constants.RC_FAILED_DECRYPT_VOUCHER, constants.RD_FAILED_DECRYPT_VOUCHER)
			return res
		}

		brand, err := db.GetBrandCode(param.ProductCode)
		if err != nil {
			fmt.Println("err")
		}
		// Use Voucher ID as a Transaction ID
		param.TrxID = couponID
		param.ExpDate = expDate
		param.CouponCode = fmt.Sprintf("%s", plaintText)
		param.VoucherLink = voucherLink

		if brand != nil && brand.Code != "" {
			param.ProductCode = brand.Code
		}

		id := utils.GenerateTokenUUID()
		go SaveDBVoucherAg(id, param.InstitutionID, coupon, code, param.AccountNumber, param.AccountId, req.CampaignID)
		go SaveTransactionVoucherAg(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "00")

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

func (t VoucherAgServices) HandleCallback(req models.CallbackRequestVoucherAg) models.Response {
	var res models.Response

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[HandleCallbackVoucherAg-Services]",
		zap.String("Institution : ", req.InstitutionID),
		zap.String("TransactionID : ", req.TransactionID),
		zap.String("OrderID : ", req.Data.OrderID),
		zap.String("VoucherID: ", req.Data.VoucherID),
		zap.String("VoucherName: ", req.Data.VoucherName),
		zap.String("VoucherCode: ", req.Data.VoucherCode),
		zap.String("Status: ", req.Data.Status),
		zap.Bool("IsRedeemed: ", req.Data.IsRedeemed),
		zap.String("RedeemedDate: ", req.Data.RedeemedDate),
		zap.String("UsedDate: ", req.Data.UsedDate),
	)

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[RedeemVoucher]")
	defer span.Finish()

	// Get TSpending
	tspending, err := db.GetVoucherAgSpending(req.Data.VoucherID, req.TransactionID)
	if err != nil {

		fmt.Println("[HandleCallbackVoucherAg]")
		fmt.Println("[FailedGetTSpending]: ", err)
		sugarLogger.Info("[VoucherAgServices]-[HandleCallback]-[CouponVoucherCustomer]")

		res = utils.GetMessageResponse(res, 422, false, err)

		return res
	}

	cekVoucher, errVoucher := opl.VoucherDetail(tspending.CampaignId)
	if errVoucher != nil || cekVoucher.CampaignID == "" {
		sugarLogger.Info("[HandleCallback]-[VoucherDetail]")
		sugarLogger.Info(fmt.Sprintf("Error : ", errVoucher))

		logs.Info("[HandleCallback]-[VoucherDetail]")
		logs.Info(fmt.Sprintf("Error : ", errVoucher))
	}

	couponCode := cekVoucher.Coupons[0]

	// Use Voucher
	use, err2 := opl.CouponVoucherCustomer(tspending.CampaignId, tspending.CouponId, couponCode, tspending.AccountId, 1)
	var useErr string
	for _, value := range use.Coupons {
		useErr = value.CouponID
	}

	if err2 != nil || useErr == "" {

		fmt.Println("[VoucherAgServices]-[HandleCallback]-[CouponVoucherCustomer]")
		fmt.Println(fmt.Sprintf("[VoucherAgServices]-[HandleCallback]-[Error : %v]", err2))
		sugarLogger.Info("[VoucherAgServices]-[HandleCallback]-[CouponVoucherCustomer]")

	}

	// Update TSpending
	_, err3 := db.UpdateVoucherAg(req.Data.RedeemedDate, req.Data.UsedDate, tspending.ID)
	if err3 != nil {

		fmt.Println("[VoucherAgServices]-[HandleCallback]-[FailedUpdate]")
		fmt.Println(fmt.Sprintf("[VoucherAgServices]-[HandleCallback]-[Error : %v]", err3))
		sugarLogger.Info("[VoucherAgServices]-[HandleCallback]-[FailedUpdateTSpending]")

	}

	go db.UpdateTSchedulerVoucherAG(req.Data.OrderID)

	res = models.Response{
		Meta: utils.ResponseMetaOK(),
	}

	return res
}

func SaveDBVoucherAg(id, institution, coupon, vouchercode, phone, custIdOPL, campaignID string) {

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
	}
}

func SaveTransactionVoucherAg(param models.Params, res interface{}, reqdata interface{}, reqOP interface{}, transType, status string) {

	fmt.Println(fmt.Sprintf("[Start-SaveDB]-[VoucherAggregator]-[%v]", transType))

	var saveStatus string
	switch status {
	case "00":
		saveStatus = constants.Success
	case "09":
		saveStatus = constants.Pending
	case "01":
		saveStatus = constants.Failed
	}

	reqUV, _ := json.Marshal(&reqdata)   // Req UV
	responseUV, _ := json.Marshal(&res)  // Response UV
	reqdataOP, _ := json.Marshal(&reqOP) // Req Service

	expDate := ""
	if param.ExpDate != "" {
		layout := "2006-01-02 15:04:05"
		parse, _ := time.Parse(layout, param.ExpDate)

		expDate = jodaTime.Format("YYYY-MM-dd", parse)
	}

	save := dbmodels.TSpending{
		ID:              utils.GenerateTokenUUID(),
		AccountNumber:   param.AccountNumber,
		RRN:             param.RRN,
		Voucher:         param.NamaVoucher,
		MerchantID:      param.MerchantID,
		TransactionId:   param.TrxID,
		ProductCode:     param.ProductCode,
		Amount:          int64(param.Amount),
		TransType:       transType,
		IsUsed:          false,
		ProductType:     param.ProductType,
		Status:          saveStatus,
		Institution:     param.InstitutionID,
		CummulativeRef:  param.CumReffnum,
		DateTime:        utils.GetTimeFormatYYMMDDHHMMSS(),
		Point:           param.Point,
		ResponderRc:     param.DataSupplier.Rc,
		ResponderRd:     param.DataSupplier.Rd,
		RequestorData:   string(reqUV),
		ResponderData:   string(responseUV),
		RequestorOPData: string(reqdataOP),
		SupplierID:      param.SupplierID,
		CouponId:        param.CouponID,
		CampaignId:      param.CampaignID,
		AccountId:       param.AccountId,
		VoucherCode:     param.CouponCode,
		VoucherLink:     param.VoucherLink,
		ExpDate:         expDate,
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

func DataParameterOrderVoucherAg() models.ParamUV {

	res := models.ParamUV{}

	nama := ""
	email := "VOUCHER_AG_EMAIL_ORDER"
	phone := "VOUCHER_AG_PHONE_ORDER"
	expired := "VOUCHER_AG_EXPIRED_VOUCHER"

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
