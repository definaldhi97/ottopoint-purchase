package v2_migrate

import (
	"fmt"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"
	"sync"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type UseVoucherOttoAgService struct {
	General models.GeneralModel
}

func (t UseVoucherOttoAgService) UseVoucherOttoAg(req models.VoucherComultaiveReq, param models.Params) {
	fmt.Println("[ >>>>>>>>>>>>>>>>>> Use Voucher OttoAg Service <<<<<<<<<<<<<<<< ]")

	// var res models.Response

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[VoucherComulative-Services]",
		zap.String("NameVoucher : ", param.NamaVoucher), zap.Int("Jumlah : ", req.Jumlah),
		zap.String("CampaignID : ", req.CampaignID), zap.String("CustID : ", req.CustID),
		zap.String("CustID2 : ", req.CustID2), zap.String("ProductCode : ", param.ProductCode),
		zap.String("AccountNumber : ", param.AccountNumber), zap.String("InstitutionID : ", param.InstitutionID))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[RedeemVoucher]")
	defer span.Finish()

	// Message_Comulative := ""
	// Code_RC_Comulative := ""

	wg := sync.WaitGroup{}

	getResp := models.RedeemComuResp{}
	// getResRedeem := models.RedeemResponse{}

	/*---- generate comulative_ref ----*/
	comulative_ref := utils.GenTransactionId()
	param.Reffnum = comulative_ref

	fmt.Println("Cumulatif reff : ", comulative_ref)

	var inqGagal int

	for i := 0; i < req.Jumlah; i++ {

		param.TrxID = utils.GenTransactionId()
		param.Total = i + 1

		getRespChan := make(chan models.RedeemComuResp)
		getErrChan := make(chan error)
		// getRespUseVouChan := make(chan models.RedeemResponse)
		// getRespUseVoucErr := make(chan error)

		go RedeemVoucherOttoAg(req, param, getRespChan, getErrChan)
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

	}
	wg.Wait()
}
