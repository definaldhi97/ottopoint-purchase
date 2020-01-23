package ottoag

import (
	"encoding/json"
	"ottopoint-purchase/db"
	ottoag "ottopoint-purchase/hosts/ottoag/host"
	redismodels "ottopoint-purchase/hosts/redis_token/models"
	"ottopoint-purchase/models"
	"ottopoint-purchase/models/dbmodels"
	ottoagmodels "ottopoint-purchase/models/ottoag"
	"ottopoint-purchase/utils"

	"github.com/astaxie/beego/logs"
)

// type PaymentBillerServices struct {
// 	General models.GeneralModel
// }

func PaymentBiller(reqdata interface{}, req models.UseRedeemRequest, dataToken redismodels.TokenResp, amount int64, rrn, MemberID, namaVoucher, expDate, category string) ottoagmodels.OttoAGPaymentRes {

	res := ottoagmodels.OttoAGPaymentRes{}

	logs.Info("[PAYMENT-SERVICES][START]")

	// switch category {
	// case constants.CategoryPulsa:
	// 	res.Data = ottoagmodels.DataPayPulsa{}
	// case constants.CategoryFreeFire, constants.CategoryMobileLegend:
	// 	res.Data = ottoagmodels.DataGame{}
	// case constants.CategoryToken:
	// 	res.Data = ottoagmodels.DataPayPLNTOKEN{}
	// }

	billerHead := ottoag.PackMessageHeader(reqdata)
	logs.Info("Nama Voucher : ", namaVoucher)
	billerDataHost, err := ottoag.Send(reqdata, billerHead, "PAYMENT")

	if err = json.Unmarshal(billerDataHost, &res); err != nil {
		logs.Error("Failed to unmarshaling json response from ottoag", err)
		res = ottoagmodels.OttoAGPaymentRes{
			Rc:  "01",
			Msg: "Payment Failed",
		}

		return res
	}

	if err != nil {
		logs.Error("Failed to connect ottoag host", err)
		res = ottoagmodels.OttoAGPaymentRes{
			Rc:  "01",
			Msg: "Payment Failed",
		}

		logs.Info("[SAVE-DB-PAYMENT-Transaksi_Redeem]")

		labelPyment1 := dbmodels.TransaksiRedeem{
			AccountNumber: dataToken.Data,
			Voucher:       namaVoucher,
			CustID:        req.CustID,
			// MerchantID:    dataToken.Data.MerchantID,
			RRN:         rrn,
			ProductCode: req.ProductCode,
			Amount:      amount,
			TransType:   "Payment",
			Status:      "01 (Gagal)",
			ExpDate:     expDate,
			Institution: "Ottopay",
			ProductType: "Pulsa",
			DateTime:    utils.GetTimeFormatYYMMDDHHMMSS(),
		}
		err1 := db.Dbcon.Create(&labelPyment1).Error
		if err1 != nil {
			logs.Info("Failed Save to database", err1)
			// return err1
		}

		return res
	}

	return res
}
