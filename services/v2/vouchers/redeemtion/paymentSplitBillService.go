package redeemtion

import (
	"encoding/json"
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	lp "ottopoint-purchase/hosts/landing_page/host"
	lpmodels "ottopoint-purchase/hosts/landing_page/models"
	"ottopoint-purchase/models"
	"ottopoint-purchase/models/dbmodels"
	"ottopoint-purchase/utils"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vjeantet/jodaTime"
)

func PaymentSplitBillServices(req models.PaymentSplitBillReq, param models.Params, balancePoint, balanceAmount int64) models.Response {

	res := models.Response{}

	nameservice := "[PackageRedeemtion]-[RedeemtionOttoAGServices]"
	logReq := fmt.Sprintf("[AccountNumber : %v,CampaignID : %v, FieldValue : %v]", param.AccountNumber, req.CampaignId, req.FieldValue)

	logrus.Info(nameservice)

	param.TrxID = utils.GenTransactionId()

	reqLG := lpmodels.LGRequestPay{
		Customerdetails: lpmodels.DataCustomerdetails{
			Email:     param.Email,
			Firstname: param.FirstName,
			Lastname:  param.LastName,
			Phone:     param.AccountNumber,
		},
		Transactiondetails: lpmodels.DataTransactiondetails{
			Amount:   int(balanceAmount),
			Currency: "idr",
			// Merchantname: ,
			Orderid: param.TrxID,
			// PaymentMethod :
			// Promocode   :
			// Vabca       :
			// Valain      :
			// Vamandiri   :
		},
	}

	landingPage, errLP := lp.PaymentLandingPage(reqLG)
	if errLP != nil || landingPage.ResponseData.StatusCode != "00" {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[PaymentLandingPage]-[Error : %v]", errLP))
		logrus.Println(logReq)

		save := savePaymentSplitBill(landingPage, reqLG, req, param, balancePoint, balanceAmount, constants.Success)
		if save != constants.KeyResponseSucceed {

			logrus.Error(nameservice)
			logrus.Error(fmt.Sprintf("[savePaymentSplitBill]-[Error : %v]", save))
			logrus.Println(logReq)

			res = utils.GetMessageFailedErrorNew(res, 500, "Internal Server Error")

			return res

		}

		res = utils.GetMessageFailedErrorNew(res, constants.RC_ERROR_FAILED_REDEEM_VOUCHER, constants.RD_ERROR_FAILED_REDEEM_VOUCHER)

		return res
	}

	param.RRN = landingPage.ResponseData.OrderID

	save := savePaymentSplitBill(landingPage, reqLG, req, param, balancePoint, balanceAmount, constants.Success)
	if save != constants.KeyResponseSucceed {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[savePaymentSplitBill]-[Error : %v]", save))
		logrus.Println(logReq)

		res = utils.GetMessageFailedErrorNew(res, 500, "Internal Server Error")

		return res

	}

	res = models.Response{
		Data: models.PaymentSplitBillResp{
			Code:       "00",
			Message:    "Success",
			Success:    1,
			Failed:     0,
			Pending:    0,
			UrlPayment: landingPage.ResponseData.EndpointURL,
		},
		Meta: utils.ResponseMetaOK(),
	}

	return res

}

func savePaymentSplitBill(resVendor lpmodels.LGResponsePay, reqVendor lpmodels.LGRequestPay, req models.PaymentSplitBillReq, param models.Params, balancePoint, balanceAmount int64, status string) string {

	nameservice := "[PackageRedeemtionServices]-[SavePaymentSplitBill]"
	logReq := fmt.Sprintf("[AccountNumber : %v, CampaignID : %v, FieldValue : %v]", param.AccountNumber, req.CampaignId, req.FieldValue)

	logrus.Info(nameservice)

	var saveStatus string
	var isUsed bool
	switch status {
	case "00":
		saveStatus = constants.Success
	case "09":
		saveStatus = constants.Pending
	case "01":
		saveStatus = constants.Failed
		isUsed = true
	}

	tspendingID := utils.GenerateTokenUUID()

	paymentCash := dbmodels.TPayment{
		ID:             utils.GenerateTokenUUID(),
		TSpendingID:    tspendingID,
		ExternalReffId: resVendor.ResponseData.OrderID,
		TransType:      constants.PaymentSplitBill,
		Value:          balanceAmount,
		ValueType:      "cash",
		Status:         saveStatus,
		ResponderRc:    resVendor.ResponseCode,
		ResponderRd:    resVendor.ResponseDesc,
		CreatedBy:      "System",
		// UpdatedBy        : ,
		CreatedAt: time.Now(),
		// UpdatedAt        : ,
	}

	errPayment := db.DbCon.Create(&paymentCash).Error
	if errPayment != nil {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[SaveTPaymentCash]-[Error : %v]", errPayment))
		logrus.Println(logReq)

		return fmt.Sprintf("%v", errPayment)

	}

	paymentPoint := dbmodels.TPayment{
		ID:          utils.GenerateTokenUUID(),
		TSpendingID: tspendingID,
		// ExternalReffId: param.TrxID,
		TransType: constants.SpendingSplitBill,
		Value:     balancePoint,
		ValueType: "point",
		Status:    saveStatus,
		// ResponderRc:    resVendor.ResponseCode,
		// ResponderRd:    resVendor.ResponseDesc,
		CreatedBy: "System",
		// UpdatedBy        : ,
		CreatedAt: time.Now(),
		// UpdatedAt        : ,
	}

	errPayment = db.DbCon.Create(&paymentPoint).Error
	if errPayment != nil {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[SaveTPaymentPoint]-[Error : %v]", errPayment))
		logrus.Println(logReq)

		return fmt.Sprintf("%v", errPayment)

	}

	reqLP, _ := json.Marshal(&reqVendor) // Req Vendor
	resLP, _ := json.Marshal(&resVendor) // Response Vendor
	reqOP, _ := json.Marshal(&req)       // Req Service

	spendingReq := dbmodels.TSpending{
		ID:            tspendingID,
		AccountNumber: param.AccountNumber,
		Voucher:       param.NamaVoucher,
		MerchantID:    param.MerchantID,
		CustID:        param.CustID,
		RRN:           param.RRN,
		TransactionId: param.TrxID,
		ProductCode:   param.ProductCode,
		Amount:        int64(param.Amount),
		TransType:     param.TransType,
		IsUsed:        isUsed,
		ProductType:   param.ProductType,
		Status:        saveStatus,
		// ExpDate:         param.ExpDate,
		// ExpDate:           utils.DefaultNulTime(ExpireDate),
		Institution:     param.InstitutionID,
		CummulativeRef:  param.Reffnum,
		DateTime:        utils.GetTimeFormatYYMMDDHHMMSS(),
		Point:           param.Point,
		ResponderRc:     param.DataSupplier.Rc,
		ResponderRd:     param.DataSupplier.Rd,
		RequestorData:   string(reqLP),
		ResponderData:   string(resLP),
		RequestorOPData: string(reqOP),
		SupplierID:      param.SupplierID,
		// RedeemAt:          utils.DefaultNulTime(redeemDate),
		CampaignId: param.CampaignID,
		// VoucherCode:       codeVoucher,
		CouponId:          param.CouponID,
		AccountId:         param.AccountId,
		ProductCategoryID: param.CategoryID,
		Comment:           param.Comment,
		MRewardID:         param.RewardID,
		MProductID:        param.ProductID,
		PointsTransferID:  param.PointTransferID,
		// UsedAt:            utils.DefaultNulTime(usedAt),
	}

	errSpending := db.DbCon.Create(&spendingReq).Error
	if errSpending != nil {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[SaveTSpending]-[Error : %v]", errSpending))
		logrus.Println(logReq)

		name := jodaTime.Format("dd-MM-YYYY", time.Now()) + ".csv"
		go utils.CreateCSVFile(spendingReq, name)

		return fmt.Sprintf("%v", errSpending)

	}

	return constants.KeyResponseSucceed

}
