package voucher

import (
	"ottopoint-purchase/models"
	ottoagmodels "ottopoint-purchase/models/ottoag"
	biller "ottopoint-purchase/services/ottoag"
	"ottopoint-purchase/utils"

	"github.com/astaxie/beego/logs"
)

func RedeemPLNComulative(req models.UseRedeemRequest, reqOP interface{}, param models.Params) models.UseRedeemResponse {
	res := models.UseRedeemResponse{}

	logs.Info("[Start]-[Package-Voucher]-[RedeemPLNComulative]")

	// ===== Payment OttoAG =====
	logs.Info("[PAYMENT-BILLER][START]")

	// payment to ottoag
	billerReq := ottoagmodels.OttoAGPaymentReq{
		Amount:      uint64(param.Amount),
		CustID:      req.CustID,
		MemberID:    utils.MemberID,
		Period:      req.CustID2,
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

	if billerRes.Rc == "09" || billerRes.Rc == "68" {
		logs.Info("[Response Payment %v]", billerRes.Rc)

		go SaveTransactionPLN(paramPay, "Payment", "09")

		res = models.UseRedeemResponse{
			Rc:  "09",
			Msg: "Request in progress",
		}
		return res
	}

	if billerRes.Rc != "00" && billerRes.Rc != "09" && billerRes.Rc != "68" {
		logs.Info("[Response Payment %v]", billerRes.Rc)

		go SaveTransactionPLN(paramPay, "Payment", "01")

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

	go SaveTransactionPLN(paramPay, "Payment", "00")

	res = models.UseRedeemResponse{
		Rc:          billerRes.Rc,
		Rrn:         billerRes.Rrn,
		Category:    param.Category,
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
