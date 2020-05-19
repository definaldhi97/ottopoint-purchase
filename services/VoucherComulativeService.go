package services

import (
	"fmt"
	"ottopoint-purchase/db"
	"ottopoint-purchase/hosts/opl/host"
	ottomart "ottopoint-purchase/hosts/ottomart/host"
	ottomartmodels "ottopoint-purchase/hosts/ottomart/models"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"
	"strconv"
	"strings"
	"sync"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type VoucherComulativeService struct {
	General models.GeneralModel
}

func (t VoucherComulativeService) VoucherComulative(req models.VoucherComultaiveReq, param models.Params) models.Response {
	var res models.Response

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[VoucherComulative-Services]",
		zap.String("NameVoucher : ", param.NamaVoucher), zap.Int("Jumlah : ", req.Jumlah),
		zap.String("CampaignID : ", req.CampaignID), zap.String("CustID : ", req.CustID),
		zap.String("CustID2 : ", req.CustID2), zap.String("ProductCode : ", param.ProductCode),
		// zap.Int("Point : ", reiq.Pont),
		zap.String("AccountNumber : ", param.AccountNumber), zap.String("InstitutionID : ", param.InstitutionID))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[RedeemVoucher]")
	defer span.Finish()

	Message_Comulative := ""
	Code_RC_Comulative := ""

	wg := sync.WaitGroup{}

	getResp := models.RedeemComuResp{}
	getResRedeem := models.RedeemResponse{}

	/*---- generate comulative_ref ----*/
	comulative_ref := utils.GenTransactionId()
	param.Reffnum = comulative_ref

	var inqGagal int
	for i := 0; i < req.Jumlah; i++ {
		// wg.Add(2)

		param.Total = i + 1
		getRespChan := make(chan models.RedeemComuResp)
		getErrChan := make(chan error)
		getRespUseVouChan := make(chan models.RedeemResponse)
		getRespUseVoucErr := make(chan error)

		go RedeemComulativeVoucher(req, param, getRespChan, getErrChan)

		if getErr := <-getErrChan; getErr != nil {
			getResp = <-getRespChan
			fmt.Println("Gagal Redeem or Inquiry")
			fmt.Println("Error Message : ", getResp.Message)
			inqGagal++
			continue
		} else {
			fmt.Println("Berhasil Redeem or Inquiry")
			getResp = <-getRespChan
		}

		fmt.Println("test Redeem or Inquiry")

		fmt.Println("Code : ", getResp.Code)
		if getResp.Code == "00" {
			wg.Add(1)
			// req.Jumlah = i
			go UseVoucherComulative(req, getResp, param, getRespUseVouChan, getRespUseVoucErr, &wg)
			getResRedeem = <-getRespUseVouChan

		}
	}
	wg.Wait()

	fmt.Println("Response OttoAG Payment 2 : ", getResRedeem)
	total := req.Jumlah * 2
	countPayment, _ := db.GetCountPyenment(comulative_ref)
	if countPayment.Count != total {
		countPayment, _ = db.GetCountPyenment(comulative_ref)
	}

	countPending, _ := db.GetCountPending_Pyenment(comulative_ref)
	if countPending.Count == 0 {
		countPending, _ = db.GetCountPending_Pyenment(comulative_ref)
	}

	// countFailed, _ := db.GetCountFailedPyenment(comulative_ref)
	// if countFailed.Count == 0 {
	// 	countFailed, _ = db.GetCountFailedPyenment(comulative_ref)
	// }

	countSuccess, _ := db.GetCountSucc_Pyenment(comulative_ref)
	if countSuccess.Count == 0 {
		countSuccess, _ = db.GetCountSucc_Pyenment(comulative_ref)
	}

	pyenmentFail := req.Jumlah - countSuccess.Count

	/* ------ Reversal to Point ----- */
	rcUseVoucher, _ := db.GetPyenmentFailed(comulative_ref)
	fmt.Println("get RC use voucher : ", rcUseVoucher)
	if rcUseVoucher.AccountNumber != "" {
		fmt.Println("============= Reversal to Point ===========")
		// get Custid from user where acount nomor
		dataUser, _ := db.CheckUser(rcUseVoucher.AccountNumber)
		fmt.Println("CustId user OPL from tb user : ", dataUser.CustID)
		Text := "OP009 - " + "Reversal point cause transaction " + param.NamaVoucher + " is failed"
		_, errReversal := host.TransferPoint(dataUser.CustID, strconv.Itoa(rcUseVoucher.Count), Text)
		if errReversal != nil {
			fmt.Println("[ERROR DB]")
			fmt.Println("[transfer point] : ", errReversal)
		}

		fmt.Println("========== Send Notif ==========")
		notifReq := ottomartmodels.NotifRequest{
			AccountNumber:    rcUseVoucher.AccountNumber,
			Title:            "Reversal Point",
			Message:          fmt.Sprintf("Point anda berhasil di reversal sebesar %v", int64(rcUseVoucher.Count)),
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

	/* ---- Message ----
	* Sukses  ( success == jumlah request )
	* Sukses sebagian (success != jumlah request)
	* Gagal (success == 0)
	 */

	fmt.Println(" jumlah transaction payment : ", countPayment.Count)
	fmt.Println(" jumlah success transaction success : ", countSuccess.Count)
	fmt.Println(" jumlah success transaction Pending : ", countPending.Count)
	fmt.Println(" jumlah success transaction failed : ", pyenmentFail)
	fmt.Println(" jumlah request : ", req.Jumlah)
	fmt.Println(" category : ", param.Category)

	respMessage := models.CommulativeResp{
		Success: countSuccess.Count,
		Pending: countPending.Count,
		Failed:  pyenmentFail,
	}

	// Sukses
	if (respMessage.Success != 0) && (respMessage.Pending == 0) && (respMessage.Failed == 0) {
		Code_RC_Comulative = "00"
		Message_Comulative = "Transaksi Berhasil"
	}

	// Sukses & Gagal
	if (respMessage.Success != 0) && (respMessage.Pending == 0) && (respMessage.Failed != 0) {
		Code_RC_Comulative = "174"
		Message_Comulative = fmt.Sprintf("%v Voucher Anda berhasil dirukar namun %v voucher tidak berhasil. Poin yang tidak digunakan akan dikembalikan ke saldo Anda", countSuccess.Count, pyenmentFail)

	}

	// Sukses & Pending
	if (respMessage.Success != 0) && (respMessage.Pending != 0) && (respMessage.Failed == 0) {
		Code_RC_Comulative = "175"
		Message_Comulative = fmt.Sprintf("%v Voucher Anda berhasil ditukar & %v Transaksi Anda sedang dalam proses", countSuccess.Count, countPending.Count)
	}

	// Sukses & Pending & Gagal
	if (respMessage.Success != 0) && (respMessage.Pending != 0) && (respMessage.Failed != 0) {
		Code_RC_Comulative = "33"
		Message_Comulative = fmt.Sprintf("%v vVucher Anda berhasil ditukar namun %v Voucher pending dan %v voucher tidak berhasil. Harap hubungi customer support untuk informasi lebih lanjut.", countSuccess.Count, countPending.Count, pyenmentFail)
		// Message_Comulative = fmt.Sprintf("%v Voucher Anda berhasil dirukar namun %v voucher tidak berhasil. Poin yang tidak digunakan akan dikembalikan ke saldo Anda", countSuccess.Count, pyenmentFail)
	}

	// Pending
	if (respMessage.Success == 0) && (respMessage.Pending != 0) && (respMessage.Failed == 0) {
		Code_RC_Comulative = "56"
		Message_Comulative = fmt.Sprintf("%v Transaksi Anda sedang dalam proses. Silahkan hubungi tim kami untuk informasi selengkapnya.", countPending.Count)
	}

	// Pending & Gagal
	if (respMessage.Success == 0) && (respMessage.Pending != 0) && (respMessage.Failed != 0) {
		Code_RC_Comulative = "57"
		Message_Comulative = fmt.Sprintf("%v Transaksi Anda sedang dalam proses & %v Transaksi Anda Gagal.Poin yang tidak digunakan akan dikembalikan ke saldo Anda", countSuccess.Count, pyenmentFail)

		// Message_Comulative = fmt.Sprintf("%v Transaksi Anda sedang dalam proses & %v Poin yang tidak digunakan akan dikembalikan ke saldo Anda", countSuccess.Count, pyenmentFail)
	}

	// Gagal
	if (respMessage.Success == 0) && (respMessage.Pending == 0) && (respMessage.Failed != 0) {
		Code_RC_Comulative = "01"
		Message_Comulative = "Transaksi Gagal"
	}

	if req.Jumlah == 1 {

		if getResRedeem.Rc == "" {

			getmsg, errmsg := db.GetResponseOttoag("OTTOAG", getResp.Redeem.Rc)

			Code_RC_Comulative = getmsg.InternalRc
			Message_Comulative = getmsg.InternalRd

			if errmsg != nil || getmsg.InternalRc == "" {

				fmt.Println("[VoucherComulativeService]-[GetResponseOttoag]")
				fmt.Println("[Failed to Get Data Mapping Response]")
				fmt.Println(fmt.Sprintf("[Data GetResponseOttoag : ]", getmsg))
				fmt.Println(fmt.Sprintf("[Error %v]", errmsg))
				// return res, err

				Code_RC_Comulative = getResp.Redeem.Rc
				Message_Comulative = getResp.Redeem.Msg

			}

		} else {

			getmsg, errmsg := db.GetResponseOttoag("OTTOAG", getResRedeem.Rc)

			Code_RC_Comulative = getmsg.InternalRc
			Message_Comulative = getmsg.InternalRd

			if errmsg != nil || getmsg.InternalRc == "" {

				fmt.Println("[VoucherComulativeService]-[GetResponseOttoag]")
				fmt.Println("[Failed to Get Data Mapping Response]")
				fmt.Println(fmt.Sprintf("[Data GetResponseOttoag : ]", getmsg))
				fmt.Println(fmt.Sprintf("[Error %v]", errmsg))
				// return res, err

				Code_RC_Comulative = getResRedeem.Rc
				Message_Comulative = getResRedeem.Msg

			}

		}

	}

	rc := Code_RC_Comulative
	msg := Message_Comulative
	if req.Jumlah > 1 {
		rc, msg = getMsgCummulative(Code_RC_Comulative, Message_Comulative)
	}
	// pyenmentFail := req.Jumlah - countSuccess.Count
	// pyenmentPending := req.Jumlah - countPending.Count

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

func getMsgCummulative(rc, msg string) (string, string) {

	var codeRc, codeMsg string

	getmsg, errmsg := db.GetResponseCummulativeOttoAG(rc)
	if errmsg != nil || getmsg.InternalRc == "" {

		fmt.Println("[VoucherComulativeService]-[GetResponseCummulativeOttoAG]")
		fmt.Println("[Failed to Get Data Mapping Response]")
		fmt.Println(fmt.Sprintf("[Data GetResponseOttoag : ]", getmsg))
		fmt.Println(fmt.Sprintf("[Error %v]", errmsg))
		// return res, err

		codeRc = rc
		codeMsg = msg

		return codeRc, codeMsg
	}

	codeRc = getmsg.InternalRc
	codeMsg = strings.Replace(getmsg.InternalRd, "[x]", "%v", 10)
	// codeMsg = getmsg.InternalRd

	return codeRc, codeMsg
}
