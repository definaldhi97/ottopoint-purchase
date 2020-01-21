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

// type InquiryBillerServices struct {
// 	General models.GeneralModel
// }

func InquiryBiller(reqdata interface{}, req models.UseRedeemRequest, dataToken redismodels.TokenResp, MemberID, namaVoucher, expDate string) ottoagmodels.OttoAGInquiryResponse {
	response := ottoagmodels.OttoAGInquiryResponse{}

	logs.Info("[INQUIRY-SERVICES][START]")

	// sugarLogger := t.General.OttoZaplog
	// sugarLogger.Info("[ottoag-Services]",
	// 	zap.String("reqdata", reqdata.AccountNumber))
	// span, _ := opentracing.StartSpanFromContext(t.General.Context, "[ottoag-Services]")
	// defer span.Finish()

	headOttoAg := ottoag.PackMessageHeader(reqdata)
	billerDataHost, err := ottoag.Send(reqdata, headOttoAg, "INQUIRY")
	if err = json.Unmarshal(billerDataHost, &response); err != nil {
		logs.Error("Failed to unmarshaling json response from ottoag", err)
		response = ottoagmodels.OttoAGInquiryResponse{
			Rc:  "01",
			Msg: "Inquiry Failed",
		}

		return response
	}

	if err != nil {
		logs.Error("Failed to connect ottoag host", err)
		response = ottoagmodels.OttoAGInquiryResponse{
			Rc:  "01",
			Msg: "Inquiry Failed",
		}

		logs.Info("[SAVE-DB-INQUIRY-Transaksi_Redeem]")

		saveInq := dbmodels.TransaksiRedeem{
			AccountNumber: dataToken.Data.AccountNumber,
			Voucher:       namaVoucher,
			CustID:        req.CustID,
			MerchantID:    dataToken.Data.MerchantID,
			RRN:           response.Rrn,
			ProductCode:   req.ProductCode,
			Amount:        response.Amount,
			TransType:     "Payment",
			Status:        "01 (Gagal)",
			ExpDate:       expDate,
			Institution:   "Ottopay",
			ProductType:   "Pulsa",
			DateTime:      utils.GetTimeFormatYYMMDDHHMMSS(),
		}
		err1 := db.Dbcon.Create(&saveInq).Error
		if err1 != nil {
			logs.Info("Failed Save to database", err1)
			// return err1
		}

		return response
	}

	return response
}
