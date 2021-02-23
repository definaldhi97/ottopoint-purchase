package redeemtion

import (
	"encoding/json"
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	ottoag "ottopoint-purchase/hosts/ottoag/host"
	kafka "ottopoint-purchase/hosts/publisher/host"
	"ottopoint-purchase/models"
	"ottopoint-purchase/models/dbmodels"
	ottoagmodels "ottopoint-purchase/models/ottoag"
	"ottopoint-purchase/services"
	"ottopoint-purchase/utils"
	"strconv"
	"strings"
	"sync"
	"time"

	"ottopoint-purchase/services/v2.1/Trx"

	"github.com/sirupsen/logrus"
	"github.com/vjeantet/jodaTime"
)

// func (t V21_VoucherOttoAgService) V21_VoucherOttoAg(req models.VoucherComultaiveReq, param models.Params, header models.RequestHeader) models.Response {
func RedeemtionOttoAG_V21_Service(req models.VoucherComultaiveReq, param models.Params, header models.RequestHeader) models.Response {
	// fmt.Println("[ >>>>>>>>>>>>>>>>>> Voucher OttoAg Service <<<<<<<<<<<<<<<< ]")

	var res models.Response

	nameservice := "[PackageRedeemtion]-[RedeemtionOttoAG_V21_Service]"
	logReq := fmt.Sprintf("[AccountNumber : %v, RewardID : %v]", param.AccountNumber, param.RewardID)

	logrus.Info(nameservice)

	Message_Comulative := ""
	Code_RC_Comulative := ""

	wg := sync.WaitGroup{}

	getResp := models.RedeemComuResp{}
	getResRedeem := models.RedeemResponse{}

	/*---- generate comulative_ref ----*/
	comulative_ref := utils.GenTransactionId()
	param.Reffnum = comulative_ref

	// logrus.Info("[ Cumulatif Reff ] : ", comulative_ref)

	var inqGagal int

	for i := 0; i < req.Jumlah; i++ {

		param.TrxID = utils.GenTransactionId()
		param.Total = i + 1

		getRespChan := make(chan models.RedeemComuResp)
		getErrChan := make(chan error)
		getRespUseVouChan := make(chan models.RedeemResponse)
		getRespUseVoucErr := make(chan error)

		go inquiryOttoAG(req, param, header, getRespChan, getErrChan)
		if getErr := <-getErrChan; getErr != nil {
			getResp = <-getRespChan
			fmt.Println("[ Failed Deduct point, Deduct voucher or Inquiry Voucher ]")
			fmt.Println("Error Message : ", getResp.Message)
			inqGagal++
			continue
		} else {
			fmt.Println("[ Success Deduct point, Deduct voucher and Inquiry Voucher")
			getResp = <-getRespChan
		}

		fmt.Println("[ Response Code RedeemVoucherOttoAg : ", getResp.Code)
		if getResp.Code == "00" {
			wg.Add(1)
			// go V21_UseVoucherOttoAg(req, getResp, param, getRespUseVouChan, getRespUseVoucErr, &wg)
			go useVoucherOttoAG(req, getResp, param, getRespUseVouChan, getRespUseVoucErr, &wg)
			getResRedeem = <-getRespUseVouChan
		}

	}
	wg.Wait()

	fmt.Println("[ Response OttoAG Payment ] : ", getResRedeem)

	countPayment, _ := db.GetCountPyenment(comulative_ref)
	if countPayment.Count != req.Jumlah*2 {
		countPayment, _ = db.GetCountPyenment(comulative_ref)
	}

	countPending, _ := db.GetCountPending_Pyenment(comulative_ref)
	if countPending.Count == 0 {
		countPending, _ = db.GetCountPending_Pyenment(comulative_ref)
	}

	countSuccess, _ := db.GetCountSucc_Pyenment(comulative_ref)
	if countSuccess.Count == 0 {
		countSuccess, _ = db.GetCountSucc_Pyenment(comulative_ref)
	}

	pyenmentFail := req.Jumlah - countSuccess.Count - countPending.Count

	/* ------ Reversal to Point ----- */
	rcUseVoucher, _ := db.GetPyenmentFailed(comulative_ref)
	fmt.Println("[ Get RC Payment T_Spending by TSP02 ] : ", rcUseVoucher)

	if rcUseVoucher.AccountNumber != "" {
		fmt.Println("============= Reversal to Point ===========")

		trxIDReversal := utils.GenTransactionId()
		param.TrxID = trxIDReversal
		resultReversal := Trx.V21_Adding_PointVoucher(param, rcUseVoucher.Count, rcUseVoucher.CountFailed, header)
		fmt.Println(resultReversal)

		fmt.Println("[ >>>>>>>>>>>>>>>>> Send Publisher Notification <<<<<<<<<<<<<<<< ]")

		pubreq := models.NotifPubreq{
			Type:           constants.CODE_REVERSAL_POINT,
			NotificationTo: param.AccountNumber,
			Institution:    param.InstitutionID,
			ReferenceId:    param.RRN,
			TransactionId:  param.Reffnum,
			Data: models.DataValue{
				RewardValue: "point",
				Value:       strconv.Itoa(rcUseVoucher.Count),
			},
		}

		bytePub, _ := json.Marshal(pubreq)

		kafkaReq := kafka.PublishReq{
			Topic: utils.TopicsNotif,
			Value: bytePub,
		}

		kafkaRes, errKafka := kafka.SendPublishKafka(kafkaReq)
		if errKafka != nil {

			logrus.Error(nameservice)
			logrus.Error(fmt.Sprintf("[SendPublishKafka]-[Error : %v]", errKafka))
			logrus.Println(logReq)
		}

		fmt.Println("Response Publisher Notification : ", kafkaRes)

	}

	/* -------------- Message -----------------------
	* Sukses  ( success == jumlah request )
	* Sukses sebagian (success != jumlah request)
	* Gagal (success == 0)
	* -----------------------------------------------
	 */

	fmt.Println(" jumlah transaction payment : ", countPayment.Count)
	fmt.Println(" jumlah success transaction success : ", countSuccess.Count)
	fmt.Println(" jumlah success transaction Pending : ", countPending.Count)
	fmt.Println(" jumlah success transaction failed : ", pyenmentFail)
	fmt.Println(" jumlah request : ", req.Jumlah)
	fmt.Println(" Category : ", param.Category)

	respMessage := models.CommulativeResp{
		Success: countSuccess.Count,
		Pending: countPending.Count,
		Failed:  pyenmentFail,
	}

	var s, p, f int

	// Sukses
	if (respMessage.Success != 0) && (respMessage.Pending == 0) && (respMessage.Failed == 0) {
		Code_RC_Comulative = "00"
		Message_Comulative = "Transaksi Berhasil"

		s = countSuccess.Count
	}

	// Sukses & Gagal
	if (respMessage.Success != 0) && (respMessage.Pending == 0) && (respMessage.Failed != 0) {
		Code_RC_Comulative = "174"
		Message_Comulative = fmt.Sprintf("%v Voucher Anda berhasil dirukar namun %v voucher tidak berhasil. Poin yang tidak digunakan akan dikembalikan ke saldo Anda", countSuccess.Count, pyenmentFail)

		s = countSuccess.Count
		f = pyenmentFail
	}

	// Sukses & Pending
	if (respMessage.Success != 0) && (respMessage.Pending != 0) && (respMessage.Failed == 0) {
		Code_RC_Comulative = "175"
		Message_Comulative = fmt.Sprintf("%v Voucher Anda berhasil ditukar & %v Transaksi Anda sedang dalam proses", countSuccess.Count, countPending.Count)

		s = countSuccess.Count
		p = countPending.Count

	}

	// Sukses & Pending & Gagal
	if (respMessage.Success != 0) && (respMessage.Pending != 0) && (respMessage.Failed != 0) {
		Code_RC_Comulative = "33"
		Message_Comulative = fmt.Sprintf("%v Vucher Anda berhasil ditukar namun %v Voucher pending dan %v voucher tidak berhasil. Harap hubungi customer support untuk informasi lebih lanjut.", countSuccess.Count, countPending.Count, pyenmentFail)
		// Message_Comulative = fmt.Sprintf("%v Voucher Anda berhasil dirukar namun %v voucher tidak berhasil. Poin yang tidak digunakan akan dikembalikan ke saldo Anda", countSuccess.Count, pyenmentFail)

		s = countSuccess.Count
		p = countPending.Count
		f = pyenmentFail
	}

	// Pending
	if (respMessage.Success == 0) && (respMessage.Pending != 0) && (respMessage.Failed == 0) {
		Code_RC_Comulative = "56"
		Message_Comulative = fmt.Sprintf("%v Transaksi Anda sedang dalam proses. Silahkan hubungi tim kami untuk informasi selengkapnya.", countPending.Count)

		p = countPending.Count

	}

	// Pending & Gagal
	if (respMessage.Success == 0) && (respMessage.Pending != 0) && (respMessage.Failed != 0) {
		Code_RC_Comulative = "57"
		Message_Comulative = fmt.Sprintf("%v Transaksi Anda sedang dalam proses & %v Transaksi Anda Gagal.Poin yang tidak digunakan akan dikembalikan ke saldo Anda", countSuccess.Count, pyenmentFail)

		p = countPending.Count
		f = pyenmentFail
	}

	// Gagal
	if (respMessage.Success == 0) && (respMessage.Pending == 0) && (respMessage.Failed != 0) {
		Code_RC_Comulative = "01"
		Message_Comulative = "Transaksi Gagal"

		f = pyenmentFail
	}

	rc := Code_RC_Comulative
	msg := Message_Comulative

	if req.Jumlah == 1 {

		if getResRedeem.Rc == "" {

			getmsg, errmsg := db.GetResponseOttoag("OTTOAG", getResp.Redeem.Rc)

			rc = getmsg.InternalRc
			msg = getmsg.InternalRd

			if errmsg != nil || getmsg.InternalRc == "" {

				logrus.Error(nameservice)
				logrus.Error(fmt.Sprintf("[GetResponseOttoag]-[Error : %v]", errmsg))
				logrus.Error(fmt.Sprintf("[Response %v]", getmsg))
				logrus.Println(logReq)
				// return res, err

				rc = getResp.Redeem.Rc
				msg = getResp.Redeem.Msg

			}

		} else {

			getmsg, errmsg := db.GetResponseOttoag("OTTOAG", getResRedeem.Rc)

			rc = getmsg.InternalRc
			msg = getmsg.InternalRd

			if errmsg != nil || getmsg.InternalRc == "" {

				logrus.Error(nameservice)
				logrus.Error(fmt.Sprintf("[GetResponseOttoag]-[Error : %v]", errmsg))
				logrus.Error(fmt.Sprintf("[Response %v]", getmsg))
				// return res, err

				rc = getResRedeem.Rc
				msg = getResRedeem.Msg

			}

		}

	}

	var m string
	if req.Jumlah > 1 {
		m = services.GetMsgCummulative(rc, msg)
	}

	if s != 0 && f != 0 && p == 0 {
		a := strings.Replace(m, "[x]", fmt.Sprintf("%v", s), 1)
		b := strings.Replace(a, "[x]", fmt.Sprintf("%v", f), 1)

		msg = b
	}

	if s != 0 && f == 0 && p != 0 {
		a := strings.Replace(m, "[x]", fmt.Sprintf("%v", s), 1)
		b := strings.Replace(a, "[x]", fmt.Sprintf("%v", p), 1)

		msg = b
	}

	if s != 0 && f != 0 && p != 0 {
		a := strings.Replace(m, "[x]", fmt.Sprintf("%v", s), 1)
		b := strings.Replace(a, "[x]", fmt.Sprintf("%v", p), 1)
		c := strings.Replace(b, "[x]", fmt.Sprintf("%v", f), 1)

		msg = c
	}

	if s == 0 && f == 0 && p != 0 {
		a := strings.Replace(m, "[x]", fmt.Sprintf("%v", p), 1)
		msg = a
	}

	/* ------ Response UseVoucher Comulative */
	fmt.Println("========== Mesage from Inquiry OTTOAG and OPL ===============")
	fmt.Println("Rc : ", getResp.Code)
	fmt.Println("Message : ", getResp.Message)
	fmt.Println("=============================================================")
	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: models.CommulativeResp{
			Code:    rc,
			Msg:     msg,
			Success: countSuccess.Count,
			Pending: countPending.Count,
			Failed:  pyenmentFail,

			//RedeemRes :
		},
	}

	return res

}

// func V21_RedeemVoucherOttoAg(req models.VoucherComultaiveReq, param models.Params, header models.RequestHeader, getResp chan models.RedeemComuResp, ErrRespRedeem chan error) {
func inquiryOttoAG(req models.VoucherComultaiveReq, param models.Params, header models.RequestHeader, getResp chan models.RedeemComuResp, ErrRespRedeem chan error) {
	fmt.Println("[ >>>>>>>>>>>>>>>>>> Redeemtion Comulative Voucher Otto AG <<<<<<<<<<<<<<<< ]")
	logrus.Info("[ Inquery OttoAG ] - [ Deduct point OPL & Deduct Voucher ]")

	defer close(getResp)
	defer close(ErrRespRedeem)

	resRedeemComu := models.RedeemComuResp{}
	redeemRes := models.RedeemComuResp{
		Code: "00",
	}

	// ==========Inquery OttoAG==========
	inqBiller := ottoagmodels.BillerInquiryDataReq{
		ProductCode: param.ProductCode,
		MemberID:    utils.MemberID,
		CustID:      req.CustID,
		Period:      req.CustID2,
	}

	inqReq := ottoagmodels.OttoAGInquiryRequest{
		TypeTrans:     "0003",
		Datetime:      utils.GetTimeFormatYYMMDDHHMMSS(),
		IssuerID:      "OTTOPAY",
		AccountNumber: param.AccountNumber,
		Data:          inqBiller,
	}

	reqInq := models.UseRedeemRequest{
		AccountNumber: param.AccountNumber,
		CustID:        req.CustID,
		CustID2:       req.CustID2,
		ProductCode:   param.ProductCode,
	}

	fmt.Println("[INQUIRY-BILLER][START]")

	resOttAG := ottoagmodels.OttoAGInquiryResponse{}
	headOttoAg := ottoag.PackMessageHeader(inqReq.Data)
	billerDataHost, errInquiry := ottoag.Send(inqReq.Data, headOttoAg, "INQUIRY")
	errInquiry = json.Unmarshal(billerDataHost, &resOttAG)

	// if errInquiry = json.Unmarshal(billerDataHost, &resOttAG); errInquiry != nil {
	// 	logrus.Info("[INQUIRY-SERVICES-01]")
	// 	logs.Error("Failed to unmarshaling json response from ottoag", errInquiry)
	// 	resOttAG = ottoagmodels.OttoAGInquiryResponse{
	// 		Rc:  "01",
	// 		Msg: "Inquiry Failed",
	// 	}

	// 	return resOttAG, errInquiry
	// }

	// resOttAG, errInquiry := biller.InquiryBiller(inqReq.Data, req, reqInq, param)

	textCommentSpending := param.TrxID + "#" + param.NamaVoucher
	param.Comment = textCommentSpending
	paramInq := models.Params{
		PaymentMethod: param.PaymentMethod,
		AccountNumber: param.AccountNumber,
		MerchantID:    param.MerchantID,
		InstitutionID: param.InstitutionID,
		CustID:        req.CustID,
		TransType:     constants.CODE_TRANSTYPE_INQUERY,
		Reffnum:       param.Reffnum, // internal
		RRN:           resOttAG.Rrn,
		TrxID:         param.TrxID,
		Amount:        resOttAG.Amount,
		NamaVoucher:   param.NamaVoucher,
		ProductType:   param.ProductType,
		ProductCode:   param.ProductCode,
		Category:      param.Category,
		Point:         param.Point,
		ExpDate:       param.ExpDate,
		SupplierID:    param.SupplierID,
		CategoryID:    param.CategoryID,
		CampaignID:    param.CampaignID,
		ProductID:     param.ProductID,
		AccountId:     param.AccountId,
		Comment:       textCommentSpending,
		RewardID:      param.RewardID,
		DataSupplier: models.Supplier{
			Rc: resOttAG.Rc,
			Rd: resOttAG.Msg,
		},
	}

	if resOttAG.Rc != constants.CODE_SUCCESS || errInquiry != nil {
		fmt.Println("[Error-DataInquiry]-[Redeem Comulative Voucher Otto AG]")
		fmt.Println("[Error : %v]", errInquiry)

		redeemRes = models.RedeemComuResp{
			Code:    "01",
			Message: "Inquiry Failed",
		}

		go saveTransactionOttoAg(paramInq, resOttAG, reqInq, req, constants.CODE_FAILED)

		ErrRespRedeem <- errInquiry

		r := models.RedeemResponse{
			Rc:          resOttAG.Rc,
			Rrn:         resOttAG.Rrn,
			CustID:      resOttAG.CustID,
			ProductCode: resOttAG.ProductCode,
			Amount:      resOttAG.Amount,
			Msg:         resOttAG.Msg,
			Uimsg:       resOttAG.Uimsg,
			// Datetime:    time.Now(),
			Data: resOttAG.Data,
		}

		resRedeemComu.Redeem = r
		resRedeemComu.Code = redeemRes.Code
		resRedeemComu.Message = redeemRes.Message

		getResp <- resRedeemComu

		return

	}

	// Time Out
	if resOttAG.Rc == "" {
		fmt.Println("[Error-DataInquiry]-[Redeem Comulative Voucher Otto AG]")
		fmt.Println("[Error : %v]", errInquiry)
		redeemRes = models.RedeemComuResp{
			Code:    "01",
			Message: "Inquiry Failed",
		}

		go saveTransactionOttoAg(paramInq, resOttAG, reqInq, req, constants.CODE_FAILED)

		ErrRespRedeem <- errInquiry

		resRedeemComu.Redeem.Rc = "01"
		resRedeemComu.Redeem.Rc = "Time Out"
		resRedeemComu.Code = redeemRes.Code
		resRedeemComu.Message = redeemRes.Message

		getResp <- resRedeemComu

		return

	}
	//ss
	// spending point and spending usage_limit voucher
	resultRedeemVouch, errRedeemVouch := Trx.V21_Redeem_PointandVoucher(1, param, header)
	fmt.Println("Response Deduct point dan voucher")
	fmt.Println(resultRedeemVouch)

	// paramInq.CouponID = resultRedeemVouch.CouponseVouch[0].CouponsID

	if resultRedeemVouch.Rc == "00" {
		paramInq.CouponID = resultRedeemVouch.CouponseVouch[0].CouponsID
		paramInq.PointTransferID = resultRedeemVouch.PointTransferID
	}

	go saveTransactionOttoAg(paramInq, resOttAG, reqInq, req, constants.CODE_SUCCESS)

	if resultRedeemVouch.Rc != "00" {
		logrus.Error("[ Error Redeem_PointandVoucher] : ", resultRedeemVouch.Rc)
		logrus.Error("[ Error Redeem_PointandVoucher] : ", resultRedeemVouch.Rd)

		redeemRes = models.RedeemComuResp{
			Code:    "01",
			Message: "Gagal Redeem",
		}

		ErrRespRedeem <- errRedeemVouch

		resRedeemComu.Redeem.Rc = "01"
		resRedeemComu.Redeem.Msg = "Maaf transaksi Anda tidak dapat dilakukan saat ini. Silahkan dicoba lagi atau hubungi tim kami untuk informasi selengkapnya"
		resRedeemComu.Code = redeemRes.Code
		resRedeemComu.Message = redeemRes.Message

		getResp <- resRedeemComu

		return
	}

	resRedeemComu.CouponCode = resultRedeemVouch.CouponseVouch[0].CouponsCode
	resRedeemComu.CouponID = resultRedeemVouch.CouponseVouch[0].CouponsID

	ErrRespRedeem <- nil

	r := models.RedeemResponse{
		Rc:          resOttAG.Rc,
		Rrn:         resOttAG.Rrn,
		CustID:      resOttAG.CustID,
		ProductCode: resOttAG.ProductCode,
		Amount:      resOttAG.Amount,
		Msg:         resOttAG.Msg,
		Uimsg:       resOttAG.Uimsg,
		// Datetime:    time.Now(),
		Data: resOttAG.Data,
	}

	resRedeemComu.Code = redeemRes.Code
	resRedeemComu.Message = redeemRes.Message
	resRedeemComu.PointTransferID = resultRedeemVouch.PointTransferID
	resRedeemComu.Comment = textCommentSpending
	resRedeemComu.Redeem = r
	getResp <- resRedeemComu

}

func useVoucherOttoAG(req models.VoucherComultaiveReq, redeemComu models.RedeemComuResp, param models.Params, getRespChan chan models.RedeemResponse, ErrRespUseVouc chan error, wg *sync.WaitGroup) {
	// fmt.Println("[ >>>>>>>>>>>>>>>>>> Use Voucher Comulative Voucher Otto AG <<<<<<<<<<<<<<<< ]")

	defer wg.Done()
	defer close(getRespChan)
	defer close(ErrRespUseVouc)

	// Reedem Use Voucher
	param.Amount = redeemComu.Redeem.Amount
	param.RRN = redeemComu.Redeem.Rrn
	param.CouponID = redeemComu.CouponID
	param.PointTransferID = redeemComu.PointTransferID
	param.Comment = redeemComu.Comment
	// resRedeem := services.RedeemUseVoucherComulative(req, param)

	res := models.RedeemResponse{}

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
		Productcode: param.ProductCode,
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

	billerRes := ottoagmodels.OttoAGPaymentRes{}

	billerHead := ottoag.PackMessageHeader(billerReq)
	logrus.Info("Nama Voucher : ", param.NamaVoucher)
	billerDataHost, errPayment := ottoag.Send(billerReq, billerHead, "PAYMENT")
	errPayment = json.Unmarshal(billerDataHost, &billerRes)

	fmt.Println(fmt.Sprintf("Response OttoAG %v Payment : %v", param.ProductType, billerRes))
	paramPay := models.Params{
		PaymentMethod:   param.PaymentMethod,
		AccountNumber:   param.AccountNumber,
		MerchantID:      param.MerchantID,
		InstitutionID:   param.InstitutionID,
		CustID:          custId,
		TransType:       constants.CODE_TRANSTYPE_REDEMPTION,
		Reffnum:         param.Reffnum, // Internal
		RRN:             billerRes.Rrn,
		TrxID:           param.TrxID,
		Amount:          int64(billerRes.Amount),
		NamaVoucher:     param.NamaVoucher,
		ProductType:     param.ProductType,
		ProductCode:     param.ProductCode,
		Category:        param.Category,
		Point:           param.Point,
		ExpDate:         param.ExpDate,
		SupplierID:      param.SupplierID,
		CategoryID:      param.CategoryID,
		CampaignID:      param.CampaignID,
		VoucherCode:     billerRes.Data.Code,
		CouponID:        param.CouponID,
		ExpireDateVidio: billerRes.Data.EndDateVidio,
		AccountId:       param.AccountId,
		ProductID:       param.ProductID,
		RewardID:        param.RewardID,
		PointTransferID: param.PointTransferID,
		Comment:         param.Comment,
		DataSupplier: models.Supplier{
			Rc: billerRes.Rc,
			Rd: billerRes.Msg,
		},
	}

	fmt.Println(fmt.Sprintf("[Payment Response : %v]", billerRes))

	// Time Out
	if billerRes.Rc == "" || errPayment != nil {
		fmt.Println(fmt.Sprintf("[Payment %v Time Out]", param.ProductType))

		save := saveTransactionOttoAg(paramPay, billerRes, billerReq, req, "09")

		fmt.Println(fmt.Sprintf("[Response Save Payment Pulsa : %v]", save))

		res = models.RedeemResponse{
			// Rc:  "09",
			// Msg: "Request in progress",
			Rc:    "68",
			Msg:   "Timeout",
			Uimsg: "Timeout",
		}

		getRespChan <- res

		return
	}

	// Pending
	if billerRes.Rc == "09" || billerRes.Rc == "68" {
		fmt.Println(fmt.Sprintf("[Payment %v Pending]", param.ProductType))

		save := saveTransactionOttoAg(paramPay, billerRes, billerReq, req, "09")
		fmt.Println(fmt.Sprintf("[Response Save Payment Pulsa : %v]", save))

		res = models.RedeemResponse{
			// Rc:  "09",
			// Msg: "Request in progress",
			Rc:    billerRes.Rc,
			Msg:   billerRes.Msg,
			Uimsg: "Request in progress",
		}
		getRespChan <- res

		return
	}

	// Gagal
	if billerRes.Rc != "00" && billerRes.Rc != "09" && billerRes.Rc != "68" {
		fmt.Println(fmt.Sprintf("[Payment %v Failed]", param.ProductType))

		save := saveTransactionOttoAg(paramPay, billerRes, billerReq, req, "01")
		fmt.Println(fmt.Sprintf("[Response Save Payment Pulsa : %v]", save))

		res = models.RedeemResponse{
			// Rc:  "01",
			// Msg: "Payment Failed",
			Rc:    billerRes.Rc,
			Msg:   billerRes.Msg,
			Uimsg: "Payment Failed",
		}

		getRespChan <- res

		return
	}

	// Notif PLN

	if param.Category == constants.CategoryPLN {

		// Format Token
		stroomToken := utils.GetFormattedToken(billerRes.Data.Tokenno)
		denom := strconv.Itoa(billerRes.Data.Amount)

		fmt.Println("data denom : ", denom)
		paramPay.VoucherCode = stroomToken
		// swtich notif app/sms
		dtaIssuer, _ := db.GetDataInstitution(param.InstitutionID)
		if dtaIssuer.NOtificationID == constants.CODE_SMS_NOTIF || dtaIssuer.NOtificationID == constants.CODE_SMS_APPS_NOTIF {
			fmt.Println("SMS Notif : ", param.Category)
			fmt.Println("Institution : ", param.InstitutionID)
			fmt.Println("Notification ID : ", dtaIssuer.NOtificationID)
			fmt.Println("========== Send Publisher ==========")
			pubreqSMSNotif := []models.NotifPubreq{}
			a := models.NotifPubreq{
				Type:           constants.CODE_REDEEM_PLN_SMS,
				NotificationTo: param.AccountNumber,
				Institution:    param.InstitutionID,
				ReferenceId:    param.RRN,
				TransactionId:  param.Reffnum,
				Data: models.DataValueSMS{
					ProductName: denom,
					Token:       stroomToken,
				},
			}
			pubreqSMSNotif = append(pubreqSMSNotif, a)

			go sendToPublisher(pubreqSMSNotif, utils.TopicNotifSMS)
		}
		if dtaIssuer.NOtificationID == constants.CODE_APPS_NOTIF || dtaIssuer.NOtificationID == constants.CODE_SMS_APPS_NOTIF {
			fmt.Println("APP Notif : ", param.Category)
			fmt.Println("Institution : ", param.InstitutionID)
			fmt.Println("Notification ID : ", dtaIssuer.NOtificationID)
			fmt.Println("========== Send Publisher ==========")

			pubreq := models.NotifPubreq{
				Type:           constants.CODE_REDEEM_PLN,
				NotificationTo: param.AccountNumber,
				Institution:    param.InstitutionID,
				ReferenceId:    param.RRN,
				TransactionId:  param.Reffnum,
				Data: models.DataValue{
					RewardValue: denom,
					Value:       stroomToken,
				},
			}
			go sendToPublisher(pubreq, utils.TopicsNotif)
		}

	}

	// Notif Vidio
	if param.Category == constants.CategoryVidio {

		denom := strconv.FormatUint(billerRes.Amount, 10)
		fmt.Println("APP Notif : ", param.Category)
		fmt.Println("Institution : ", param.InstitutionID)
		fmt.Println("data denom : ", denom)
		fmt.Println("========== Send Publisher ==========")

		pubreq := models.NotifPubreq{
			Type:           constants.CODE_REDEEM_VIDIO,
			NotificationTo: param.AccountNumber,
			Institution:    param.InstitutionID,
			ReferenceId:    param.RRN,
			TransactionId:  param.Reffnum,
			Data: models.DataValue{
				RewardValue: denom,
				Value:       billerRes.Data.Code,
			},
		}
		go sendToPublisher(pubreq, utils.TopicsNotif)
	}

	fmt.Println(fmt.Sprintf("[Payment %v Success]", param.ProductType))
	save := saveTransactionOttoAg(paramPay, billerRes, billerReq, req, "00")
	fmt.Println(fmt.Sprintf("[Response Save Payment %v : %v]", param.ProductType, save))

	res = models.RedeemResponse{
		Rc:  billerRes.Rc,
		Rrn: billerRes.Rrn,
		// Category:    param.Category,
		CustID:      billerReq.CustID,
		ProductCode: billerReq.Productcode,
		Amount:      int64(billerRes.Amount),
		Msg:         billerRes.Msg,
		Uimsg:       "SUCCESS",
		Data:        billerRes.Data,
		Datetime:    utils.GetTimeFormatYYMMDDHHMMSS(),
	}

	getRespChan <- res

	return
}

func sendToPublisher(pubreq interface{}, topic string) {

	bytePub, _ := json.Marshal(pubreq)

	kafkaReq := kafka.PublishReq{
		Topic: topic,
		Value: bytePub,
	}

	kafkaRes, err := kafka.SendPublishKafka(kafkaReq)
	if err != nil {
		fmt.Println("Gagal Send Publisher")
		fmt.Println("Error : ", err)
	}

	fmt.Println("Response Publisher : ", kafkaRes)
}

func saveTransactionOttoAg(param models.Params, res interface{}, reqdata interface{}, req interface{}, status string) string {

	fmt.Println(fmt.Sprintf("[Start-SaveDB]-[%v]", param.ProductType))

	// validasi vidio is_used -> false
	isUsed := true
	// codeVoucher := param.VoucherCode
	var codeVoucher string
	var ExpireDate time.Time
	var redeemDate time.Time
	var usedAt time.Time
	trxID := utils.GenTransactionId()

	if param.TransType == constants.CODE_TRANSTYPE_REDEMPTION {
		// timeRedeem := jodaTime.Format("dd-MM-YYYY HH:mm:ss", time.Now())
		// redeemDate = timeRedeem
		codeVoucher = utils.EncryptVoucherCode(param.VoucherCode, param.CouponID)
		isUsed = true
		ExpireDate = utils.ExpireDateVoucherAGt(constants.EXPDATE_VOUCHER)
		redeemDate = time.Now()
		trxID = param.TrxID
		usedAt = time.Now()

		if param.Category == constants.CategoryVidio {
			isUsed = false // isUsed status untuk used
			usedAt = time.Time{}
		}

	}

	var saveStatus string
	switch status {
	case "00":
		saveStatus = constants.Success
	case "09":
		saveStatus = constants.Pending
	case "01":
		saveStatus = constants.Failed
		isUsed = true
	}

	reqOttoag, _ := json.Marshal(&reqdata)  // Req Ottoag
	responseOttoag, _ := json.Marshal(&res) // Response Ottoag
	reqdataOP, _ := json.Marshal(&req)      // Req Service

	if param.PaymentMethod == constants.SplitBillMethod {

		_, errUpdate := db.UpdateTransactionSplitBill(isUsed, param.TrxID, saveStatus, param.DataSupplier.Rc, param.DataSupplier.Rd, responseOttoag, reqOttoag, reqdataOP)
		if errUpdate != nil {

			logrus.Error(fmt.Sprintf("[UpdateTransactionSplitBill]-[SaveTransactionOttoAg]"))
			logrus.Error(fmt.Sprintf("[TrxID : %v]-[Error : %v]", trxID, errUpdate))

			return fmt.Sprintf("Gagal Updated " + trxID)
		}

		return fmt.Sprintf("Berhasil Updated " + trxID)

	}

	save := dbmodels.TSpending{
		ID:            utils.GenerateTokenUUID(),
		AccountNumber: param.AccountNumber,
		Voucher:       param.NamaVoucher,
		MerchantID:    param.MerchantID,
		CustID:        param.CustID,
		RRN:           param.RRN,
		// TransactionId: param.TrxID,
		TransactionId: trxID,
		ProductCode:   param.ProductCode,
		Amount:        int64(param.Amount),
		TransType:     param.TransType,
		// IsUsed:          true,
		IsUsed:      isUsed,
		ProductType: param.ProductType,
		Status:      saveStatus,
		// ExpDate:         param.ExpDate,
		ExpDate:           utils.DefaultNulTime(ExpireDate),
		Institution:       param.InstitutionID,
		CummulativeRef:    param.Reffnum,
		DateTime:          utils.GetTimeFormatYYMMDDHHMMSS(),
		Point:             param.Point,
		ResponderRc:       param.DataSupplier.Rc,
		ResponderRd:       param.DataSupplier.Rd,
		RequestorData:     string(reqOttoag),
		ResponderData:     string(responseOttoag),
		RequestorOPData:   string(reqdataOP),
		SupplierID:        param.SupplierID,
		RedeemAt:          utils.DefaultNulTime(redeemDate),
		CampaignId:        param.CampaignID,
		VoucherCode:       codeVoucher,
		CouponId:          param.CouponID,
		AccountId:         param.AccountId,
		ProductCategoryID: param.CategoryID,
		Comment:           param.Comment,
		MRewardID:         param.RewardID,
		MProductID:        param.ProductID,
		PointsTransferID:  param.PointTransferID,
		UsedAt:            utils.DefaultNulTime(usedAt),
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
