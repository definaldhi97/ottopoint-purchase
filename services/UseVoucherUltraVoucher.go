package services

import (
	"fmt"
	"ottopoint-purchase/db"
	"ottopoint-purchase/hosts/opl/host"
	opl "ottopoint-purchase/hosts/opl/host"
	ottomart "ottopoint-purchase/hosts/ottomart/host"
	ottomartmodels "ottopoint-purchase/hosts/ottomart/models"
	uv "ottopoint-purchase/hosts/ultra_voucher/host"
	uvmodels "ottopoint-purchase/hosts/ultra_voucher/models"
	"ottopoint-purchase/models"
	"ottopoint-purchase/models/dbmodels"
	"ottopoint-purchase/utils"
	"strconv"
	"time"

	"github.com/opentracing/opentracing-go"
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

	param.Reffnum = utils.GenTransactionId()

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
		InstitutionRefno:  param.Reffnum,
		ExpireDateVoucher: expired,
		ReceiverName:      nama,
		ReceiverEmail:     dataorder.Email,
		ReceiverPhone:     dataorder.Phone,
	}

	// order to u
	fmt.Println("[OrderVoucher]")
	order, errOrder := uv.OrderVoucher(reqOrder, param.InstitutionID)

	param.DataSupplier.Rd = order.ResponseDesc
	param.DataSupplier.Rc = order.ResponseCode

	if errOrder != nil || order.ResponseCode == "" {

		fmt.Println("Error : ", errOrder)
		fmt.Println("Response OrderVoucher : ", order)
		fmt.Println("[UltraVoucherServices]-[OrderVoucher]")
		fmt.Println("[Failed Order Voucher]-[Gagal Order Voucher]")

		// sugarLogger.Info("Internal Server Error : ", errOrder)
		sugarLogger.Info("[UltraVoucherServices]-[OrderVoucher]")
		sugarLogger.Info("[Failed Order Voucher]-[Gagal Order Voucher]")

		fmt.Println("[CheckStatusOrder]")
		time.Sleep(5 * time.Second)

		reqCheckStatus := map[string]interface{}{
			"InstitutionRefno": reqOrder.InstitutionRefno,
			"InstitutionID":    param.InstitutionID,
		}

		checkOrder, errCheck := chekStatus(reqOrder.InstitutionRefno, param.InstitutionID)

		param.DataSupplier.Rd = checkOrder.ResponseDesc
		param.DataSupplier.Rc = checkOrder.ResponseCode
		// checkOrder, errCheck := uv.CheckStatusOrder(reqOrder.InstitutionRefno, param.InstitutionID)

		fmt.Println("Response Check Status : ", checkOrder)

		if checkOrder.ResponseCode == "" || checkOrder.ResponseCode == "18" {
			fmt.Println("Error : ", errCheck)
			fmt.Println("Response CheckStatusOrder : ", checkOrder)
			fmt.Println("[UltraVoucherServices]-[CheckStatusOrder]")
			fmt.Println("[TimeOunt or Pending CheckStatusOrder]-[Gagal CheckStatusOrder Voucher]")

			// sugarLogger.Info("Internal Server Error : ", errOrder)
			sugarLogger.Info("[UltraVoucherServices]-[CheckStatusOrder]")
			sugarLogger.Info("[TimeOunt or Pending CheckStatusOrder]-[Gagal CheckStatusOrder]")

			for i := req.Jumlah; i > 0; i-- {

				fmt.Println(fmt.Sprintf("[Line Save DB : %v]", i))

				t := i - 1
				coupon := redeem.Coupons[t].Id
				param.CouponID = coupon

				go SaveTransactionUV(param, checkOrder, reqCheckStatus, req, "Reedemtion", "09")
			}

			// res = utils.GetMessageResponse(res, 500, false, errors.New("Transaksi Anda sedang dalam proses. Silahkan hubungi tim kami untuk informasi selengkapnya."))
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

		if checkOrder.ResponseCode != "00" {
			fmt.Println("Error : ", errCheck)
			fmt.Println("Response CheckStatusOrder : ", checkOrder)
			fmt.Println("[UltraVoucherServices]-[CheckStatusOrder]")
			fmt.Println("[Failed CheckStatusOrder]-[Gagal CheckStatusOrder Voucher]")

			// sugarLogger.Info("Internal Server Error : ", errOrder)
			sugarLogger.Info("[UltraVoucherServices]-[CheckStatusOrder]")
			sugarLogger.Info("[Failed CheckStatusOrder]-[Gagal CheckStatusOrder]")

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

			fmt.Println("========== Reversal  ==========")
			Text := "OP009 - " + "Reversal point cause transaction " + param.NamaVoucher + " is failed"
			point := param.Point * req.Jumlah
			totalPoint := strconv.Itoa(point)
			_, errReversal := host.TransferPoint(param.AccountId, totalPoint, Text)
			if errReversal != nil {
				fmt.Println("Internal Server Error : ", errReversal)
				fmt.Println("[UltraVoucherServices]-[TransferPoint]")
				fmt.Println("[Failed Order Voucher]-[Gagal Reversal Point]")
			}

			fmt.Println("========== Send Notif ==========")
			notifReq := ottomartmodels.NotifRequest{
				AccountNumber:    param.AccountNumber,
				Title:            "Reversal Point",
				Message:          fmt.Sprintf("Point anda berhasil di reversal sebesar %v", int64(point)),
				NotificationType: 3,
			}

			// send notif & inbox
			dataNotif, errNotif := ottomart.NotifAndInbox(notifReq)
			if errNotif != nil {
				fmt.Println("Error to send Notif & Inbox")
			}

			if dataNotif.RC != "00" {
				fmt.Println(fmt.Sprintf("[Response : %v\n]", dataNotif))
				fmt.Println("Gagal Send Notif & Inbox")
				fmt.Println("Error : ", errNotif)
			}

			for i := req.Jumlah; i > 0; i-- {

				fmt.Println(fmt.Sprintf("[Line Save DB : %v]", i))

				t := i - 1
				coupon := redeem.Coupons[t].Id
				param.CouponID = coupon

				go SaveTransactionUV(param, checkOrder, reqCheckStatus, req, "Reedemtion", "01")
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

		fmt.Println("Response UV checkOrder : ", checkOrder)
		for i := req.Jumlah; i > 0; i-- {

			fmt.Println(fmt.Sprintf("[Line Save DB : %v]", i))

			t := i - 1
			coupon := redeem.Coupons[t].Id
			code := checkOrder.Data.VouchersCode[t].Code

			param.CouponID = coupon

			id := utils.GenerateTokenUUID()
			go SaveDB(id, param.InstitutionID, coupon, code, param.AccountNumber, param.AccountId, req.CampaignID)
			go SaveTransactionUV(param, checkOrder, reqCheckStatus, req, "Reedemtion", "00")
		}

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

		Text := "OP009 - " + "Reversal point cause transaction " + param.NamaVoucher + " is failed"
		point := param.Point * req.Jumlah
		totalPoint := strconv.Itoa(point)
		_, errReversal := host.TransferPoint(param.AccountId, totalPoint, Text)
		if errReversal != nil {
			fmt.Println("Internal Server Error : ", errReversal)
			fmt.Println("[UltraVoucherServices]-[TransferPoint]")
			fmt.Println("[Failed Order Voucher]-[Gagal Reversal Point]")
		}

		fmt.Println("========== Send Notif ==========")
		notifReq := ottomartmodels.NotifRequest{
			AccountNumber:    param.AccountNumber,
			Title:            "Reversal Point",
			Message:          fmt.Sprintf("Point anda berhasil di reversal sebesar %v", int64(point)),
			NotificationType: 3,
		}

		// send notif & inbox
		dataNotif, errNotif := ottomart.NotifAndInbox(notifReq)
		if errNotif != nil {
			fmt.Println("Error to send Notif & Inbox")
		}

		if dataNotif.RC != "00" {
			fmt.Println(fmt.Sprintf("[Response : %v\n]", dataNotif))
			fmt.Println("Gagal Send Notif & Inbox")
			fmt.Println("Error : ", errNotif)
		}

		for i := req.Jumlah; i > 0; i-- {

			fmt.Println(fmt.Sprintf("[Line Save DB : %v]", i))

			t := i - 1
			coupon := redeem.Coupons[t].Id
			param.CouponID = coupon

			go SaveTransactionUV(param, order, reqOrder, req, "Reedemtion", "01")
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

		Text := "OP009 - " + "Reversal point cause transaction " + param.NamaVoucher + " is failed"
		point := param.Point * req.Jumlah
		totalPoint := strconv.Itoa(point)
		reversal, errReversal := host.TransferPoint(param.AccountId, totalPoint, Text)
		if errReversal != nil || reversal.PointsTransferId == "" {
			fmt.Println("Internal Server Error : ", errReversal)
			fmt.Println("[UltraVoucherServices]-[TransferPoint]")
			fmt.Println("[Failed Order Voucher]-[Gagal Reversal Point]")
		}

		fmt.Println("========== Send Notif ==========")
		notifReq := ottomartmodels.NotifRequest{
			AccountNumber:    param.AccountNumber,
			Title:            "Reversal Point",
			Message:          fmt.Sprintf("Point anda berhasil di reversal sebesar %v", int64(point)),
			NotificationType: 3,
		}

		// send notif & inbox
		dataNotif, errNotif := ottomart.NotifAndInbox(notifReq)
		if errNotif != nil {
			fmt.Println("Error to send Notif & Inbox")
		}

		if dataNotif.RC != "00" {
			fmt.Println(fmt.Sprintf("[Response : %v\n]", dataNotif))
			fmt.Println("Gagal Send Notif & Inbox")
			fmt.Println("Error : ", errNotif)
		}

		for i := req.Jumlah; i > 0; i-- {

			fmt.Println(fmt.Sprintf("[Line Save DB : %v]", i))

			t := i - 1
			coupon := redeem.Coupons[t].Id
			param.CouponID = coupon

			go SaveTransactionUV(param, order, reqOrder, req, "Reedemtion", "01")
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

		fmt.Println(">>> Order Success <<<")

		fmt.Println(fmt.Sprintf("[Line Save DB : %v]", i))

		t := i - 1
		coupon := redeem.Coupons[t].Id
		code := order.Data.VouchersCode[t].Code

		param.CouponID = coupon

		id := utils.GenerateTokenUUID()
		go SaveDB(id, param.InstitutionID, coupon, code, param.AccountNumber, param.AccountId, req.CampaignID)

		go SaveTransactionUV(param, order, reqOrder, req, "Reedemtion", "00")
	}

	fmt.Println("Response UV : ", order)
	// res = models.Response{
	// 	Meta: utils.ResponseMetaOK(),
	// 	Data: models.UltraVoucherResp{
	// 		Success: req.Jumlah,
	// 		Failed:  0,
	// 		Total:   req.Jumlah,
	// 		Voucher: param.NamaVoucher,
	// 	},
	// }
	// go SaveTransactionUV(param, order, reqOrder, req, "Reedemtion", "00", order.ResponseCode)

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
