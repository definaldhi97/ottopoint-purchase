package ottoag

import (
	models "ottopoint-purchase/models/ottoag"

	"github.com/astaxie/beego/logs"
)

// Validasi Data Inquiry
func ValidateDataInq(req models.OttoAGInquiryRequest) bool {
	logs.Info("[ValidateDataInq]")
	if req.IssuerID == "" {
		logs.Error("IssuerID is empty")
		return false
	}

	if req.AccountNumber == "" {
		logs.Error("Account Number is empty")
		return false
	}
	return true
}
