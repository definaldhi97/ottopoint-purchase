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
	"sync"

	"github.com/astaxie/beego/logs"
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
			logs.Info("Gagal Redeem or Inquiry")
			logs.Info("Error Message : ", getResp.Message)
			inqGagal++
			continue
		} else {
			getResp = <-getRespChan
		}

		getResp = <-getRespChan
		if getResp.Code == "00" {
			wg.Add(1)
			// req.Jumlah = i
			go UseVoucherComulative(req, getResp, param, getRespUseVouChan, getRespUseVoucErr, &wg)
			getResRedeem = <-getRespUseVouChan

		}
	}
	wg.Wait()

	logs.Info("Response OttoAG Payment : ", getResRedeem)
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

	// if inqGagal != 0 {

	// }
	// countInqFailed, _ := db.GetCountInquiryGagal(comulative_ref)

	// if countSuccess.Count == 0 && countFailed.Count == 0 && countPending.Count == 0 {
	// 	countFailed.Count = countInqFailed.Count
	// }

	/* ------ Reversal to Point ----- */
	rcUseVoucher, _ := db.GetPyenmentFailed(comulative_ref)
	logs.Info("get RC use voucher : ", rcUseVoucher)
	if rcUseVoucher.AccountNumber != "" {
		logs.Info("============= Reversal to Point ===========")
		// get Custid from user where acount nomor
		dataUser, _ := db.CheckUser(rcUseVoucher.AccountNumber)
		logs.Info("CustId user OPL from tb user : ", dataUser.CustID)
		Text := "OP009 - " + "Reversal point cause transaction " + param.NamaVoucher + " is failed"
		_, errReversal := host.TransferPoint(dataUser.CustID, strconv.Itoa(rcUseVoucher.Count), Text)
		if errReversal != nil {
			logs.Info("[ERROR DB]")
			logs.Info("[transfer point] : ", errReversal)
		}

		logs.Info("========== Send Notif ==========")
		notifReq := ottomartmodels.NotifRequest{
			AccountNumber:    rcUseVoucher.AccountNumber,
			Title:            "Reversal Point",
			Message:          fmt.Sprintf("Point anda berhasil di reversal sebesar %v", int64(rcUseVoucher.Count)),
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

	}

	/* ---- Message ----
	* Sukses  ( success == jumlah request )
	* Sukses sebagian (success != jumlah request)
	* Gagal (success == 0)
	 */
	logs.Info(" jumlah transaction payment : ", countPayment.Count)
	logs.Info(" jumlah success transaction success : ", countSuccess.Count)
	logs.Info(" jumlah success transaction Pending : ", countPending.Count)
	logs.Info(" jumlah success transaction failed : ", pyenmentFail)
	logs.Info(" jumlah request : ", req.Jumlah)

	respMessage := models.CommulativeResp{
		Success: countSuccess.Count,
		Pending: countPending.Count,
		Failed:  pyenmentFail,
	}
	// sukses 	pending 	failed
	if (respMessage.Success != 0) && (respMessage.Pending == 0) && (respMessage.Failed == 0) {
		Code_RC_Comulative = "00"
		Message_Comulative = "Sukses semua"
	}
	if (respMessage.Success != 0) && (respMessage.Pending != 0) && (respMessage.Failed == 0) {
		Code_RC_Comulative = "33"
		Message_Comulative = "Sukses sebagian"
	}
	if (respMessage.Success != 0) && (respMessage.Pending != 0) && (respMessage.Failed != 0) {
		Code_RC_Comulative = "33"
		Message_Comulative = "Sukses sebagian"
	}
	if (respMessage.Success == 0) && (respMessage.Pending != 0) && (respMessage.Failed == 0) {
		Code_RC_Comulative = "56"
		Message_Comulative = "Transaksi pending"
	}
	if (respMessage.Success == 0) && (respMessage.Pending != 0) && (respMessage.Failed != 0) {
		Code_RC_Comulative = "56"
		Message_Comulative = "Transaksi pending"
	}
	if (respMessage.Success == 0) && (respMessage.Pending == 0) && (respMessage.Failed != 0) {
		Code_RC_Comulative = "01"
		Message_Comulative = "Gagal"
	}

	// pyenmentFail := req.Jumlah - countSuccess.Count
	// pyenmentPending := req.Jumlah - countPending.Count

	/* ------ Response UseVoucher Comulative */
	logs.Info("========== Mesage from Inquiry OTTOAG and OPL ===============")
	logs.Info("Code : ", getResp.Code)
	logs.Info("Message : ", getResp.Message)
	logs.Info("=============================================================")
	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: models.CommulativeResp{
			Code:    Code_RC_Comulative,
			Msg:     Message_Comulative,
			Success: countSuccess.Count,
			Pending: countPending.Count,
			Failed:  pyenmentFail,

			//RedeemRes :
		},
	}

	return res
}
