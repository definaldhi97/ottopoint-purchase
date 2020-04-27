package voucher

import (
	"encoding/json"
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

func RedeemPulsa(req models.UseRedeemRequest, reqOP interface{}, param models.Params) models.UseRedeemResponse {
	res := models.UseRedeemResponse{}

	logs.Info("[Start]-[Package-Voucher]-[RedeemPulsa]")

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
		MemberID:    utils.MemberID,
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

	logs.Info("[INQUIRY-BILLER][START]")
	inqRespOttoag, errinqRespOttoag := biller.InquiryBiller(inqReq.Data, reqOP, req, param)

	paramInq := models.Params{
		AccountNumber: req.AccountNumber,
		MerchantID:    param.MerchantID,
		InstitutionID: param.InstitutionID,
		CustID:        req.CustID,
		// Reffnum
		RRN:         inqRespOttoag.Rrn,
		Amount:      inqRespOttoag.Amount,
		NamaVoucher: param.NamaVoucher,
		ProductType: "Pulsa",
		ProductCode: req.ProductCode,
		Category:    "Pulsa",
		Point:       param.Point,
		ExpDate:     param.ExpDate,
		SupplierID:  param.SupplierID,
	}

	// Validate if RC != 00
	if inqRespOttoag.Rc != "00" {
		logs.Info("[Error-InquiryResponse]-[RedeemPulsa]")
		logs.Info("[Error : %v]", errinqRespOttoag)

		go SaveTransactionPulsa(paramInq, inqRespOttoag, inqBiller, reqOP, "Inquiry", "01", inqRespOttoag.Rc)

		res = models.UseRedeemResponse{
			Rc:  "01",
			Msg: "Inquiry Failed",
		}

		return res
	}

	// Save DB if RC == 00
	logs.Info("[Response Inquiry %v]", inqRespOttoag.Rc)
	go SaveTransactionPulsa(paramInq, inqRespOttoag, inqBiller, reqOP, "Inquiry", "00", inqRespOttoag.Rc)

	// ===== Payment OttoAG =====
	logs.Info("[PAYMENT-BILLER][START]")
	// refnum := utils.GetRrn()

	// payment to ottoag
	billerReq := ottoagmodels.OttoAGPaymentReq{
		Amount:      uint64(inqRespOttoag.Amount),
		CustID:      req.CustID,
		MemberID:    utils.MemberID,
		Period:      req.CustID2,
		Productcode: req.ProductCode,
		Rrn:         inqRespOttoag.Rrn,
	}

	billerRes := biller.PaymentBiller(billerReq, reqOP, req, param)

	paramPay := models.Params{
		AccountNumber: req.AccountNumber,
		MerchantID:    param.MerchantID,
		InstitutionID: param.InstitutionID,
		CustID:        req.CustID,
		Reffnum:       param.Reffnum,
		RRN:           billerRes.Rrn,
		Amount:        int64(billerRes.Amount),
		NamaVoucher:   param.NamaVoucher,
		ProductType:   "Pulsa",
		ProductCode:   req.ProductCode,
		Category:      "Pulsa",
		Point:         param.Point,
		ExpDate:       param.ExpDate,
		SupplierID:    param.SupplierID,
	}

	if billerRes.Rc == "09" || billerRes.Rc == "68" {
		logs.Info("[Response Payment %v]", billerRes.Rc)

		go SaveTransactionPulsa(paramPay, billerRes, billerReq, reqOP, "Payment", "09", billerRes.Rc)

		res = models.UseRedeemResponse{
			Rc:  "09",
			Msg: "Request in progress",
		}
		return res
	}

	if billerRes.Rc != "00" && billerRes.Rc != "09" && billerRes.Rc != "68" {
		logs.Info("[Response Payment %v]", billerRes.Rc)

		go SaveTransactionPulsa(paramPay, billerRes, billerReq, reqOP, "Payment", "01", billerRes.Rc)

		res = models.UseRedeemResponse{
			Rc:  "01",
			Msg: "Payment Failed",
		}

		return res
	}

	logs.Info("[Response Payment %v]", billerRes.Rc)
	go SaveTransactionPulsa(paramPay, billerRes, billerReq, reqOP, "Payment", "00", billerRes.Rc)

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

// func SaveTransactionPulsa(AccountNumber, voucher, CustID, RRN, ProductCode, trasnType, status, instituion string, amount int64) {
func SaveTransactionPulsa(param models.Params, res interface{}, reqdata interface{}, reqOP interface{}, trasnType, status, rc string) {

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

	reqOttoag, _ := json.Marshal(&reqdata)  // Req Ottoag
	responseOttoag, _ := json.Marshal(&res) // Response Ottoag
	reqdataOP, _ := json.Marshal(&reqOP)    // Req Service

	save := dbmodels.TransaksiRedeem{
		AccountNumber:   param.AccountNumber,
		Voucher:         param.NamaVoucher,
		MerchantID:      param.MerchantID,
		CustID:          param.CustID,
		RRN:             param.RRN,
		ProductCode:     param.ProductCode,
		Amount:          int64(param.Amount),
		TransType:       trasnType,
		ProductType:     "Pulsa",
		Status:          saveStatus,
		ExpDate:         param.ExpDate,
		Institution:     param.InstitutionID,
		CummulativeRef:  param.Reffnum,
		DateTime:        utils.GetTimeFormatYYMMDDHHMMSS(),
		ResponderData:   status,
		Point:           param.Point,
		ResponderRc:     rc,
		RequestorData:   string(reqOttoag),
		ResponderData2:  string(responseOttoag),
		RequestorOPData: string(reqdataOP),
		SupplierID:      param.SupplierID,
	}

	err := db.DbCon.Create(&save).Error
	if err != nil {
		logs.Info("[Failed Save to DB ]", err)
		logs.Info("[Package-Voucher]-[Service-RedeemPulsa]")
		// return err

	}
}
