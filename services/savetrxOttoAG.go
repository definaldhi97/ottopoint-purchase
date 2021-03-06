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

func SaveTransactionOttoAg(param models.Params, res interface{}, reqdata interface{}, reqOP interface{}, status string) string {

	nameservice := fmt.Sprintf("[PackageServices]-[Start]-[SaveTransactionOttoAg][%v]", param.TransType)
	logReq := fmt.Sprintf("[RRN : %v || MReward]", param.RRN, param.RewardID)

	logrus.Info(nameservice)
	logrus.Info(logReq)

	// validasi vidio is_used -> false
	isUsed := true
	// codeVoucher := param.VoucherCode
	var codeVoucher, invNum string
	var paymethod int
	var ExpireDate time.Time
	var redeemDate time.Time
	var usedAt time.Time
	trxID := utils.GenTransactionId()

	if param.TransType == constants.CODE_TRANSTYPE_REDEMPTION {
		// timeRedeem := jodaTime.Format("dd-MM-YYYY HH:mm:ss", time.Now())
		// redeemDate = timeRedeem
		codeVoucher = utils.EncryptVoucherCode(param.VoucherCode, param.CouponID)
		isUsed = true
		ExpireDate = utils.ExpireDateVoucherAGt(constants.EXPDATE_VOUCHER)
		redeemDate = time.Now()
		trxID = param.TrxID
		usedAt = time.Now()

		invNum = param.InvoiceNumber
		paymethod = 2

		if param.Category == constants.CategoryVidio {
			isUsed = false // isUsed status untuk used
			usedAt = time.Time{}
		}

	}

	var saveStatus string
	switch status {
	case "00":
		saveStatus = constants.Success
	case "09":
		saveStatus = constants.Pending
	case "01":
		saveStatus = constants.Failed
		isUsed = true
	}

	reqOttoag, _ := json.Marshal(&reqdata)  // Req Ottoag
	responseOttoag, _ := json.Marshal(&res) // Response Ottoag
	reqdataOP, _ := json.Marshal(&reqOP)    // Req Service

	idSpending := utils.GenerateTokenUUID()

	save := dbmodels.TSpending{
		ID:            idSpending,
		AccountNumber: param.AccountNumber,
		Voucher:       param.NamaVoucher,
		MerchantID:    param.MerchantID,
		CustID:        param.CustID,
		RRN:           param.RRN,
		// TransactionId: param.TrxID,
		TransactionId: trxID,
		ProductCode:   param.ProductCode,
		Amount:        int64(param.Amount),
		TransType:     param.TransType,
		// IsUsed:          true,
		IsUsed:      isUsed,
		ProductType: param.ProductType,
		Status:      saveStatus,
		// ExpDate:         param.ExpDate,
		ExpDate:           utils.DefaultNulTime(ExpireDate),
		Institution:       param.InstitutionID,
		CummulativeRef:    param.CumReffnum,
		DateTime:          utils.GetTimeFormatYYMMDDHHMMSS(),
		Point:             param.Point,
		ResponderRc:       param.DataSupplier.Rc,
		ResponderRd:       param.DataSupplier.Rd,
		RequestorData:     string(reqOttoag),
		ResponderData:     string(responseOttoag),
		RequestorOPData:   string(reqdataOP),
		SupplierID:        param.SupplierID,
		RedeemAt:          utils.DefaultNulTime(redeemDate),
		CampaignId:        param.CampaignID,
		VoucherCode:       codeVoucher,
		CouponId:          param.CouponID,
		AccountId:         param.AccountId,
		ProductCategoryID: param.CategoryID,
		Comment:           param.Comment,
		MRewardID:         &param.RewardID,
		MProductID:        &param.ProductID,
		PointsTransferID:  param.PointTransferID,
		UsedAt:            utils.DefaultNulTime(usedAt),
		IsCallback:        false,
		PaymentMethod: paymethod,
		InvoiceNumber: invNum,
	}

	err := db.DbCon.Create(&save).Error
	if err != nil {

		fmt.Println(fmt.Sprintf("[Error : %v]", err))
		fmt.Println("[Failed saveTransactionOttoAg to DB]")
		fmt.Println(fmt.Sprintf("[TransType : %v || RRN : %v]", param.TransType, param.RRN))

		name := jodaTime.Format("dd-MM-YYYY", time.Now()) + ".csv"
		go utils.CreateCSVFile(save, name)

		return "Gagal Save"

	}

	if param.TransType == constants.CODE_TRANSTYPE_REDEMPTION {
		savePayment := dbmodels.TPayment{
			ID:             utils.GenerateUUID(),
			TSpendingID:    idSpending,
			ExternalReffId: param.RRN,
			TransType:      param.TransType,
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

			// return "Gagal Save"

		}
	}

	return "Berhasil Save"
}
