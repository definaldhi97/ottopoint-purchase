package services

import (
	"ottopoint-purchase/db"
	"ottopoint-purchase/hosts/opl/host"
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

	for i := 0; i < req.Jumlah; i++ {
		// wg.Add(2)

		getRespChan := make(chan models.RedeemComuResp)
		getErrChan := make(chan error)
		getRespUseVouChan := make(chan models.RedeemResponse)
		getRespUseVoucErr := make(chan error)

		go RedeemComulativeVoucher(req, param, getRespChan, getErrChan)

		if getErr := <-getErrChan; getErr != nil {
			logs.Info("=========== Error redeem voucher ===========")
			// return
			continue
		} else {
			logs.Info("=========== Continue to use voucher ===========")
			getResp = <-getRespChan
		}

		// wg.Add(1)
		logs.Info("========== response code redeem : ", getResp.Code)
		if getResp.Code == "00" {
			wg.Add(1)
			logs.Info("================ doing use voucher ================")
			go UseVoucherComulative(req, getResp, param, getRespUseVouChan, getRespUseVoucErr, &wg)
			getResRedeem = <-getRespUseVouChan

		}
	}
	wg.Wait()

	/* Get total count pyenment success dan failed dari DB */
	logs.Info("===== comulative_ref : ", comulative_ref)
	countSuccess, _ := db.GetCountSucc_Pyenment(comulative_ref)
	countPending, _ := db.GetCountPending_Pyenment(comulative_ref)
	pyenmentFail := req.Jumlah - countSuccess.Count

	/* ------ Reversal to Point ----- */
	rcUseVoucher, _ := db.GetPyenmentFailed(comulative_ref)
	logs.Info("get RC use voucher : ", rcUseVoucher)
	if rcUseVoucher.AccountNumber != "" {
		logs.Info("============= Reversal to Point ===========")
		logs.Info("Response Code from Payment OttoAG : ", getResRedeem.Rc)
		// get Custid from user where acount nomor
		dataUser, _ := db.CheckUser(rcUseVoucher.AccountNumber)
		logs.Info("CustId user OPL from tb user : ", dataUser.CustID)
		Text := "Reversal point cause transaction " + param.NamaVoucher + " is failed"
		_, errReversal := host.TransferPoint(dataUser.CustID, strconv.Itoa(rcUseVoucher.Count), Text)
		if errReversal != nil {
			logs.Info("[ERROR DB]")
			logs.Info("[transfer point] : ", errReversal)
		}

	}

	/* ---- Message ----
	* Sukses  ( success == jumlah request )
	* Sukses sebagian (success != jumlah request)
	* Gagal (success == 0)
	 */
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
