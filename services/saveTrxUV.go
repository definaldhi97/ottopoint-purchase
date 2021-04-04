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

func SaveTransactionUV(param models.Params, res interface{}, reqdata interface{}, reqOP interface{}, trasnType, status string, expVoucher int) {

	fmt.Println(fmt.Sprintf("[Start-SaveDB]-[UltraVoucher]-[%v]", trasnType))

	fmt.Println("Param.Mreward : ", param.RewardID)

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

	isUsed := false
	if status == "01" {
		isUsed = true
	}

	if trasnType == constants.CODE_TRANSTYPE_REDEMPTION {
		ExpireDate = utils.ExpireDateVoucherAGt(expVoucher)
		redeemDate = time.Now()
	}

	reqUV, _ := json.Marshal(&reqdata)   // Req UV
	responseUV, _ := json.Marshal(&res)  // Response UV
	reqdataOP, _ := json.Marshal(&reqOP) // Req Service

	// if param.PaymentMethod == constants.SplitBillMethod {

	// 	_, errUpdate := db.UpdateTransactionSplitBill(isUsed, param.TrxID, saveStatus, param.DataSupplier.Rc, param.DataSupplier.Rd, responseOttoag, reqOttoag, reqdataOP)
	// 	if errUpdate != nil {

	// 		logrus.Error(fmt.Sprintf("[UpdateTransactionSplitBill]-[SaveTransactionOttoAg]"))
	// 		logrus.Error(fmt.Sprintf("[TrxID : %v]-[Error : %v]", trxID, errUpdate))

	// 		return
	// 	}

	// 	return

	// }

	save := dbmodels.TSpending{
		ID:            utils.GenerateTokenUUID(),
		AccountNumber: param.AccountNumber,
		Voucher:       param.NamaVoucher,
		MerchantID:    param.MerchantID,
		// CustID:          param.CustID,
		RRN:               param.RRN,
		TransactionId:     param.TrxID,
		ProductCode:       param.ProductCode,
		Amount:            int64(param.Point),
		TransType:         trasnType,
		IsUsed:            isUsed,
		ProductType:       param.ProductType,
		Status:            saveStatus,
		ExpDate:           utils.DefaultNulTime(ExpireDate),
		Institution:       param.InstitutionID,
		CummulativeRef:    param.CumReffnum,
		DateTime:          utils.GetTimeFormatYYMMDDHHMMSS(),
		Point:             param.Point,
		ResponderRc:       param.DataSupplier.Rc,
		ResponderRd:       param.DataSupplier.Rd,
		RequestorData:     string(reqUV),
		ResponderData:     string(responseUV),
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
		CreatedAT:         param.TrxTime,
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

func SaveDB(id, institution, coupon, vouchercode, phone, custIdOPL, campaignID string) {
	fmt.Println("[SaveDB]-[UltraVoucherServices]")
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
		fmt.Println("[Package-Services]-[UltraVoucherServices]")
		// return err
	}
}
