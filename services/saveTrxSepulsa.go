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

	fmt.Println(fmt.Sprintf("[Start-SaveDB]-[Sepulsa]-[%v]", transType))

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

	save := dbmodels.TSpending{
		ID:                utils.GenerateTokenUUID(),
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
