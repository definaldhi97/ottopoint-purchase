package services

import (
	"fmt"
	"ottopoint-purchase/db"
	"ottopoint-purchase/utils"

	"github.com/sirupsen/logrus"
	"github.com/vjeantet/jodaTime"
)

type PatchingInvoice struct {
	Success int `json:"success"`
	Failed  int `json:"failed"`
	Total   int `json:"total"`
}

func PatchingInvoiceNumberService() PatchingInvoice {
	res := PatchingInvoice{}

	var success, failed int
	getData, errData := db.GeDataPatching()
	res.Total = len(getData)
	if errData != nil {
		logrus.Error("[PackageServices]-[PatchingInvoiceNumberService]")
		logrus.Error(fmt.Sprintf("[GeDataPatching]-[Error : %v]", errData))

		return res

	}

	for i := 0; i < len(getData); i++ {

		inv := "INV" + jodaTime.Format("YYYYMMdd", getData[i].CreatedAT) + utils.GenTransactionId()[7:11]

		update := db.UpdateDataPatching(inv, getData[i].TransactionId)
		if update != nil {

			logrus.Error("[PackageServices]-[PatchingInvoiceNumberService]")
			logrus.Error(fmt.Sprintf("[GeDataPatching]-[Error : %v]-[TrxId : %v]", errData, getData[i].TransactionId))

			failed++
			continue

		}

		success++

	}

	res.Success = success
	res.Failed = failed
	return res

}
