package voucher

import (
	"fmt"
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

// type RedeemPulsaServices struct {
// 	General models.GeneralModel
// }

func RedeemPulsa(req models.UseRedeemRequest, AccountNumber, InstitutionID, MemberID, namaVoucher, expDate, category string) models.UseRedeemResponse {
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
	inqRespOttoag = biller.InquiryBiller(inqReq.Data, req, AccountNumber, MemberID, namaVoucher, expDate)

	// Validate if RC != 00
	if inqRespOttoag.Rc != "00" {
		logs.Info("[Response Inquiry %v]", inqRespOttoag.Rc)

		go SaveTransactionPulsa(AccountNumber, namaVoucher, inqRespOttoag.CustID, inqRespOttoag.Rrn, inqRespOttoag.ProductCode, "Inquiry", "01", InstitutionID, inqRespOttoag.Amount)

		res = models.UseRedeemResponse{
			Rc:  "01",
			Msg: "Inquiry Failed",
		}

		return res
	}

	// Save DB if RC == 00
	logs.Info("[Response Inquiry %v]", inqRespOttoag.Rc)
	go SaveTransactionPulsa(AccountNumber, namaVoucher, inqRespOttoag.CustID, inqRespOttoag.Rrn, inqRespOttoag.ProductCode, "Inquiry", "00", InstitutionID, inqRespOttoag.Amount)

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

		go SaveTransactionPulsa(AccountNumber, namaVoucher, billerRes.Custid, billerRes.Rrn, billerRes.Productcode, "Payment", "09", InstitutionID, int64(billerRes.Amount))

		res = models.UseRedeemResponse{
			Rc:  "09",
			Msg: "Request in progress",
		}
		return res
	}

	if billerRes.Rc != "00" && billerRes.Rc != "09" && billerRes.Rc != "68" {
		logs.Info("[Response Payment %v]", billerRes.Rc)

		go SaveTransactionPulsa(AccountNumber, namaVoucher, billerRes.Custid, billerRes.Rrn, billerRes.Productcode, "Payment", "01", InstitutionID, int64(billerRes.Amount))

		res = models.UseRedeemResponse{
			Rc:  "01",
			Msg: "Payment Failed",
		}

		return res
	}

	logs.Info("[Response Payment %v]", billerRes.Rc)
	go SaveTransactionPulsa(AccountNumber, namaVoucher, billerRes.Custid, billerRes.Rrn, billerRes.Productcode, "Payment", "00", InstitutionID, int64(billerRes.Amount))

	res = models.UseRedeemResponse{
		Rc:          billerRes.Rc,
		Rrn:         billerRes.Rrn,
		Category:    "PULSA",
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

func SaveTransactionPulsa(AccountNumber, voucher, CustID, RRN, ProductCode, trasnType, status, instituion string, amount int64) {

	logs.Info("[Start-SaveDB]-[Pulsa]")

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
		logs.Info("[Package-Voucher]-[Service-RedeemPulsa]")
		// return err

	}
}
