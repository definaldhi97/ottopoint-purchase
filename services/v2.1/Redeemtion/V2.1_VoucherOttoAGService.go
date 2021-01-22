package Redeemtion

import (
	"encoding/json"
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	kafka "ottopoint-purchase/hosts/publisher/host"
	"ottopoint-purchase/models"
	"ottopoint-purchase/services"
	"ottopoint-purchase/services/v2.1/Trx"
	"ottopoint-purchase/utils"
	"strconv"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type V21_VoucherOttoAgService struct {
	General models.GeneralModel
}

func (t V21_VoucherOttoAgService) V21_VoucherOttoAg(req models.VoucherComultaiveReq, param models.Params, header models.RequestHeader) models.Response {
	fmt.Println("[ >>>>>>>>>>>>>>>>>> Voucher OttoAg Service <<<<<<<<<<<<<<<< ]")

	var res models.Response

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[VoucherComulative-Services]",
		zap.String("NameVoucher : ", param.NamaVoucher), zap.Int("Jumlah : ", req.Jumlah),
		zap.String("CampaignID : ", req.CampaignID), zap.String("CustID : ", req.CustID),
		zap.String("CustID2 : ", req.CustID2), zap.String("ProductCode : ", param.ProductCode),
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

	logrus.Info("[ Cumulatif Reff ] : ", comulative_ref)

	var inqGagal int

	for i := 0; i < req.Jumlah; i++ {

		param.TrxID = utils.GenTransactionId()
		param.Total = i + 1

		getRespChan := make(chan models.RedeemComuResp)
		getErrChan := make(chan error)
		getRespUseVouChan := make(chan models.RedeemResponse)
		getRespUseVoucErr := make(chan error)

		go V21_RedeemVoucherOttoAg(req, param, header, getRespChan, getErrChan)
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
			go V21_UseVoucherOttoAg(req, getResp, param, getRespUseVouChan, getRespUseVoucErr, &wg)
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

		kafkaRes, err := kafka.SendPublishKafka(kafkaReq)
		if err != nil {
			fmt.Println("Failed Send Publisher Notification")
			fmt.Println("Error : ", err)
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

				fmt.Println("[VoucherComulativeService]-[GetResponseOttoag]")
				fmt.Println("[Failed to Get Data Mapping Response]")
				fmt.Println(fmt.Sprintf("[Data GetResponseOttoag : ]", getmsg))
				fmt.Println(fmt.Sprintf("[Error %v]", errmsg))
				// return res, err

				rc = getResp.Redeem.Rc
				msg = getResp.Redeem.Msg

			}

		} else {

			getmsg, errmsg := db.GetResponseOttoag("OTTOAG", getResRedeem.Rc)

			rc = getmsg.InternalRc
			msg = getmsg.InternalRd

			if errmsg != nil || getmsg.InternalRc == "" {

				fmt.Println("[VoucherComulativeService]-[GetResponseOttoag]")
				fmt.Println("[Failed to Get Data Mapping Response]")
				fmt.Println(fmt.Sprintf("[Data GetResponseOttoag : ]", getmsg))
				fmt.Println(fmt.Sprintf("[Error %v]", errmsg))
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
