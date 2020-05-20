package voucher

import (
	"fmt"
	ottomart "ottopoint-purchase/hosts/ottomart/host"
	ottomartmodels "ottopoint-purchase/hosts/ottomart/models"
	"ottopoint-purchase/models"
	ottoagmodels "ottopoint-purchase/models/ottoag"
	biller "ottopoint-purchase/services/ottoag"
	"ottopoint-purchase/utils"
)

func RedeemPLNComulative(req models.UseRedeemRequest, reqOP interface{}, param models.Params) models.UseRedeemResponse {
	res := models.UseRedeemResponse{}

	fmt.Println("[Start]-[Package-Voucher]-[RedeemPLNComulative]")

	// ===== Payment OttoAG =====
	fmt.Println("[PAYMENT-BILLER][START]")

	// payment to ottoag
	billerReq := ottoagmodels.OttoAGPaymentReq{
		Amount:   uint64(param.Amount),
		CustID:   req.CustID,
		MemberID: utils.MemberID,
		// Period:      req.CustID2,
		Productcode: req.ProductCode,
		Rrn:         param.RRN,
	}

	billerRes := biller.PaymentBiller(billerReq, reqOP, req, param)

	paramPay := models.Params{
		AccountNumber: param.AccountNumber,
		MerchantID:    param.MerchantID,
		InstitutionID: param.InstitutionID,
		CustID:        req.CustID,
		Reffnum:       param.Reffnum, // Internal
		RRN:           billerRes.Rrn,
		Amount:        int64(billerRes.Amount),
		NamaVoucher:   param.NamaVoucher,
		ProductType:   param.ProductType,
		ProductCode:   req.ProductCode,
		Category:      param.Category,
		Point:         param.Point,
		ExpDate:       param.ExpDate,
		SupplierID:    param.SupplierID,
	}

	fmt.Println(fmt.Sprintf("[Payment Response : %v]", billerRes))

	// Time Out
	if billerRes.Rc == "" {
		fmt.Println("[Payment Time Out]")

		go SaveTransactionGame(paramPay, billerRes, billerReq, reqOP, "Payment", "09", billerRes.Rc)

		res = models.UseRedeemResponse{
			// Rc:  "09",
			// Msg: "Request in progress",
			Rc:    billerRes.Rc,
			Msg:   "Timeout",
			Uimsg: "Timeout",
		}

		return res
	}

	// Pending
	if billerRes.Rc == "09" || billerRes.Rc == "68" || billerRes.Rc == "" {
		fmt.Println("[Payment Pending]")

		go SaveTransactionGame(paramPay, billerRes, billerReq, reqOP, "Payment", "09", billerRes.Rc)

		res = models.UseRedeemResponse{
			// Rc:  "09",
			// Msg: "Request in progress",
			Rc:    billerRes.Rc,
			Msg:   billerRes.Msg,
			Uimsg: "Request in progress",
		}

		return res
	}

	// Gagal
	if billerRes.Rc != "00" && billerRes.Rc != "09" && billerRes.Rc != "68" {
		fmt.Println("[Payment Failed]")

		go SaveTransactionGame(paramPay, billerRes, billerReq, reqOP, "Payment", "01", billerRes.Rc)

		res = models.UseRedeemResponse{
			// Rc:  "01",
			// Msg: "Payment Failed",
			Rc:    billerRes.Rc,
			Msg:   billerRes.Msg,
			Uimsg: "Payment Failed",
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
		fmt.Println("Error to send Notif & Inbox")
	}

	if dataNotif.RC != "00" {
		fmt.Println("[Response Notif PLN]")
		fmt.Println("Gagal Send Notif & Inbox")
		fmt.Println("Error : ", errNotif)
	}

	fmt.Println("[Payment Success]")

	go SaveTransactionPLN(paramPay, billerRes, billerReq, reqOP, "Payment", "00", billerRes.Rc)

	res = models.UseRedeemResponse{
		Rc:          billerRes.Rc,
		Rrn:         billerRes.Rrn,
		Category:    param.Category,
		CustID:      billerRes.Custid,
		Amount:      int64(billerRes.Amount),
		ProductCode: billerRes.Productcode,
		Msg:         billerRes.Msg,
		Uimsg:       "SUCCESS",
		Data:        billerRes.Data,
		Datetime:    utils.GetTimeFormatYYMMDDHHMMSS(),
	}

	return res
}
