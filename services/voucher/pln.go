package voucher

import (
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	"ottopoint-purchase/models"
	"ottopoint-purchase/models/dbmodels"
	ottoagmodels "ottopoint-purchase/models/ottoag"
	"ottopoint-purchase/services/ottoag"
	biller "ottopoint-purchase/services/ottoag"
	"ottopoint-purchase/utils"

	"github.com/astaxie/beego/logs"
)

func RedeemPLN(req models.UseRedeemRequest, AccountNumber, InstitutionID, MemberID, namaVoucher, expDate, category string) models.UseRedeemResponse {
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
	inqRespOttoag = biller.InquiryBiller(inqReq.Data, req, AccountNumber, MemberID, namaVoucher, expDate)

	if inqRespOttoag.Rc != "00" {

		logs.Info("[Response Inquiry %v]", inqRespOttoag.Rc)
		go SaveTransactionPLN(AccountNumber, namaVoucher, inqRespOttoag.CustID, inqRespOttoag.Rrn, inqRespOttoag.ProductCode, "Inquiry", "01", InstitutionID, inqRespOttoag.Amount)

		res = models.UseRedeemResponse{
			Rc:  "01",
			Msg: "Inquiry Failed",
		}

		return res
	}

	logs.Info("[Response Inquiry %v]", inqRespOttoag.Rc)
	go SaveTransactionPLN(AccountNumber, namaVoucher, inqRespOttoag.CustID, inqRespOttoag.Rrn, inqRespOttoag.ProductCode, "Inquiry", "00", InstitutionID, inqRespOttoag.Amount)

	// ===== Payment OttoAG =====
	logs.Info("[PAYMENT-BILLER][START]")

	// payment to ottoag
	billerReq := ottoagmodels.OttoAGPaymentReq{
		Amount:      uint64(inqRespOttoag.Amount),
		CustID:      req.CustID,
		MemberID:    MemberID,
		Period:      req.CustID2,
		Productcode: req.ProductCode,
		Rrn:         inqRespOttoag.Rrn,
	}

	billerRes := biller.PaymentBiller(billerReq, req, AccountNumber, inqRespOttoag.Amount, inqRespOttoag.Rrn, MemberID, namaVoucher, expDate, category)

	if billerRes.Rc == "09" || billerRes.Rc == "68" {
		logs.Info("[Response Payment %v]", billerRes.Rc)

		go SaveTransactionPLN(AccountNumber, namaVoucher, billerRes.Custid, billerRes.Rrn, billerRes.Productcode, "Payment", "09", InstitutionID, int64(billerRes.Amount))

		res = models.UseRedeemResponse{
			Rc:  "09",
			Msg: "Request in progress",
		}
		return res
	}

	if billerRes.Rc != "00" && billerRes.Rc != "09" && billerRes.Rc != "68" {
		logs.Info("[Response Payment %v]", billerRes.Rc)

		go SaveTransactionPLN(AccountNumber, namaVoucher, billerRes.Custid, billerRes.Rrn, billerRes.Productcode, "Payment", "01", InstitutionID, int64(billerRes.Amount))

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

	logs.Info("[Response Payment %v]", billerRes.Rc)

	go SaveTransactionPLN(AccountNumber, namaVoucher, billerRes.Custid, billerRes.Rrn, billerRes.Productcode, "Payment", "00", InstitutionID, int64(billerRes.Amount))

	res = models.UseRedeemResponse{
		Rc:          billerRes.Rc,
		Rrn:         billerRes.Rrn,
		Category:    "PLN",
		CustID:      billerRes.Custid,
		Amount:      int64(billerRes.Amount),
		ProductCode: billerRes.Productcode,
		Msg:         "SUCCESS",
		Uimsg:       "SUCCESS",
		Data:        billerRes.Data,
		Datetime:    utils.GetTimeFormatYYMMDDHHMMSS(),
	}

	return res
}

func SaveTransactionPLN(AccountNumber, voucher, CustID, RRN, ProductCode, trasnType, status, instituion string, amount int64) {

	logs.Info("[Start-SaveDB]-[PLN]")

	var saveStatus string
	switch status {
	case "00":
		saveStatus = constants.Success
	case "09":
		saveStatus = constants.Pending
	case "01":
		saveStatus = constants.Failed
	}

	save := dbmodels.TransaksiRedeem{
		AccountNumber: AccountNumber,
		Voucher:       voucher,
		CustID:        CustID,
		// MerchantID:    AccountNumber.MerchantID,
		RRN:         RRN,
		ProductCode: ProductCode,
		Amount:      amount,
		TransType:   trasnType,
		Status:      saveStatus,
		// ExpDate:     expDate,
		Institution: instituion,
		ProductType: "Pulsa",
		DateTime:    utils.GetTimeFormatYYMMDDHHMMSS(),
	}

	err := db.DbCon.Create(&save).Error
	if err != nil {
		logs.Info("[Failed Save to DB ]", err)
		logs.Info("[Package-Voucher]-[Service-RedeemPLN]")
		// return err

	}
}
