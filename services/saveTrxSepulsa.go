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

func SaveTransactionSepulsa(param models.Params, res interface{}, reqdata interface{}, reqOP models.VoucherComultaiveReq, transType, status string) {

	nameservice := fmt.Sprintf("[PackageServices]-[Start]-[SaveTransactionSepulsa][%v]", transType)
	logReq := fmt.Sprintf("[RRN : %v || MReqard]", param.RRN, param.RewardID)

	logrus.Info(nameservice)
	logrus.Info(logReq)

	var ExpireDate time.Time
	var redeemDate time.Time

	var saveStatus string
	switch status {
	case "00":
		saveStatus = constants.Success
	case "09":
		saveStatus = constants.Pending
	case "01":
		saveStatus = constants.Failed
	}

	if transType == constants.CODE_TRANSTYPE_REDEMPTION {
		ExpireDate = utils.ExpireDateVoucherAGt(constants.EXPDATE_VOUCHER)
		redeemDate = time.Now()
	}

	reqSepulsa, _ := json.Marshal(&reqdata)
	responseSepulsa, _ := json.Marshal(&res)
	reqdataOP, _ := json.Marshal(&reqOP)

	// timeRedeem := jodaTime.Format("dd-MM-YYYY HH:mm:ss", time.Now())

	idSpending := utils.GenerateTokenUUID()

	save := dbmodels.TSpending{
		ID:                idSpending,
		AccountNumber:     param.AccountNumber,
		Voucher:           param.NamaVoucher,
		MerchantID:        param.MerchantID,
		CustID:            reqOP.CustID,
		RRN:               param.RRN,
		TransactionId:     param.TrxID,
		ProductCode:       param.ProductCode,
		Amount:            int64(param.Amount),
		TransType:         transType,
		IsUsed:            true,
		UsedAt:            &redeemDate,
		ProductType:       param.ProductType,
		Status:            saveStatus,
		ExpDate:           utils.DefaultNulTime(ExpireDate),
		Institution:       param.InstitutionID,
		CummulativeRef:    param.CumReffnum,
		DateTime:          utils.GetTimeFormatYYMMDDHHMMSS(),
		Point:             param.Point,
		ResponderRc:       param.DataSupplier.Rc,
		ResponderRd:       param.DataSupplier.Rd,
		RequestorData:     string(reqSepulsa),
		ResponderData:     string(responseSepulsa),
		RequestorOPData:   string(reqdataOP),
		SupplierID:        param.SupplierID,
		CouponId:          param.CouponID,
		CampaignId:        param.CampaignID,
		AccountId:         param.AccountId,
		RedeemAt:          utils.DefaultNulTime(redeemDate),
		Comment:           param.Comment,
		MRewardID:         &param.RewardID,
		ProductCategoryID: param.CategoryID,
		MProductID:        &param.ProductID,
		PointsTransferID:  param.PointTransferID,
		IsCallback:        false,
		PaymentMethod: 2,
		InvoiceNumber: param.InvoiceNumber,
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

func SaveDBSepulsa(id, institution, coupon, vouchercode, phone, custIdOPL, campaignID string) {
	fmt.Println("[SaveDB]-[SepulsaVoucherService]")
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
		fmt.Println("[Failed Save to DB]", err)
		fmt.Println("[Package-Service]-[SepulsaService]")
	}
}
