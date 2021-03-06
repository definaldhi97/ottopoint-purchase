package redeemtion

import (
	"encoding/json"
	"errors"
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	vgmodels "ottopoint-purchase/hosts/voucher_aggregator/models"
	"ottopoint-purchase/models"
	"ottopoint-purchase/models/dbmodels"
	"ottopoint-purchase/utils"
	"time"

	"ottopoint-purchase/services/v2.1/Trx"

	"github.com/sirupsen/logrus"
	"github.com/vjeantet/jodaTime"
)

func RedeemtionDummyService(req models.VoucherComultaiveReq, param models.Params, header models.RequestHeader) models.Response {
	res := models.Response{
		Meta: utils.ResponseMetaOK(),
	}

	param.TrxID = utils.GenTransactionId()
	param.CumReffnum = utils.GenTransactionId()
	textCommentSpending := param.TrxID + "#" + param.NamaVoucher
	param.Comment = textCommentSpending

	nameservice := "[PackageRedeemtion_V21_Services]-[RedeemtionDummyService]"
	logReq := fmt.Sprintf("[CampaignID : %v, AccountNumber : %v]", req.CampaignID, param.AccountNumber)

	redeem, errRedeem := Trx.V21_Redeem_PointandVoucher(1, param, header)
	param.PointTransferID = redeem.PointTransferID

	var coupon string
	for _, val := range redeem.CouponseVouch {
		coupon = val.CouponsID
	}

	param.CouponID = coupon

	if errRedeem != nil || redeem.Rc != "00" {
		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[V21_Redeem_PointandVoucher]-[Error : %v]", errRedeem))
		logrus.Println(logReq)

		save := saveTrxRedeemtionDUmmy(param, req, constants.Failed)
		logrus.Info("[Response Save : %v]", save)

		res = utils.GetMessageResponse(res, 500, false, errors.New("Internal Server Error"))

		return res
	}

	save := saveTrxRedeemtionDUmmy(param, req, constants.Success)
	logrus.Info("[Response Save : %v]", save)

	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: vgmodels.ResponseVoucherAg{
			Code:    "00",
			Msg:     "Success",
			Success: req.Jumlah,
			Failed:  0,
			Pending: 0,
		},
	}

	return res
}

func saveTrxRedeemtionDUmmy(param models.Params, req interface{}, status string) string {

	logrus.Info("[Start-SaveTrxRedeemtionDUmmy]")

	var codeVoucher string
	var ExpireDate time.Time

	var saveStatus string
	switch status {
	case "00":
		saveStatus = constants.Success
	case "09":
		saveStatus = constants.Pending
	case "01":
		saveStatus = constants.Failed
	}

	reqdataOP, _ := json.Marshal(&req) // Req Service

	save := dbmodels.TSpending{
		ID:             utils.GenerateTokenUUID(),
		AccountNumber:  param.AccountNumber,
		Voucher:        param.NamaVoucher,
		MerchantID:     param.MerchantID,
		CustID:         param.CustID,
		RRN:            param.RRN,
		TransactionId:  param.TrxID,
		ProductCode:    param.ProductCode,
		Amount:         int64(param.Amount),
		TransType:      constants.CODE_TRANSTYPE_REDEMPTION,
		IsUsed:         true,
		ProductType:    param.ProductType,
		Status:         saveStatus,
		ExpDate:        utils.DefaultNulTime(ExpireDate),
		Institution:    param.InstitutionID,
		CummulativeRef: param.CumReffnum,
		DateTime:       utils.GetTimeFormatYYMMDDHHMMSS(),
		Point:          param.Point,
		ResponderRc:    param.DataSupplier.Rc,
		ResponderRd:    param.DataSupplier.Rd,
		// RequestorData:     string(reqOttoag),
		// ResponderData:     string(responseOttoag),
		RequestorOPData:   string(reqdataOP),
		SupplierID:        param.SupplierID,
		RedeemAt:          utils.DefaultNulTime(time.Now()),
		CampaignId:        param.CampaignID,
		VoucherCode:       codeVoucher,
		CouponId:          param.CouponID,
		AccountId:         param.AccountId,
		ProductCategoryID: param.CategoryID,
		Comment:           param.Comment,
		MRewardID:         &param.RewardID,
		MProductID:        &param.ProductID,
		PointsTransferID:  param.PointTransferID,
		UsedAt:            utils.DefaultNulTime(time.Now()),

		PaymentMethod: 2,
		InvoiceNumber: param.InvoiceNumber,
	}

	err := db.DbCon.Create(&save).Error
	if err != nil {

		logrus.Error("[PackageRedeemtion_V21_Services]-[SaveTrxRedeemtionDUmmy]")
		logrus.Error(fmt.Sprintf("[SaveTrxRedeemtionDUmmy]-[Error : %v]", err))
		logrus.Println(fmt.Sprintf("[TransType : %v || TrxID : %v]", param.TransType, param.TrxID))

		name := jodaTime.Format("dd-MM-YYYY", time.Now()) + ".csv"
		go utils.CreateCSVFile(save, name)

		return "Gagal Save"

	}

	return "Berhasil Save"
}
