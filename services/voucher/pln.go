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

func RedeemPLN(req models.UseRedeemRequest, dataToken redismodels.TokenResp, MemberID, namaVoucher, expDate, category string) models.UseRedeemResponse {
	res := models.UseRedeemResponse{}

	logs.Info("[Start]-[Package-Voucher]-[RedeemPLN]")
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

	// // Format Token
	// stroomToken := utils.GetFormattedToken(billerRes.Data.Tokenno)

	// Format Struct notif
	// notifReq := ottomartmodels.NotifReq{
	// 	CollapseKey: "type_c",
	// 	Title:       "Transaksi Berhasil",
	// 	Body:        fmt.Sprintf("Mitra OttoPay, transaksi pembelian voucher PLN telah berhasil. Silakan masukan kode berikut %v ke meteran listrik kamu. Nilai kwh yang diperoleh sesuai dengan PLN. Terima kasih.", stroomToken),
	// 	Target:      "earning point",
	// 	Phone:       "",
	// 	Rc:          "00",
	// }

	// // send notif & inbox
	// _, errNotif := ottomart.NotifInboxOttomart(notifReq, header)
	// if errNotif != nil {
	// 	logs.Info("Error to send Notif & Inbox")
	// }

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
	errDB := db.DbCon.Create(&labelPyment1).Error
	if errDB != nil {
		logs.Info("Failed Save to database", errDB)
	}

	res = models.UseRedeemResponse{
		Rc:          billerRes.Rc,
		Rrn:         billerRes.Rrn,
		CustID:      billerRes.Custid,
		Amount:      int64(billerRes.Amount),
		ProductCode: billerRes.Productcode,
		Msg:         billerRes.Msg,
		Uimsg:       billerRes.Uimsg,
		Data:        billerRes.Data,
		Datetime:    utils.GetTimeFormatYYMMDDHHMMSS(),
	}

	return res
}
