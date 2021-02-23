package ottoag

import (
	models "ottopoint-purchase/models/ottoag"

	"github.com/astaxie/beego/logs"
	"github.com/sirupsen/logrus"
)

// Validasi Data Inquiry
func ValidateDataInq(req models.OttoAGInquiryRequest) bool {
	logrus.Info("[ValidateDataInq]")
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
