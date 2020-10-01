package services

import (
	"encoding/json"
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	"ottopoint-purchase/hosts/opl/host"
	kafka "ottopoint-purchase/hosts/publisher/host"
	"ottopoint-purchase/models"
	"ottopoint-purchase/models/dbmodels"
	"ottopoint-purchase/utils"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/vjeantet/jodaTime"
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

		// TrxID
		param.TrxID = utils.GenTransactionId()

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
	fmt.Println("get RC use voucher : ", rcUseVoucher)
	if rcUseVoucher.AccountNumber != "" {
		fmt.Println("============= Reversal to Point ===========")
		// get Custid from user where acount nomor

		Text := param.TrxID + param.InstitutionID + constants.CodeReversal + "OP009 - Reversal point cause transaction " + param.NamaVoucher + " is failed"
		// Text := "OP009 - " + "Reversal point cause transaction " + param.NamaVoucher + " is failed"
		sendReversal, errReversal := host.TransferPoint(param.AccountId, strconv.Itoa(rcUseVoucher.Count), Text)
		if errReversal != nil {
			fmt.Println("[ERROR DB]")
			fmt.Println("[transfer point] : ", errReversal)
		}

		expired := ExpiredPointService()

		saveReversal := dbmodels.TEarning{
			ID: utils.GenerateTokenUUID(),
			// EarningRule     :,
			// EarningRuleAdd  :,
			PartnerId: param.InstitutionID,
			// ReferenceId     : ,
			TransactionId: param.TrxID,
			// ProductCode     :,
			// ProductName     :,
			AccountNumber: param.AccountNumber,
			// Amount          :,
			Point: int64(rcUseVoucher.Count),
			// Remark          :,
			Status:           constants.Success,
			StatusMessage:    "Success",
			PointsTransferId: sendReversal.PointsTransferId,
			// RequestorData   :,
			// ResponderData   :,
			TransType:       constants.CodeReversal,
			AccountId:       param.AccountId,
			ExpiredPoint:    expired,
			TransactionTime: time.Now(),
		}

		errSaveReversal := db.DbCon.Create(&saveReversal).Error
		if errSaveReversal != nil {

			fmt.Println(fmt.Sprintf("[Failed Save Reversal to DB]-[Error : %v]", errSaveReversal))
			fmt.Println("[PackageServices]-[SaveEarning]")

			fmt.Println(">>> Save CSV <<<")
			name := jodaTime.Format("dd-MM-YYYY", time.Now()) + ".csv"
			go utils.CreateCSVFile(saveReversal, name)

		}

		fmt.Println("========== Send Publisher ==========")

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
			fmt.Println("Gagal Send Publisher")
			fmt.Println("Error : ", err)
		}

		fmt.Println("Response Publisher : ", kafkaRes)

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
		m = getMsgCummulative(rc, msg)

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

func getMsgCummulative(rc, msg string) string {

	var codeMsg string

	getmsg, errmsg := db.GetResponseCummulativeOttoAG(rc)
	if errmsg != nil || getmsg.InternalRc == "" {

		fmt.Println("[VoucherComulativeService]-[GetResponseCummulativeOttoAG]")
		fmt.Println("[Failed to Get Data Mapping Response]")
		fmt.Println(fmt.Sprintf("[Data GetResponseOttoag : ]", getmsg))
		fmt.Println(fmt.Sprintf("[Error %v]", errmsg))
		// return res, err

		codeMsg = msg

		return codeMsg
	}

	// codeRc = getmsg.InternalRc
	// codeMsg = strings.Replace(getmsg.InternalRd, "[x]", "%v", 10)
	codeMsg = getmsg.InternalRd

	return codeMsg
}

func ExpiredPointService() string {

	fmt.Println(">>> ExpiredPointService <<<")

	get, err := host.SettingsOPL()
	if err != nil || get.Settings.ProgramName == "" {

		fmt.Println(fmt.Sprintf("[Error : %v]", err))
		fmt.Println("[PackageBulkService]-[ExpiredPointService]")

	}

	data := get.Settings.PointsDaysActiveCount + 1

	expired := utils.FormatTimeString(time.Now(), 0, 0, data)

	return expired

}
