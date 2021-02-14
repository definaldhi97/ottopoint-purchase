package redeemtion

import (
	"encoding/json"
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	kafka "ottopoint-purchase/hosts/publisher/host"
	signature "ottopoint-purchase/hosts/signature/host"
	vgmodels "ottopoint-purchase/hosts/voucher_aggregator/models"
	"ottopoint-purchase/models"
	"ottopoint-purchase/models/dbmodels"
	"ottopoint-purchase/services/v2/Trx"
	"ottopoint-purchase/utils"
	"reflect"
	"strconv"
	"time"

	vg "ottopoint-purchase/hosts/voucher_aggregator/host"

	"github.com/astaxie/beego/logs"
	"github.com/sirupsen/logrus"
	"github.com/vjeantet/jodaTime"
	ODU "ottodigital.id/library/utils"
)

var (
	HostPurcahse              = ODU.GetEnv("OTTOPOINT_PURCHASE_HOST_CALLBACK_VOUCHERAG", "http://13.228.25.85:8006")
	CallbackOttoPointPurchase = ODU.GetEnv("OTTOPOINT_PURCHASE_CALLBACK_VOUCHERAG", "/transaction/v2/redeem/voucherag")
)

// func (t V2_VoucherAgServices) VoucherAg(req models.VoucherComultaiveReq, param models.Params, head models.RequestHeader) models.Response {
func RedeemtionAggServices(req models.VoucherComultaiveReq, param models.Params, head models.RequestHeader) models.Response {
	// fmt.Println("[ >>>>>>>>>>>>>>>>>> V2 Migrate Voucher Agregator Service <<<<<<<<<<<<<<<< ]")

	nameservice := "[PackageRedeemtion]-[RedeemtionAggServices]"
	logReq := fmt.Sprintf("[AccountNumber : %v, RewardID : %v]", param.AccountNumber, param.RewardID)

	logrus.Info(nameservice)

	var res models.Response

	dataOrder := DataParameterOrderVoucherAg()

	timeExp, _ := strconv.Atoi(dataOrder.Expired)

	param.CumReffnum = utils.GenTransactionId()

	param.Amount = int64(param.Point)

	// spending point and spending usage_limit voucher
	textCommentSpending := param.CumReffnum + "#" + param.NamaVoucher
	param.Comment = textCommentSpending

	RedeemVouchAG, errRedeemVouchAG := Trx.V2_Redeem_PointandVoucher(req.Jumlah, param)

	param.PointTransferID = RedeemVouchAG.PointTransferID

	if RedeemVouchAG.Rd == "Invalid JWT Token" {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[V2_Redeem_PointandVoucher]-[Response : %v]", RedeemVouchAG))
		logrus.Println("[ ResponseCode ] : ", RedeemVouchAG.Rc)
		logrus.Println("[ ResponseDesc ] : ", RedeemVouchAG.Rd)
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

	if RedeemVouchAG.Rd == "not enough points" {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[V2_Redeem_PointandVoucher]-[Response : %v]", RedeemVouchAG))
		logrus.Println("[ ResponseCode ] : ", RedeemVouchAG.Rc)
		logrus.Println("[ ResponseDesc ] : ", RedeemVouchAG.Rd)
		logrus.Println(logReq)

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

	if RedeemVouchAG.Rc == "208" {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[V2_Redeem_PointandVoucher]-[Response : %v]", RedeemVouchAG))
		logrus.Println("[ ResponseCode ] : ", RedeemVouchAG.Rc)
		logrus.Println("[ ResponseDesc ] : ", RedeemVouchAG.Rd)
		logrus.Println(logReq)

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

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[V2_Redeem_PointandVoucher]-[Response : %v]", RedeemVouchAG))
		logrus.Println("[ ResponseCode ] : ", RedeemVouchAG.Rc)
		logrus.Println("[ ResponseDesc ] : ", RedeemVouchAG.Rd)
		logrus.Println(logReq)

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

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[V2_Redeem_PointandVoucher]-[Error : %v]", errRedeemVouchAG))
		logrus.Println("[ ResponseCode ] : ", RedeemVouchAG.Rc)
		logrus.Println("[ ResponseDesc ] : ", RedeemVouchAG.Rd)
		logrus.Println(logReq)

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
		ProductCode:    param.ProductCodeInternal,
		Qty:            req.Jumlah,
		OrderID:        param.CumReffnum,
		CustomerName:   nama,
		CustomerEmail:  dataOrder.Email,
		CustomerPhone:  dataOrder.Phone,
		DeliveryMethod: 1,
		RedeemCallback: HostPurcahse + CallbackOttoPointPurchase,
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

	order, errorder := vg.OrderVoucher(reqOrder, head)

	param.DataSupplier.Rd = order.ResponseDesc
	param.DataSupplier.Rc = order.ResponseCode

	if errorder != nil || order.ResponseCode == "" {

		// Reversal Start Here
		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[OrderVoucher]-[Error : %v]", errorder))
		logrus.Println(logReq)

		for i := req.Jumlah; i > 0; i-- {
			// TrxId
			param.TrxID = utils.GenTransactionId()

			fmt.Println(fmt.Sprintf("[Line Save DB : %v]", i))

			t := i - 1
			coupon := RedeemVouchAG.CouponseVouch[t].CouponsID
			param.CouponID = coupon
			go SaveTSchedulerRetry(param.TrxID, constants.CodeSchedulerVoucherAG)
			go SaveTransactionVoucherAgMigrate(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "09", timeExp)
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

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[OrderVoucher]-[Response : %v]", order))
		logrus.Println("[ ResponseCode ] : ", order.ResponseCode)
		logrus.Println("[ ResponseDesc ] : ", order.ResponseDesc)
		logrus.Println(logReq)

		totalPoint := param.Point * req.Jumlah
		param.TrxID = utils.GenTransactionId()
		resultReversal := Trx.V2_Adding_PointVoucher(param, totalPoint, req.Jumlah)
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

			go SaveTransactionVoucherAgMigrate(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "01", timeExp)
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

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[OrderVoucher]-[Response : %v]", order))
		logrus.Println("[ ResponseCode ] : ", order.ResponseCode)
		logrus.Println("[ ResponseDesc ] : ", order.ResponseDesc)
		logrus.Println(logReq)

		for i := req.Jumlah; i > 0; i-- {
			fmt.Println(fmt.Sprintf("[Line : %v]", i))

			// TrxId
			param.TrxID = utils.GenTransactionId()

			fmt.Println(fmt.Sprintf("[Line Save DB : %v]", i))

			t := i - 1
			coupon := RedeemVouchAG.CouponseVouch[t].CouponsID
			param.CouponID = coupon
			go SaveTSchedulerRetry(param.TrxID, constants.CodeSchedulerVoucherAG)
			go SaveTransactionVoucherAgMigrate(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "09", timeExp)

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
		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[OrderVoucher]-[Response : %v]", order))
		logrus.Println("[ ResponseCode ] : ", order.ResponseCode)
		logrus.Println("[ ResponseDesc ] : ", order.ResponseDesc)
		logrus.Println(logReq)

		totalPoint := param.Point * req.Jumlah
		param.TrxID = utils.GenTransactionId()
		resultReversal := Trx.V2_Adding_PointVoucher(param, totalPoint, req.Jumlah)
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

		_, errKafka := kafka.SendPublishKafka(kafkaReq)
		if errKafka != nil {

			logrus.Error(nameservice)
			logrus.Error(fmt.Sprintf("[SendPublishKafka]-[Error : %v]", errKafka))
			logrus.Println(logReq)
		}

		// fmt.Println("Response Publisher : ", kafkaRes)

		for i := req.Jumlah; i > 0; i-- {

			fmt.Println(fmt.Sprintf("[Line Save DB : %v]", i))

			// TrxID
			param.TrxID = utils.GenTransactionId()

			t := i - 1
			coupon := RedeemVouchAG.CouponseVouch[t].CouponsID
			param.CouponID = coupon

			go SaveTransactionVoucherAgMigrate(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "01", timeExp)

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

	// Check Order Status
	statusOrder, errStatus := vg.CheckStatusOrder(vgmodels.RequestCheckOrderStatus{
		OrderID:       param.CumReffnum,
		RecordPerPage: fmt.Sprintf("%d", req.Jumlah),
		CurrentPage:   "1",
	}, head)

	if errStatus != nil {

		// Handle Error Here
		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[CheckStatusOrder]-[Error : %v]", errStatus))
		logrus.Println(logReq)

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

		// a := []rune(param.CouponID)
		// key32 := string(a[0:32])
		// key := []byte(key32)
		// chiperText := []byte(voucherCode)
		// plaintText, err := utils.EncryptAES(chiperText, key)
		// if err != nil {
		// 	res = utils.GetMessageFailedErrorNew(res, constants.RC_FAILED_DECRYPT_VOUCHER, constants.RD_FAILED_DECRYPT_VOUCHER)
		// 	return res
		// }

		// Use Voucher ID as a Transaction ID
		param.TrxID = utils.GenTransactionId()
		param.ExpDate = expDate
		// param.CouponCode = fmt.Sprintf("%s", plaintText)
		param.CouponCode = voucherCode
		param.VoucherLink = voucherLink

		id := utils.GenerateTokenUUID()
		go SaveDBVoucherAgMigrate(id, param.InstitutionID, param.CouponID, voucherCode, param.AccountNumber, param.AccountId, req.CampaignID)
		go SaveTransactionVoucherAgMigrate(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "00", timeExp)

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

func DataParameterOrderVoucherAg() models.ParamUV {

	res := models.ParamUV{}

	datanama, errnama := db.ParamData(constants.CODE_CONFIG_AGG_GROUP, constants.CODE_CONFIG_AGG_NAME)
	if errnama != nil {
		fmt.Println("[Error get data from Db m_paramaters]")
		fmt.Println("Error : ", errnama)
		fmt.Println("Code :", constants.CODE_CONFIG_AGG_NAME)
	}

	dataemail, erremail := db.ParamData(constants.CODE_CONFIG_AGG_GROUP, constants.CODE_CONFIG_AGG_EMAIL)
	if erremail != nil {
		fmt.Println("[Error get data from Db m_paramaters]")
		fmt.Println("Error : ", erremail)
		fmt.Println("Code :", constants.CODE_CONFIG_AGG_EMAIL)
	}

	dataphone, errphone := db.ParamData(constants.CODE_CONFIG_AGG_GROUP, constants.CODE_CONFIG_AGG_PHONE)
	if errphone != nil {
		fmt.Println("[Error get data from Db m_paramaters]")
		fmt.Println("Error : ", errphone)
		fmt.Println("Code :", constants.CODE_CONFIG_AGG_PHONE)
	}

	dataexpired, errexpired := db.ParamData(constants.CODE_CONFIG_AGG_GROUP, constants.CODE_CONFIG_AGG_EXPD)
	if errexpired != nil {
		fmt.Println("[Error get data from Db m_paramaters]")
		fmt.Println("Error : ", errexpired)
		fmt.Println("Code :", constants.CODE_CONFIG_AGG_EXPD)
	}

	res = models.ParamUV{
		Nama:    datanama.Value,
		Email:   dataemail.Value,
		Phone:   dataphone.Value,
		Expired: dataexpired.Value,
	}

	return res

}

func SaveTransactionVoucherAgMigrate(param models.Params, res interface{}, reqdata interface{}, reqOP interface{}, transType, status string, timeExpVouc int) {

	fmt.Println(fmt.Sprintf("[Start-SaveDB]-[VoucherAggregator]-[%v]", transType))

	var ExpireDate time.Time
	var redeemDate time.Time

	var saveStatus string
	isUsed := false
	switch status {
	case "00":
		saveStatus = constants.Success
	case "09":
		saveStatus = constants.Pending
	case "01":
		saveStatus = constants.Failed
		isUsed = true
	}

	reqUV, _ := json.Marshal(&reqdata)   // Req UV
	responseUV, _ := json.Marshal(&res)  // Response UV
	reqdataOP, _ := json.Marshal(&reqOP) // Req Service

	// expDate := ""
	// if param.ExpDate != "" {
	// 	layout := "2006-01-02 15:04:05"
	// 	parse, _ := time.Parse(layout, param.ExpDate)

	// 	expDate = jodaTime.Format("YYYY-MM-dd", parse)
	// }

	if transType == constants.CODE_TRANSTYPE_REDEMPTION {
		ExpireDate = utils.ExpireDateVoucherAGt(timeExpVouc)
		redeemDate = time.Now()
	}

	save := dbmodels.TSpending{
		ID:              utils.GenerateTokenUUID(),
		AccountNumber:   param.AccountNumber,
		RRN:             param.RRN,
		Voucher:         param.NamaVoucher,
		MerchantID:      param.MerchantID,
		TransactionId:   param.TrxID,
		ProductCode:     param.ProductCodeInternal,
		Amount:          int64(param.Amount),
		TransType:       transType,
		IsUsed:          isUsed,
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
		ExpDate:         utils.DefaultNulTime(ExpireDate),
		RedeemAt:        utils.DefaultNulTime(redeemDate),

		Comment:           param.Comment,
		MRewardID:         param.RewardID,
		ProductCategoryID: param.CategoryID,
		MProductID:        param.ProductID,
		PointsTransferID:  param.PointTransferID,
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

func SaveDBVoucherAgMigrate(id, institution, coupon, vouchercode, phone, custIdOPL, campaignID string) {

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
