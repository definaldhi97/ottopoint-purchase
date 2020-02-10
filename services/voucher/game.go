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

func RedeemGame(req models.UseRedeemRequest, AccountNumber, InstitutionID, MemberID, namaVoucher, expDate, category string) models.UseRedeemResponse {
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
	inqRespOttoag = biller.InquiryBiller(inqReq.Data, req, AccountNumber, MemberID, namaVoucher, expDate)

	if inqRespOttoag.Rc != "00" {

		logs.Info("[Response Inquiry %v]", inqRespOttoag.Rc)
		go SaveTransactionGame(AccountNumber, namaVoucher, inqRespOttoag.CustID, inqRespOttoag.Periode, inqRespOttoag.Rrn, inqRespOttoag.ProductCode, "Inquiry", "01", InstitutionID, inqRespOttoag.Amount)

		res = models.UseRedeemResponse{
			Rc:  "01",
			Msg: "Inquiry Failed",
		}

		return res
	}

	logs.Info("[Response Inquiry %v]", inqRespOttoag.Rc)
	go SaveTransactionGame(AccountNumber, namaVoucher, inqRespOttoag.CustID, inqRespOttoag.Periode, inqRespOttoag.Rrn, inqRespOttoag.ProductCode, "Inquiry", "01", InstitutionID, inqRespOttoag.Amount)

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

	billerRes := biller.PaymentBiller(billerReq, req, AccountNumber, inqRespOttoag.Amount, inqRespOttoag.Rrn, MemberID, namaVoucher, expDate, category)

	if billerRes.Rc == "09" || billerRes.Rc == "68" {
		logs.Info("[Response Payment %v]", billerRes.Rc)

		go SaveTransactionGame(AccountNumber, namaVoucher, billerRes.Custid, billerRes.Period, billerRes.Rrn, billerRes.Productcode, "Payment", "09", InstitutionID, int64(billerRes.Amount))

		res = models.UseRedeemResponse{
			Rc:  "09",
			Msg: "Request in progress",
		}
		return res
	}

	if billerRes.Rc != "00" && billerRes.Rc != "09" && billerRes.Rc != "68" {
		logs.Info("[Response Payment %v]", billerRes.Rc)

		go SaveTransactionGame(AccountNumber, namaVoucher, billerRes.Custid, billerRes.Period, billerRes.Rrn, billerRes.Productcode, "Payment", "01", InstitutionID, int64(billerRes.Amount))

		res = models.UseRedeemResponse{
			Rc:  "01",
			Msg: "Payment Failed",
		}

		return res
	}

	logs.Info("[Response Payment %v]", billerRes.Rc)
	go SaveTransactionGame(AccountNumber, namaVoucher, billerRes.Custid, billerRes.Period, billerRes.Rrn, billerRes.Productcode, "Payment", "00", InstitutionID, int64(billerRes.Amount))

	res = models.UseRedeemResponse{
		Rc:          billerRes.Rc,
		Rrn:         billerRes.Rrn,
		Category:    "GAME",
		CustID:      billerRes.Custid,
		CustID2:     billerRes.Period,
		ProductCode: billerRes.Productcode,
		Amount:      int64(billerRes.Amount),
		Msg:         "SUCCESS",
		Uimsg:       "SUCCESS",
		Data:        billerRes.Data,
		Datetime:    utils.GetTimeFormatYYMMDDHHMMSS(),
	}

	return res
}

func SaveTransactionGame(AccountNumber, voucher, CustID1, CustID2, RRN, ProductCode, trasnType, status, instituion string, amount int64) {

	logs.Info("[Start-SaveDB]-[Game]")

	var saveStatus string
	switch status {
	case "00":
		saveStatus = constants.Success
	case "09":
		saveStatus = constants.Pending
	case "01":
		saveStatus = constants.Failed
	}

	cust_id := CustID1 + " || " + CustID2

	save := dbmodels.TransaksiRedeem{
		AccountNumber: AccountNumber,
		Voucher:       voucher,
		CustID:        cust_id,
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
		logs.Info("[Package-Voucher]-[Service-RedeemGame]")
		// return err

	}
}
