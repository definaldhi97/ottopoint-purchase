package voucher

import (
	"fmt"
	"ottopoint-purchase/db"
	redismodels "ottopoint-purchase/hosts/redis_token/models"
	"ottopoint-purchase/models"
	"ottopoint-purchase/models/dbmodels"
	ottoagmodels "ottopoint-purchase/models/ottoag"
	"ottopoint-purchase/services/ottoag"
	biller "ottopoint-purchase/services/ottoag"
	"ottopoint-purchase/utils"

	"github.com/astaxie/beego/logs"
)

// type RedeemPulsaServices struct {
// 	General models.GeneralModel
// }

func RedeemPulsa(req models.UseRedeemRequest, dataToken redismodels.TokenResp, MemberID, namaVoucher, expDate, category string) models.UseRedeemResponse {
	res := models.UseRedeemResponse{}

	// sugarLogger := t.General.OttoZaplog
	// sugarLogger.Info("[UseVoucher-Services]",
	// 	zap.String("category", req.AccountNumber))

	// span, _ := opentracing.StartSpanFromContext(t.General.Context, "[RedeemVoucher]")
	// defer span.Finish()

	// Validasi Prefix
	dataPrefix, errPrefix := db.GetOperatorCodebyPrefix(req.CustID)
	if errPrefix != nil {

		logs.Info("[ErrorPrefix]-[services-RedeemPulsa]")
		logs.Info(fmt.Sprintf("dataPrefix = %v", dataPrefix))
		logs.Info(fmt.Sprintf("Prefix tidak ditemukan %v", errPrefix))

		res = models.UseRedeemResponse{
			Msg: "Prefix Failed",
		}
		return res
	}

	// prefix := utils.Operator(dataPrefix.OperatorCode)

	prefixErr := ValidatePrefix(dataPrefix.OperatorCode, req.CustID, req.ProductCode)
	if prefixErr == false {
		res.Msg = "Prefix Failed"
		return res
	}

	// types := utils.TypeProduct(req.ProductCode[0:4])

	// ===== Inquiry OttoAG =====

	inqBiller := ottoagmodels.BillerInquiryDataReq{
		ProductCode: req.ProductCode,
		MemberID:    MemberID,
		CustID:      req.CustID,
	}

	inqReq := ottoagmodels.OttoAGInquiryRequest{
		TypeTrans:     "0003",
		Datetime:      utils.GetTimeFormatYYMMDDHHMMSS(),
		IssuerID:      "OTTOPAY",
		AccountNumber: req.AccountNumber,
		Data:          inqBiller,
	}

	if !ottoag.ValidateDataInq(inqReq) {
		res = models.UseRedeemResponse{
			Rc:  "01",
			Msg: "Inquiry Failed",
		}

		return res
	}

	inqRespOttoag := ottoagmodels.OttoAGInquiryResponse{}
	inqRespOttoag = biller.InquiryBiller(inqReq.Data, req, dataToken, MemberID, namaVoucher, expDate)

	if inqRespOttoag.Rc != "00" {
		logs.Info("[Error Inq Failed]")
		res = models.UseRedeemResponse{
			Rc:  "01",
			Msg: "Inquiry Failed",
		}

		// save DB
		labelInq := dbmodels.TransaksiRedeem{
			AccountNumber: dataToken.Data,
			Voucher:       namaVoucher,
			CustID:        req.CustID,
			// MerchantID:    dataToken.Data.MerchantID,
			RRN:         inqRespOttoag.Rrn,
			ProductCode: req.ProductCode,
			Amount:      inqRespOttoag.Amount,
			TransType:   "Inquiry",
			Status:      "01 (Gagal)",
			ExpDate:     expDate,
			Institution: "Ottopay",
			ProductType: "Pulsa",
			DateTime:    utils.GetTimeFormatYYMMDDHHMMSS(),
		}
		err1 := db.DbCon.Create(&labelInq).Error
		if err1 != nil {
			logs.Info("Failed Save to database", err1)
			// return err1
		}

		return res
	}

	logs.Info("[SAVE-DB-Transaksi_Redeem]")
	labelInq1 := dbmodels.TransaksiRedeem{
		AccountNumber: dataToken.Data,
		Voucher:       namaVoucher,
		CustID:        req.CustID,
		// MerchantID:    dataToken.Data.MerchantID,
		RRN:         inqRespOttoag.Rrn,
		ProductCode: req.ProductCode,
		Amount:      inqRespOttoag.Amount,
		TransType:   "Inquiry",
		ProductType: "Pulsa",
		Status:      "00 (Success)",
		ExpDate:     expDate,
		Institution: "Ottopay", // sementara
		DateTime:    utils.GetTimeFormatYYMMDDHHMMSS(),
	}

	err2 := db.DbCon.Create(&labelInq1).Error
	if err2 != nil {
		logs.Info("Failed Save to database", err2)
		// return err1
	}

	// ===== Payment OttoAG =====
	logs.Info("[PAYMENT-BILLER][START]")
	// refnum := utils.GetRrn()

	// payment to ottoag
	billerReq := ottoagmodels.OttoAGPaymentReq{
		Amount:      uint64(inqRespOttoag.Amount),
		CustID:      req.CustID,
		MemberID:    MemberID,
		Period:      req.CustID2,
		Productcode: req.ProductCode,
		Rrn:         inqRespOttoag.Rrn,
	}

	billerRes := biller.PaymentBiller(billerReq, req, dataToken, inqRespOttoag.Amount, inqRespOttoag.Rrn, MemberID, namaVoucher, expDate, category)

	if billerRes.Rc == "09" {
		// save to DB transaski_redeem
		labelPyment1 := dbmodels.TransaksiRedeem{
			AccountNumber: dataToken.Data,
			Voucher:       namaVoucher,
			CustID:        req.CustID,
			// MerchantID:    dataToken.Data.MerchantID,
			RRN:         inqRespOttoag.Rrn,
			ProductCode: req.ProductCode,
			Amount:      inqRespOttoag.Amount,
			TransType:   "Payment",
			Status:      "09 (Pending)",
			ExpDate:     expDate,
			Institution: "Ottopay",
			ProductType: "Pulsa",
			DateTime:    utils.GetTimeFormatYYMMDDHHMMSS(),
		}
		err1 := db.DbCon.Create(&labelPyment1).Error
		if err1 != nil {
			logs.Info("Failed Save to database", err1)
			// return err1
		}

		res = models.UseRedeemResponse{
			Rc:  "09",
			Msg: "Request in progress",
		}
		return res
	}

	if billerRes.Rc != "00" {

		// save to DB transaski_redeem
		labelPyment1 := dbmodels.TransaksiRedeem{
			AccountNumber: dataToken.Data,
			Voucher:       namaVoucher,
			CustID:        req.CustID,
			// MerchantID:    dataToken.Data.MerchantID,
			RRN:         inqRespOttoag.Rrn,
			ProductCode: req.ProductCode,
			Amount:      inqRespOttoag.Amount,
			TransType:   "Payment",
			Status:      "01 (Gagal)",
			ExpDate:     expDate,
			Institution: "Ottopay",
			ProductType: "Pulsa",
			DateTime:    utils.GetTimeFormatYYMMDDHHMMSS(),
		}
		err1 := db.DbCon.Create(&labelPyment1).Error
		if err1 != nil {
			logs.Info("Failed Save to database", err1)
			// return err1
		}

		res = models.UseRedeemResponse{
			Rc:  "01",
			Msg: "Payment Failed",
		}
		return res
	}

	res = models.UseRedeemResponse{
		Rc:  "00",
		Msg: "Payment Success",
	}

	// save to DB transaski_redeem
	labelPyment1 := dbmodels.TransaksiRedeem{
		AccountNumber: dataToken.Data,
		Voucher:       namaVoucher,
		CustID:        req.CustID,
		// MerchantID:    dataToken.Data.MerchantID,
		RRN:         inqRespOttoag.Rrn,
		ProductCode: req.ProductCode,
		Amount:      inqRespOttoag.Amount,
		TransType:   "Payment",
		Status:      "00 (Success)",
		ExpDate:     expDate,
		Institution: "Ottopay",
		ProductType: "Pulsa",
		DateTime:    utils.GetTimeFormatYYMMDDHHMMSS(),
	}
	err1 := db.DbCon.Create(&labelPyment1).Error
	if err1 != nil {
		logs.Info("Failed Save to database", err1)
		// return err1
	}

	if billerRes.Rc != "00" && billerRes.Rc != "09" {
		res = models.UseRedeemResponse{
			Rc:  "01",
			Msg: "Payment Failed",
		}

		return res
	}

	res = models.UseRedeemResponse{
		Rc:          billerRes.Rc,
		Rrn:         billerRes.Rrn,
		CustID:      billerReq.CustID,
		ProductCode: billerReq.Productcode,
		Amount:      int64(billerRes.Amount),
		Msg:         "SUCCESS",
		Uimsg:       "SUCCESS",
		Data:        billerRes.Data,
		Datetime:    utils.GetTimeFormatYYMMDDHHMMSS(),
	}

	return res
}

func ValidatePrefix(OperatorCode int, custID, productCode string) bool {

	logs.Info("===== Req.Product =====", productCode)
	// no, _ := strconv.Atoi(custID)
	logs.Info("===== NO =====", custID)
	prefix := utils.Operator(OperatorCode)
	logs.Info("===== Prefix =====", prefix)
	product := utils.ProductPulsa(productCode[0:4])
	logs.Info("===== Product =====", product)

	logs.Info("===== CustID =====", custID)
	if len(custID) < 4 {

		logs.Info("[FAILED]-[Prefix-ottopoint]-[services-RedeemPulsa]")
		logs.Info(fmt.Sprintf("invalid Prefix %s", custID))

		return false
	}

	// Jika nomor kurang dari 11
	if len(custID) <= 10 || len(custID) > 15 {

		logs.Info("[FAILED]-[Prefix-ottopoint]-[services-RedeemPulsa]")
		logs.Info(fmt.Sprintf("invalid Prefix %s", custID))

		return false

	}

	// Jika Nomor tidak sesuai dengan operator
	if prefix != product {

		logs.Info("[FAILED]-[Prefix-ottopoint]-[services-RedeemPulsa]")
		logs.Info(fmt.Sprintf("invalid Prefix %s", prefix))

		return false

	}

	return true
}
