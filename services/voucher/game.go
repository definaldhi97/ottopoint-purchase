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

func RedeemGame(req models.UseRedeemRequest, reqOP interface{}, param models.Params) models.UseRedeemResponse {
	res := models.UseRedeemResponse{}

	logs.Info("[Start]-[Package-Voucher]-[RedeemGame]")

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

	var custIDGame string
	if req.CustID2 != "" {
		custIDGame = req.CustID + " || " + req.CustID2
	} else {
		custIDGame = req.CustID
	}

	// custIDGame := req.CustID + " || " + req.CustID2

	paramInq := models.Params{
		AccountNumber: req.AccountNumber,
		MerchantID:    param.MerchantID,
		InstitutionID: param.InstitutionID,
		CustID:        custIDGame,
		// Reffnum
		RRN:         inqRespOttoag.Rrn,
		Amount:      inqRespOttoag.Amount,
		NamaVoucher: param.NamaVoucher,
		ProductType: "Game",
		ProductCode: req.ProductCode,
		Category:    "Game",
		Point:       param.Point,
		ExpDate:     param.ExpDate,
		SupplierID:  param.SupplierID,
	}

	if inqRespOttoag.Rc != "00" {

		logs.Info("[Error-InquiryResponse]-[RedeemGame]")
		logs.Info("[Error : %v]", errinqRespOttoag)

		go SaveTransactionGame(paramInq, "Inquiry", "01")

		res = models.UseRedeemResponse{
			Rc:  "01",
			Msg: "Inquiry Failed",
		}

		return res
	}

	logs.Info("[Response Inquiry %v]", inqRespOttoag.Rc)
	go SaveTransactionGame(paramInq, "Inquiry", "00")

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
		CustID:        custIDGame,
		// Reffnum
		RRN:         billerRes.Rrn,
		Amount:      int64(billerRes.Amount),
		NamaVoucher: param.NamaVoucher,
		ProductType: "Game",
		ProductCode: req.ProductCode,
		Category:    "Game",
		Point:       param.Point,
		ExpDate:     param.ExpDate,
		SupplierID:  param.SupplierID,
	}

	if billerRes.Rc == "09" || billerRes.Rc == "68" {
		logs.Info("[Response Payment %v]", billerRes.Rc)

		go SaveTransactionGame(paramPay, "Payment", "09")

		res = models.UseRedeemResponse{
			Rc:  "09",
			Msg: "Request in progress",
		}
		return res
	}

	if billerRes.Rc != "00" && billerRes.Rc != "09" && billerRes.Rc != "68" {
		logs.Info("[Response Payment %v]", billerRes.Rc)

		go SaveTransactionGame(paramPay, "Payment", "01")

		res = models.UseRedeemResponse{
			Rc:  "01",
			Msg: "Payment Failed",
		}

		return res
	}

	logs.Info("[Response Payment %v]", billerRes.Rc)
	go SaveTransactionGame(paramPay, "Payment", "00")

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

func SaveTransactionGame(param models.Params, trasnType, status string) {

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

	save := dbmodels.TransaksiRedeem{
		AccountNumber: param.AccountNumber,
		Voucher:       param.NamaVoucher,
		CustID:        param.CustID,
		MerchantID:    param.MerchantID,
		RRN:           param.RRN,
		ProductCode:   param.ProductCode,
		Amount:        param.Amount,
		TransType:     trasnType,
		Status:        saveStatus,
		ExpDate:       param.ExpDate,
		Institution:   param.InstitutionID,
		ProductType:   param.ProductType,
		DateTime:      utils.GetTimeFormatYYMMDDHHMMSS(),
	}

	err := db.DbCon.Create(&save).Error
	if err != nil {
		logs.Info("[Failed Save to DB ]", err)
		logs.Info("[Package-Voucher]-[Service-RedeemGame]")
		// return err

	}
}
