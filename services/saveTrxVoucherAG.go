package services

import (
	"encoding/json"
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	"ottopoint-purchase/models"
	"ottopoint-purchase/models/dbmodels"
	"ottopoint-purchase/utils"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vjeantet/jodaTime"
)

func SaveTransactionVoucherAgMigrate(param models.Params, res interface{}, reqdata interface{}, reqOP interface{}, transType, status string, timeExpVouc int) {

	nameservice := fmt.Sprintf("[PackageServices]-[Start]-[SaveTransactionVoucherAgMigrate][%v]", transType)
	logReq := fmt.Sprintf("[RRN : %v]", param.RRN)

	logrus.Info(nameservice)
	logrus.Info(logReq)

	var ExpireDate time.Time
	var redeemDate time.Time

	var saveStatus string
	isUsed := true
	switch status {
	case "00":
		saveStatus = constants.Success
	case "09":
		saveStatus = constants.Pending
	case "01":
		saveStatus = constants.Failed
		isUsed = true
	}

	if param.SupplierID == constants.CODE_VENDOR_UV {
		isUsed = false
	}

	reqVAG, _ := json.Marshal(&reqdata)  // Req UV
	responseVAG, _ := json.Marshal(&res) // Response UV
	reqdataOP, _ := json.Marshal(&reqOP) // Req Service

	if transType == constants.CODE_TRANSTYPE_REDEMPTION {
		ExpireDate = utils.ExpireDateVoucherAGt(timeExpVouc)
		redeemDate = time.Now()
	}

	idSpending := utils.GenerateTokenUUID()

	save := dbmodels.TSpending{
		ID:              idSpending,
		AccountNumber:   param.AccountNumber,
		RRN:             param.RRN,
		Voucher:         param.NamaVoucher,
		MerchantID:      param.MerchantID,
		CustID:          param.CustID,
		TransactionId:   param.TrxID,
		ProductCode:     param.ProductCodeInternal,
		Amount:          int64(param.Amount),
		TransType:       transType,
		IsUsed:          isUsed,
		ProductType:     param.ProductType,
		Status:          saveStatus,
		Institution:     param.InstitutionID,
		CummulativeRef:  param.CumReffnum,
		DateTime:        utils.GetTimeFormatYYMMDDHHMMSS(),
		Point:           param.Point,
		ResponderRc:     param.DataSupplier.Rc,
		ResponderRd:     param.DataSupplier.Rd,
		RequestorData:   string(reqVAG),
		ResponderData:   string(responseVAG),
		RequestorOPData: string(reqdataOP),
		SupplierID:      param.SupplierID,
		CouponId:        param.CouponID,
		CampaignId:      param.CampaignID,
		AccountId:       param.AccountId,
		VoucherCode:     param.CouponCode,
		VoucherLink:     param.VoucherLink,
		ExpDate:         utils.DefaultNulTime(ExpireDate),
		RedeemAt:        utils.DefaultNulTime(redeemDate),

		Comment:           param.Comment,
		MRewardID:         &param.RewardID,
		ProductCategoryID: param.CategoryID,
		MProductID:        &param.ProductID,
		PointsTransferID:  param.PointTransferID,
		IsCallback:        false,
		InvoiceNumber: param.InvoiceNumber,
		PaymentMethod: 2,
	}

	err := db.DbCon.Create(&save).Error
	if err != nil {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[TSpending]-[Error : %v]", err))
		logrus.Println(logReq)

		name := jodaTime.Format("YYYY-MM-dd", time.Now()) + ".csv"
		go utils.CreateCSVFile(save, name)

		// return err

	}

	savePayment := dbmodels.TPayment{
		ID:             utils.GenerateUUID(),
		TSpendingID:    idSpending,
		ExternalReffId: param.RRN,
		TransType:      transType,
		Value:          int64(param.Point),
		ValueType:      constants.TypePoint,
		Status:         status,
		// ResponderRc   : ,
		// ResponderRd   : ,
		CreatedBy: constants.CreatedbySystem,
		// UpdatedBy     : ,
		CreatedAt: time.Now(),
		// UpdatedAt     : ,
	}

	errPayment := db.DbCon.Create(&savePayment).Error
	if errPayment != nil {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[SavePayment]-[Error : %v]", errPayment))
		logrus.Println(logReq)

		name := jodaTime.Format("YYYY-MM-dd", time.Now()) + ".csv"
		go utils.CreateCSVFile(save, name)

		// return

	}

}

func SaveDBVoucherAgMigrate(id, institution, coupon, vouchercode, phone, custIdOPL, campaignID string) {

	nameservice := "[PackagePayment]-[Start]-[SaveDBVoucherAgMigrate]"
	logReq := fmt.Sprintf("[Coupon : %v || VoucherCode : %v]", coupon, vouchercode)

	save := dbmodels.UserMyVocuher{
		ID:            id,
		InstitutionID: institution,
		CouponID:      coupon,
		VoucherCode:   vouchercode,
		Phone:         phone,
		AccountId:     custIdOPL,
		CampaignID:    campaignID,
	}

	err := db.DbCon.Create(&save).Error
	if err != nil {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[UserMyVocuher]-[Error : %v]", err))
		logrus.Println(logReq)

		name := jodaTime.Format("YYYY-MM-dd", time.Now()) + ".csv"
		go utils.CreateCSVFile(save, name)

	}
}
