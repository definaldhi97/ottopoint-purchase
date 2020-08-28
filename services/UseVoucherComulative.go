package services

import (
	"encoding/json"
	"fmt"
	"ottopoint-purchase/constants"
	db "ottopoint-purchase/db"
	opl "ottopoint-purchase/hosts/opl/host"
	ottomart "ottopoint-purchase/hosts/ottomart/host"
	ottomartmodels "ottopoint-purchase/hosts/ottomart/models"
	"ottopoint-purchase/models"
	"ottopoint-purchase/models/dbmodels"
	ottoagmodels "ottopoint-purchase/models/ottoag"
	biller "ottopoint-purchase/services/ottoag"
	"ottopoint-purchase/utils"
	"sync"
	"time"

	"github.com/vjeantet/jodaTime"
)

// func UseVoucherService1(header models.RequestHeader, req models.VoucherComultaiveReq, dataToken ottomartmodels.ResponseToken, amount int64, rrn string) models.Response {
func UseVoucherComulative(req models.VoucherComultaiveReq, redeemComu models.RedeemComuResp, param models.Params, getRespChan chan models.RedeemResponse, ErrRespUseVouc chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(getRespChan)
	defer close(ErrRespUseVouc)
	// resRedeemComu := models.RedeemComuResp{}
	fmt.Println("[UseVoucherComulative]-[Package-Services]")

	// get CustID
	dataUser, errUser := db.CheckUser(param.AccountNumber)
	if errUser != nil {
		fmt.Println("User Belum Eligible, Error : ", errUser)
	} else {
		// Use Voucher to Openloyalty
		fmt.Println("campaignId : ", req.CampaignID)
		fmt.Println("couponId : ", redeemComu.CouponID)
		fmt.Println("code : ", redeemComu.CouponCode)
		fmt.Println("used : ", 1)
		fmt.Println("customerId : ", dataUser.CustID)

		// _, err2 := opl.CouponVoucherCustomer(req.CampaignID, redeemComu.CouponID, redeemComu.CouponCode, dataUser.CustID, 1)
		// fmt.Println("================ doing use voucher couponId : ", redeemComu.CouponID)
		// if err2 != nil {
		// 	fmt.Println("================ doing use voucher couponId Error: ", redeemComu.CouponID)

		// } else {

		// 	// Reedem Use Voucher
		// 	param.Amount = redeemComu.Redeem.Amount
		// 	param.RRN = redeemComu.Redeem.Rrn
		// 	param.CouponID = redeemComu.CouponID
		// 	resRedeem := RedeemUseVoucherComulative(req, param)

		// 	getRespChan <- resRedeem
		// }

		if param.Category != "vidio" {
			_, err2 := opl.CouponVoucherCustomer(req.CampaignID, redeemComu.CouponID, redeemComu.CouponCode, dataUser.CustID, 1)
			fmt.Println("================ doing use voucher couponId : ", redeemComu.CouponID)
			if err2 != nil {
				fmt.Println("================ doing use voucher couponId Error: ", redeemComu.CouponID)
			}
		}

		// Reedem Use Voucher
		param.Amount = redeemComu.Redeem.Amount
		param.RRN = redeemComu.Redeem.Rrn
		param.CouponID = redeemComu.CouponID
		resRedeem := RedeemUseVoucherComulative(req, param)

		getRespChan <- resRedeem

	}

}

// function reedem use voucher
func RedeemUseVoucherComulative(req models.VoucherComultaiveReq, param models.Params) models.RedeemResponse {
	res := models.RedeemResponse{}

	fmt.Println("[RedeemUseVoucherComulative]-[Package-Services]")

	reqRedeem := models.UseRedeemRequest{
		AccountNumber: param.AccountNumber,
		CustID:        req.CustID,
		CustID2:       req.CustID2,
		ProductCode:   param.ProductCode,
		Jumlah:        param.Total,
	}

	// resRedeem := models.UseRedeemResponse{}
	param.CampaignID = req.CampaignID
	resRedeem := PaymentVoucherOttoAg(reqRedeem, req, param)

	// switch category {
	// case constants.CategoryPulsa:
	// 	resRedeem = voucher.RedeemPulsaComulative(reqRedeem, req, param)
	// case constants.CategoryPLN:
	// 	resRedeem = voucher.RedeemPLNComulative(reqRedeem, req, param)
	// case constants.CategoryMobileLegend, constants.CategoryFreeFire:
	// 	resRedeem = voucher.RedeemGameComulative(reqRedeem, req, param)
	// }

	res = models.RedeemResponse{
		Rc:          resRedeem.Rc,
		Rrn:         resRedeem.Rrn,
		CustID:      resRedeem.CustID,
		ProductCode: resRedeem.ProductCode,
		Amount:      resRedeem.Amount,
		Msg:         resRedeem.Msg,
		Uimsg:       resRedeem.Uimsg,
		Datetime:    resRedeem.Datetime,
		Data:        resRedeem.Data,
	}

	fmt.Println("[Test vidio")
	fmt.Println(res.Data)

	fmt.Println("data vidio")
	fmt.Println("StartDate : " + resRedeem.Data.StartDateVidio)
	fmt.Println("EndDate : " + resRedeem.Data.EndDateVidio)
	fmt.Println("code : " + resRedeem.Data.Code)
	fmt.Println("description : " + resRedeem.Data.Description)

	return res

}

func PaymentVoucherOttoAg(req models.UseRedeemRequest, reqOP interface{}, param models.Params) models.UseRedeemResponse {
	res := models.UseRedeemResponse{}
	// ===== Payment OttoAG =====
	fmt.Println(fmt.Sprintf("[PAYMENT-%v][START]", param.ProductType))
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

	var custId string

	if param.Category == constants.CategoryGame {
		if req.CustID2 != "" {
			custId = req.CustID + " || " + req.CustID2
		} else {
			custId = req.CustID
		}
	} else {
		custId = req.CustID
	}

	billerRes := biller.PaymentBiller(billerReq, reqOP, req, param)

	fmt.Println(fmt.Sprintf("Response OttoAG %v Payment : %v", param.ProductType, billerRes))
	paramPay := models.Params{
		AccountNumber: param.AccountNumber,
		MerchantID:    param.MerchantID,
		InstitutionID: param.InstitutionID,
		CustID:        custId,
		TransType:     constants.CODE_TRANSTYPE_REDEMPTION,
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

		CampaignID:      param.CampaignID,
		VoucherCode:     billerRes.Data.Code,
		CouponID:        param.CouponID,
		ExpireDateVidio: billerRes.Data.EndDateVidio,
		DataSupplier: models.Supplier{
			Rc: billerRes.Rc,
			Rd: billerRes.Msg,
		},
	}

	fmt.Println(fmt.Sprintf("[Payment Response : %v]", billerRes))

	// Time Out
	if billerRes.Rc == "" {
		fmt.Println(fmt.Sprintf("[Payment %v Time Out]", param.ProductType))

		save := saveTransactionOttoAg(paramPay, billerRes, billerReq, reqOP, "09")
		fmt.Println(fmt.Sprintf("[Response Save Payment Pulsa : %v]", save))

		res = models.UseRedeemResponse{
			// Rc:  "09",
			// Msg: "Request in progress",
			Rc:    "68",
			Msg:   "Timeout",
			Uimsg: "Timeout",
		}
		return res
	}

	// Pending
	if billerRes.Rc == "09" || billerRes.Rc == "68" {
		fmt.Println(fmt.Sprintf("[Payment %v Pending]", param.ProductType))

		save := saveTransactionOttoAg(paramPay, billerRes, billerReq, reqOP, "09")
		fmt.Println(fmt.Sprintf("[Response Save Payment Pulsa : %v]", save))

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
		fmt.Println(fmt.Sprintf("[Payment %v Failed]", param.ProductType))

		save := saveTransactionOttoAg(paramPay, billerRes, billerReq, reqOP, "01")
		fmt.Println(fmt.Sprintf("[Response Save Payment Pulsa : %v]", save))

		res = models.UseRedeemResponse{
			// Rc:  "01",
			// Msg: "Payment Failed",
			Rc:    billerRes.Rc,
			Msg:   billerRes.Msg,
			Uimsg: "Payment Failed",
		}

		return res
	}

	// Notif PLN
	if param.Category == constants.CategoryPLN {
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

	}

	fmt.Println(fmt.Sprintf("[Payment %v Success]", param.ProductType))
	save := saveTransactionOttoAg(paramPay, billerRes, billerReq, reqOP, "00")
	fmt.Println(fmt.Sprintf("[Response Save Payment %v : %v]", param.ProductType, save))

	res = models.UseRedeemResponse{
		Rc:          billerRes.Rc,
		Rrn:         billerRes.Rrn,
		Category:    param.Category,
		CustID:      billerReq.CustID,
		ProductCode: billerReq.Productcode,
		Amount:      int64(billerRes.Amount),
		Msg:         billerRes.Msg,
		Uimsg:       "SUCCESS",
		Data:        billerRes.Data,
		Datetime:    utils.GetTimeFormatYYMMDDHHMMSS(),
	}

	return res
}

func saveTransactionOttoAg(param models.Params, res interface{}, reqdata interface{}, reqOP interface{}, status string) string {

	fmt.Println(fmt.Sprintf("[Start-SaveDB]-[%v]", param.ProductType))

	// validasi vidio is_used -> false
	isUsed := true
	// codeVoucher := param.VoucherCode
	var codeVoucher string
	ExpireDate := param.ExpDate
	var redeemDate string

	if param.TransType == constants.CODE_TRANSTYPE_REDEMPTION {
		timeRedeem := jodaTime.Format("dd-MM-YYYY HH:mm:ss", time.Now())
		redeemDate = timeRedeem

		codeVoucher = EncryptVoucherCode(param.VoucherCode, param.CouponID)
	}

	if param.Category == "vidio" && param.TransType == constants.CODE_TRANSTYPE_REDEMPTION {
		isUsed = false // isUsed status untuk used

		// a := []rune(param.CouponID)
		// key32 := string(a[0:32])
		// screetKey := []byte(key32)
		// codeVidio := []byte(param.VoucherCode)
		// chiperText, _ := utils.EncryptAES(codeVidio, screetKey)
		// codeVoucher = string(chiperText)
	}

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
		AccountNumber: param.AccountNumber,
		Voucher:       param.NamaVoucher,
		MerchantID:    param.MerchantID,
		CustID:        param.CustID,
		RRN:           param.RRN,
		ProductCode:   param.ProductCode,
		Amount:        int64(param.Amount),
		TransType:     param.TransType,
		// IsUsed:        true,
		IsUsed:      isUsed,
		ProductType: param.ProductType,
		Status:      saveStatus,
		// ExpDate:         param.ExpDate,
		ExpDate:         ExpireDate,
		Institution:     param.InstitutionID,
		CummulativeRef:  param.Reffnum,
		DateTime:        utils.GetTimeFormatYYMMDDHHMMSS(),
		ResponderData:   status,
		Point:           param.Point,
		ResponderRc:     param.DataSupplier.Rc,
		ResponderRd:     param.DataSupplier.Rd,
		RequestorData:   string(reqOttoag),
		ResponderData2:  string(responseOttoag),
		RequestorOPData: string(reqdataOP),
		SupplierID:      param.SupplierID,
		RedeemAt:        redeemDate,
		CampaignId:      param.CampaignID,
		VoucherCode:     codeVoucher,
		CouponId:        param.CouponID,
	}

	err := db.DbCon.Create(&save).Error
	if err != nil {

		fmt.Println(fmt.Sprintf("[Error : %v]", err))
		fmt.Println("[Failed saveTransactionOttoAg to DB]")
		fmt.Println(fmt.Sprintf("[TransType : %v || RRN : %v]", param.TransType, param.RRN))

		name := jodaTime.Format("dd-MM-YYYY", time.Now()) + ".csv"
		go utils.CreateCSVFile(save, name)

		return "Gagal Save"

	}

	return "Berhasil Save"
}

func EncryptVoucherCode(data, key string) string {

	var codeVoucher string
	if data == "" {
		return codeVoucher
	}

	a := []rune(key)
	key32 := string(a[0:32])
	screetKey := []byte(key32)
	codeByte := []byte(data)
	chiperText, _ := utils.EncryptAES(codeByte, screetKey)
	codeVoucher = string(chiperText)
	return codeVoucher
}
