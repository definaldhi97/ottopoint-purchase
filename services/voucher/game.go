package voucher

import (
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

func RedeemGame(req models.UseRedeemRequest, dataToken redismodels.TokenResp, MemberID, namaVoucher, expDate, category string) models.UseRedeemResponse {
	res := models.UseRedeemResponse{}

	// ===== Inquiry OttoAG =====

	inqBiller := ottoagmodels.BillerInquiryDataReq{
		ProductCode: req.ProductCode,
		MemberID:    MemberID,
		CustID:      req.CustID,
		Period:      req.CustID2,
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
	inqRespOttoag = biller.InquiryBiller(inqReq, req, dataToken, MemberID, namaVoucher, expDate)

	if inqRespOttoag.Rc != "00" {
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

	// billerRes = models.OttoAGPaymentRes{
	// 	Data: models.DataPayPulsa{},
	// }

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
		// AccountNumber: dataToken.Data,
		Voucher: namaVoucher,
		CustID:  req.CustID,
		// // MerchantID:    dataToken.Data.MerchantID,
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
		CustID:      billerRes.Custid,
		CustID2:     billerRes.Period,
		ProductCode: billerRes.Productcode,
		Amount:      int64(billerRes.Amount),
		Msg:         billerRes.Msg,
		Uimsg:       billerRes.Uimsg,
		Data:        billerRes.Data,
		Datetime:    utils.GetTimeFormatYYMMDDHHMMSS(),
	}

	return res
}
