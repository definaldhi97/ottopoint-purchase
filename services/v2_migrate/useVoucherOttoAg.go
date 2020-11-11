package v2_migrate

import (
	"fmt"
	"ottopoint-purchase/models"
	"sync"
)

func UseVoucherOttoAg(req models.VoucherComultaiveReq, redeemComu models.RedeemComuResp, param models.Params, getRespChan chan models.RedeemResponse, ErrRespUseVouc chan error, wg *sync.WaitGroup) {
	fmt.Println("[ >>>>>>>>>>>>>>>>>> Use Voucher Comulative Voucher Otto AG <<<<<<<<<<<<<<<< ]")

	defer wg.Done()
	defer close(getRespChan)
	defer close(ErrRespUseVouc)

	// get CustID
	// dataUser, errUser := db.CheckUser(param.AccountNumber)

}
