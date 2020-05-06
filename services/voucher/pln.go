package voucher

import (
	"encoding/json"
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	ottomart "ottopoint-purchase/hosts/ottomart/host"
	ottomartmodels "ottopoint-purchase/hosts/ottomart/models"
	"ottopoint-purchase/models"
	"ottopoint-purchase/models/dbmodels"
	ottoagmodels "ottopoint-purchase/models/ottoag"
	"ottopoint-purchase/services/ottoag"
	biller "ottopoint-purchase/services/ottoag"
	"ottopoint-purchase/utils"

	"github.com/astaxie/beego/logs"
)

func RedeemPLN(req models.UseRedeemRequest, reqOP interface{}, param models.Params) models.UseRedeemResponse {
	res := models.UseRedeemResponse{}

	logs.Info("[Start]-[Package-Voucher]-[RedeemPLN]")
	// ===== Inquiry OttoAG =====

	inqBiller := ottoagmodels.BillerInquiryDataReq{
		ProductCode: req.ProductCode,
		MemberID:    utils.MemberID,
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

	// inqRespOttoag := ottoagmodels.OttoAGInquiryResponse{}
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
		ProductType: "PLN",
		ProductCode: req.ProductCode,
		Category:    "PLN",
		Point:       param.Point,
		ExpDate:     param.ExpDate,
		SupplierID:  param.SupplierID,
	}

	if inqRespOttoag.Rc != "00" {

		logs.Info("[Error-InquiryResponse]-[RedeemPLN]")
		logs.Info("[Error : %v]", errinqRespOttoag)

		go SaveTransactionPLN(paramInq, inqRespOttoag, inqBiller, reqOP, "Inquiry", "01", inqRespOttoag.Rc)

		res = models.UseRedeemResponse{
			Rc:  "01",
			Msg: "Inquiry Failed",
		}

		return res
	}

	logs.Info("[Response Inquiry %v]", inqRespOttoag.Rc)
	go SaveTransactionPLN(paramInq, inqRespOttoag, inqBiller, reqOP, "Inquiry", "00", inqRespOttoag.Rc)

	// ===== Payment OttoAG =====
	logs.Info("[PAYMENT-BILLER][START]")

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
		// Reffnum
		RRN:         billerRes.Rrn,
		Amount:      int64(billerRes.Amount),
		NamaVoucher: param.NamaVoucher,
		ProductType: "PLN",
		ProductCode: req.ProductCode,
		Category:    "PLN",
		Point:       param.Point,
		ExpDate:     param.ExpDate,
		SupplierID:  param.SupplierID,
	}

	if billerRes.Rc == "09" || billerRes.Rc == "68" {
		logs.Info("[Response Payment %v]", billerRes.Rc)

		go SaveTransactionPLN(paramPay, billerRes, billerReq, reqOP, "Payment", "09", billerRes.Rc)

		res = models.UseRedeemResponse{
			Rc:  "09",
			Msg: "Request in progress",
		}
		return res
	}

	if billerRes.Rc != "00" && billerRes.Rc != "09" && billerRes.Rc != "68" {
		logs.Info("[Response Payment %v]", billerRes.Rc)

		go SaveTransactionPLN(paramPay, billerRes, billerReq, reqOP, "Payment", "01", billerRes.Rc)

		res = models.UseRedeemResponse{
			Rc:  "01",
			Msg: "Payment Failed",
		}

		return res
	}

	// Format Token
	stroomToken := utils.GetFormattedToken(billerRes.Data.Tokenno)

	notifReq := ottomartmodels.NotifRequest{
		AccountNumber:    req.AccountNumber,
		Title:            "Transaksi Berhasil",
		Message:          fmt.Sprintf("Mitra OttoPay, transaksi pembelian voucher PLN telah berhasil. Silakan masukan kode berikut %v ke meteran listrik kamu. Nilai kwh yang diperoleh sesuai dengan PLN. Terima kasih.", stroomToken),
		NotificationType: 3,
	}

	// send notif & inbox
	dataNotif, errNotif := ottomart.NotifAndInbox(notifReq)
	if errNotif != nil {
		logs.Info("Error to send Notif & Inbox")
	}

	if dataNotif.RC != "00" {
		logs.Info("[Response Notif PLN]")
		logs.Info("Gagal Send Notif & Inbox")
		logs.Info("Error : ", errNotif)
	}

	logs.Info("[Response Payment %v]", billerRes.Rc)

	go SaveTransactionPLN(paramPay, billerRes, billerReq, reqOP, "Payment", "00", billerRes.Rc)

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

func SaveTransactionPLN(param models.Params, res interface{}, reqdata interface{}, reqOP interface{}, trasnType, status, rc string) {

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

	reqOttoag, _ := json.Marshal(&reqdata)
	responseOttoag, _ := json.Marshal(&res)
	reqdataOP, _ := json.Marshal(&reqOP)

	save := dbmodels.TransaksiRedeem{
		AccountNumber:   param.AccountNumber,
		Voucher:         param.NamaVoucher,
		MerchantID:      param.MerchantID,
		CustID:          param.CustID,
		RRN:             param.RRN,
		ProductCode:     param.ProductCode,
		Amount:          int64(param.Amount),
		TransType:       trasnType,
		ProductType:     "PLN",
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
		logs.Info("[Package-Voucher]-[Service-RedeemPLN]")
		// return err

	}
}
