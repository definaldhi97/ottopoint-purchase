package ottoag

import (
	"encoding/json"
	ottoag "ottopoint-purchase/hosts/ottoag/host"
	"ottopoint-purchase/models"
	ottoagmodels "ottopoint-purchase/models/ottoag"

	"github.com/astaxie/beego/logs"
)

// type InquiryBillerServices struct {
// 	General models.GeneralModel
// }

func InquiryBiller(reqdata interface{}, reqOP interface{}, req models.UseRedeemRequest, param models.Params) (ottoagmodels.OttoAGInquiryResponse, error) {
	resOttAG := ottoagmodels.OttoAGInquiryResponse{}

	logs.Info("[InquiryBiller-SERVICES][START]")

	// sugarLogger := t.General.OttoZaplog
	// sugarLogger.Info("[ottoag-Services]",
	// 	zap.String("reqdata", reqdata.AccountNumber))
	// span, _ := opentracing.StartSpanFromContext(t.General.Context, "[ottoag-Services]")
	// defer span.Finish()

	logs.Info("[InquiryBiller-SERVICES][REQUEST :]", reqdata)
	headOttoAg := ottoag.PackMessageHeader(reqdata)
	billerDataHost, err := ottoag.Send(reqdata, headOttoAg, "INQUIRY")
	if err = json.Unmarshal(billerDataHost, &resOttAG); err != nil {
		logs.Info("[INQUIRY-SERVICES-01]")
		logs.Error("Failed to unmarshaling json response from ottoag", err)
		resOttAG = ottoagmodels.OttoAGInquiryResponse{
			Rc:  "01",
			Msg: "Inquiry Failed",
		}

		return resOttAG, err
	}

	if err != nil {
		logs.Info("[INQUIRY-SERVICES-02]")
		logs.Error("Failed to connect ottoag host", err)
		resOttAG = ottoagmodels.OttoAGInquiryResponse{
			Rc:  "01",
			Msg: "Inquiry Failed",
		}

		// logs.Info("[SAVE-DB-INQUIRY-Transaksi_Redeem]")

		// reqOttoag, _ := json.Marshal(&reqdata)
		// responseOttoag, _ := json.Marshal(&resOttAG)
		// reqdataOP, _ := json.Marshal(&reqOP)

		// saveInq := dbmodels.TSpending{
		// 	AccountNumber:   param.AccountNumber,
		// 	Voucher:         param.NamaVoucher,
		// 	MerchantID:      param.MerchantID,
		// 	CustID:          req.CustID,
		// 	RRN:             resOttAG.Rrn,
		// 	ProductCode:     req.ProductCode,
		// 	Amount:          resOttAG.Amount,
		// 	TransType:       "Inquiry",
		// 	ProductType:     "Pulsa",
		// 	Status:          "01 (Gagal)",
		// 	ExpDate:         param.ExpDate,
		// 	Institution:     param.InstitutionID,
		// 	CummulativeRef:  param.Reffnum,
		// 	DateTime:        utils.GetTimeFormatYYMMDDHHMMSS(),
		// 	ResponderData:   "01",
		// 	Point:           param.Point,
		// 	ResponderRc:     resOttAG.Rc,
		// 	RequestorData:   string(reqOttoag),
		// 	ResponderData2:  string(responseOttoag),
		// 	RequestorOPData: string(reqdataOP),
		// 	SupplierID:      param.SupplierID,
		// }
		// err1 := db.DbCon.Create(&saveInq).Error
		// if err1 != nil {
		// 	logs.Info("Failed Save to database", err1)
		// 	// return err1
		// }

		return resOttAG, err
	}

	return resOttAG, nil
}
