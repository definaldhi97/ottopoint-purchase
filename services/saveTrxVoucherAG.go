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

	fmt.Println(fmt.Sprintf("[Start]-[SaveTransactionVoucherAgMigrate]-[%v]", transType))

	var ExpireDate time.Time
	var redeemDate time.Time

	var saveStatus string
	isUsed := false
	switch status {
	case "00":
		saveStatus = constants.Success
	case "09":
		saveStatus = constants.Pending
	case "01":
		saveStatus = constants.Failed
		isUsed = true
	}

	reqUV, _ := json.Marshal(&reqdata)   // Req UV
	responseUV, _ := json.Marshal(&res)  // Response UV
	reqdataOP, _ := json.Marshal(&reqOP) // Req Service

	// expDate := ""
	// if param.ExpDate != "" {
	// 	layout := "2006-01-02 15:04:05"
	// 	parse, _ := time.Parse(layout, param.ExpDate)

	// 	expDate = jodaTime.Format("YYYY-MM-dd", parse)
	// }

	if transType == constants.CODE_TRANSTYPE_REDEMPTION {
		ExpireDate = utils.ExpireDateVoucherAGt(timeExpVouc)
		redeemDate = time.Now()
	}

	save := dbmodels.TSpending{
		ID:              utils.GenerateTokenUUID(),
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
		RequestorData:   string(reqUV),
		ResponderData:   string(responseUV),
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
		MRewardID:         param.RewardID,
		ProductCategoryID: param.CategoryID,
		MProductID:        param.ProductID,
		PointsTransferID:  param.PointTransferID,
	}

	err := db.DbCon.Create(&save).Error
	if err != nil {
		logrus.Info(fmt.Sprintf("[Error : %v]", err))
		logrus.Info("[Failed Save to DB]")

		name := jodaTime.Format("dd-MM-YYYY", time.Now()) + ".csv"
		go utils.CreateCSVFile(save, name)

		// return err

	}

}

func SaveDBVoucherAgMigrate(id, institution, coupon, vouchercode, phone, custIdOPL, campaignID string) {

	fmt.Println("[SaveDB]-[SaveDBVoucherAgMigrate]")

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
		fmt.Println("[Failed Save to DB ]", err)
		fmt.Println("[Package-Services]-[SaveDBVoucherAgMigrate]")
	}
}
