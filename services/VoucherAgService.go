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
	callbackOttoPointPurchase = ODU.GetEnv("OTTOPOINT_PURCHASE_CALLBACK_VOUCHERAG", "http://34.101.119.111:8006/transaction/v2/callback/redeem/voucherag")
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
	reqOrder := vgmodels.RequestOrderVoucherAg{
		ProductCode:    param.ProductCode,
		Qty:            req.Jumlah,
		OrderID:        param.CumReffnum,
		CustomerName:   nama,
		CustomerEmail:  dataOrder.Email,
		CustomerPhone:  dataOrder.Phone,
		DeliveryMethod: 1,
		RedeemCallback: callbackOttoPointPurchase,
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

	}

	// Handle Stock Not Available
	if order.ResponseCode == "04" {

		// Start Reversal Here

		for i := req.Jumlah; i > 0; i-- {

			fmt.Println(fmt.Sprintf("[Line Save DB : %v]", i))

			// Generate TransactionID
			param.TrxID = utils.GenTransactionId()

			go SaveTransactionVoucherAg(param, order, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "01")

		}

		res = models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: vgmodels.ResponseVoucherAg{
				Code:    "176",
				Msg:     "Voucher tidak tersedia",
				Success: 0,
				Failed:  req.Jumlah,
				Pending: 0,
			},
		}

		return res

	}

	// Handle General Error
	if order.ResponseCode != "00" {

		// Reversal Start Here

		for i := req.Jumlah; i > 0; i-- {

			fmt.Println(fmt.Sprintf("[Line Save DB : %v]", i))

			// Generate TransactionID
			param.TrxID = utils.GenTransactionId()

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
		fmt.Println("[Failed Order Voucher]-[Gagal Order Voucher]")

		// sugarLogger.Info("Internal Server Error : ", errOrder)
		sugarLogger.Info("[VoucherAggregator]-[OrderVoucher]")
		sugarLogger.Info("[Failed Order Voucher]-[Gagal Order Voucher]")

		for i := req.Jumlah; i > 0; i-- {

			// TrxId
			param.TrxID = utils.GenTransactionId()

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

	// Check Order Status
	statusOrder, errStatus := vg.CheckStatusOrder(vgmodels.RequestCheckOrderStatus{
		OrderID:       param.CumReffnum,
		RecordPerPage: fmt.Sprintf("%d", req.Jumlah),
		CurrentPage:   "1",
	}, head)
	if errStatus != nil {

		// Handle Error Here

	}

	param.RRN = statusOrder.Data.TransactionID

	fmt.Println("[OrderStatus] : ", statusOrder)
	// Handle General Error
	if order.ResponseCode != "00" {

		// Reversal Start Here

		for i := req.Jumlah; i > 0; i-- {

			fmt.Println(fmt.Sprintf("[Line Save DB : %v]", i))

			// Generate TransactionID
			param.TrxID = utils.GenTransactionId()

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
		fmt.Println("[Failed Order Voucher]-[Gagal Order Voucher]")

		// sugarLogger.Info("Internal Server Error : ", errOrder)
		sugarLogger.Info("[VoucherAggregator]-[OrderVoucher]")
		sugarLogger.Info("[Failed Order Voucher]-[Gagal Order Voucher]")

		for i := req.Jumlah; i > 0; i-- {

			// TrxId
			param.TrxID = utils.GenTransactionId()

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

	for i := req.Jumlah; i > 0; i-- {

		fmt.Println(fmt.Sprintf("[Line Save DB : %v]", i))

		// Generate TransactionID
		param.TrxID = utils.GenTransactionId()

		t := i - 1
		coupon := statusOrder.Data.Vouchers[t].VoucherID
		code := statusOrder.Data.Vouchers[t].VoucherCode
		expDate := statusOrder.Data.Vouchers[t].ExpiredDate
		voucherLink := statusOrder.Data.Vouchers[t].Link

		param.ExpDate = expDate
		param.CouponID = coupon
		param.VoucherLink = voucherLink

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
		// return err
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
		CustID:          param.AccountId,
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
