package voucher

import (
	"fmt"
	"ottopoint-purchase/models"
	ottoagmodels "ottopoint-purchase/models/ottoag"
	biller "ottopoint-purchase/services/ottoag"
	"ottopoint-purchase/utils"
)

func RedeemPulsaComulative(req models.UseRedeemRequest, reqOP interface{}, param models.Params) models.UseRedeemResponse {
	res := models.UseRedeemResponse{}

	fmt.Println("[Start]-[Package-Voucher]-[RedeemPulsaComulative]")

	// ===== Payment OttoAG =====
	fmt.Println("[PAYMENT-BILLER][START]")
	// fmt.Println("Param : ", param)
	// refnum := utils.GetRrn()

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

	// fmt.Println("Jumlah Voucher : ", param.Total)
	fmt.Println("Response OttoAG Payment : ", billerRes)
	// fmt.Println("[Response Payment %v]", billerRes.Rc)
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

	fmt.Println("[Payment Response : %v]", billerRes)
	if billerRes.Rc == "09" || billerRes.Rc == "68" {
		fmt.Println("[Payment Pending]")

		go SaveTransactionPulsa(paramPay, billerRes, billerReq, reqOP, "Payment", "09", billerRes.Rc)

		res = models.UseRedeemResponse{
			Rc:  "09",
			Msg: "Request in progress",
		}
		return res
	}

	if billerRes.Rc != "00" && billerRes.Rc != "09" && billerRes.Rc != "68" {
		fmt.Println("[Payment Failed]")

		go SaveTransactionPulsa(paramPay, billerRes, billerReq, reqOP, "Payment", "01", billerRes.Rc)

		res = models.UseRedeemResponse{
			Rc:  "01",
			Msg: "Payment Failed",
		}

		return res
	}

	fmt.Println("[Payment Success]", billerRes.Rc)
	// fmt.Println("[Response Payment %v]", billerRes.Rc)
	go SaveTransactionPulsa(paramPay, billerRes, billerReq, reqOP, "Payment", "00", billerRes.Rc)

	res = models.UseRedeemResponse{
		Rc:          billerRes.Rc,
		Rrn:         billerRes.Rrn,
		Category:    param.Category,
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
