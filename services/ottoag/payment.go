package ottoag

import (
	"encoding/json"
	ottoag "ottopoint-purchase/hosts/ottoag/host"
	"ottopoint-purchase/models"
	ottoagmodels "ottopoint-purchase/models/ottoag"

	"github.com/astaxie/beego/logs"
)

// type PaymentBillerServices struct {
// 	General models.GeneralModel
// }

func PaymentBiller(reqdata interface{}, reqOP interface{}, req models.UseRedeemRequest, param models.Params) ottoagmodels.OttoAGPaymentRes {

	res := ottoagmodels.OttoAGPaymentRes{}

	logs.Info("[PaymentBiller-SERVICES][START]")

	billerHead := ottoag.PackMessageHeader(reqdata)
	logs.Info("Nama Voucher : ", param.NamaVoucher)
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
		// logs.Info("[SAVE-DB-PAYMENT-Transaksi_Redeem]")

		// reqOttoag, _ := json.Marshal(&reqdata)
		// responseOttoag, _ := json.Marshal(&res)
		// reqdataOP, _ := json.Marshal(&reqOP)

		// savePay := dbmodels.TransaksiRedeem{
		// 	AccountNumber:   param.AccountNumber,
		// 	Voucher:         param.NamaVoucher,
		// 	MerchantID:      param.MerchantID,
		// 	CustID:          req.CustID,
		// 	RRN:             res.Rrn,
		// 	ProductCode:     res.Productcode,
		// 	Amount:          int64(res.Amount),
		// 	TransType:       "Payment",
		// 	ProductType:     "Pulsa",
		// 	Status:          "01 (Gagal)",
		// 	ExpDate:         param.ExpDate,
		// 	Institution:     param.InstitutionID,
		// 	CummulativeRef:  param.Reffnum,
		// 	DateTime:        utils.GetTimeFormatYYMMDDHHMMSS(),
		// 	ResponderData:   "01",
		// 	Point:           param.Point,
		// 	ResponderRc:     res.Rc,
		// 	RequestorData:   string(reqOttoag),
		// 	ResponderData2:  string(responseOttoag),
		// 	RequestorOPData: string(reqdataOP),
		// 	SupplierID:      param.SupplierID,
		// }
		// err1 := db.DbCon.Create(&savePay).Error
		// if err1 != nil {
		// 	logs.Info("Failed Save to database", err1)
		// 	// return err1
		// }

		return res
	}

	return res
}
