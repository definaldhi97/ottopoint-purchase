package payment

import (
	"errors"
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	"ottopoint-purchase/models"
	"ottopoint-purchase/models/dbmodels"
	sp "ottopoint-purchase/models/v2/payment"
	"ottopoint-purchase/utils"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vjeantet/jodaTime"
)

func SpendingPaymentService(req sp.SpendingPaymentReq, param models.Params, header models.RequestHeader) models.Response {
	res := models.Response{}

	nameservice := "[PackagePayment]-[SpendingPaymentService]"

	invoiceNumber := "INV" + jodaTime.Format("YYYYMMdd", time.Now()) + utils.GenTransactionId()[7:11]

	logReq := fmt.Sprintf("[AccountNumber : %v || ReferenceId : %v || InvoiceNumber : %v]", req.AccountNumber, req.ReferenceId, invoiceNumber)

	logrus.Info(nameservice)

	_, errCheck := db.CheckReffIdSplitBill(req.ReferenceId)
	if errCheck == nil {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[CheckReffIdSplitBill]-[Duplicate Reference ID]"))
		logrus.Println(logReq)

		res = utils.GetMessageResponse(res, 172, false, errors.New("Duplicate Reference ID"))

		return res
	}

	param.TrxID = utils.GenTransactionId()
	idSpending := utils.GenerateUUID()
	param.Comment = param.TrxID + header.InstitutionID + "#" + req.ProductName

	spend, errSpend := SpendPointService(param, header)
	param.PointTransferID = spend.PointTransferID
	if errSpend != nil || spend.PointTransferID == "" {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[SpendPointService]-[Error : %v]", errSpend))
		logrus.Println(logReq)

		go saveTSpending(req, param, idSpending, invoiceNumber, constants.Failed)

		res = utils.GetMessageResponse(res, 209, false, errors.New("Gagal Spend Point"))

		return res
	}

	go saveTSpending(req, param, idSpending, invoiceNumber, constants.Success)

	if req.TransType == constants.CodeSplitBill {
		go saveTPayemnt(req, param, idSpending, constants.Success)
	}

	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: map[string]interface{}{
			"referenceId":   req.ReferenceId,
			"invoiceNumber": invoiceNumber,
		},
	}

	return res

}

func saveTSpending(req sp.SpendingPaymentReq, param models.Params, idSpending, invoiceNumber, status string) {

	nameservice := "[PackagePayment]-[SaveTSpending]"
	logReq := fmt.Sprintf("[AccountNumber : %v || ReferenceId : %v]", req.AccountNumber, req.ReferenceId)

	redeem := time.Now()

	save := dbmodels.TSpending{
		ID:            idSpending,
		AccountNumber: param.AccountNumber,
		Voucher:       req.ProductName,
		// MerchantID       : ,
		// CustID           : ,
		RRN:           param.RRN,
		TransactionId: param.TrxID,
		// ProductCode      : ,
		Amount:    int64(req.Amount),
		TransType: req.TransType,
		// IsUsed           : ,
		// ProductType      : ,
		Status: status,
		// ExpDate          : ,
		Institution: param.InstitutionID,
		// CummulativeRef   : ,
		DateTime: req.TransactionTime,
		// ResponderData    : ,
		Point: req.Point,
		// ResponderRc      : ,
		// ResponderRd      : ,
		// RequestorData    : ,
		// RequestorOPData  : ,
		// SupplierID       : ,
		// CouponId         : ,
		// CampaignId       : ,
		AccountId: param.AccountId,
		RedeemAt:  &redeem,
		UsedAt:    &redeem,
		CreatedAT: time.Now(),
		// UpdatedAT        : ,
		// VoucherCode      : ,
		// ProductCategoryID: ,
		Comment: req.Comment,
		// MRewardID        : ,
		// MProductID       : ,
		// VoucherLink      : ,
		PointsTransferID: param.PointTransferID,
		InvoiceNumber:    invoiceNumber, // INVYYYYMMDDXXXX > 15 digit total
		PaymentMethod:    req.PaymentMethod,
	}

	err := db.DbCon.Create(&save).Error
	if err != nil {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[saveTPayemnt]-[Error : %v]", err))
		logrus.Println(logReq)

		return

	}

	return

}

func saveTPayemnt(req sp.SpendingPaymentReq, param models.Params, idSpending, status string) {

	nameservice := "[PackagePayment]-[SaveTPayemnt]"
	logReq := fmt.Sprintf("[ReferenceId : %v]", req.ReferenceId)

	savePoint := dbmodels.TPayment{
		ID:             utils.GenerateUUID(),
		TSpendingID:    idSpending,
		ExternalReffId: param.TrxID,
		TransType:      req.TransType,
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

	errPoint := db.DbCon.Create(&savePoint).Error
	if errPoint != nil {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[SaveTPayemntPoint]-[Error : %v]", errPoint))
		logrus.Println(logReq)

		// return

	}

	saveCash := dbmodels.TPayment{
		ID:             utils.GenerateUUID(),
		TSpendingID:    idSpending,
		ExternalReffId: req.ReferenceId,
		TransType:      req.TransType,
		Value:          int64(req.Cash),
		ValueType:      constants.TypeCash,
		Status:         status,
		// ResponderRc   : ,
		// ResponderRd   : ,
		CreatedBy: constants.CreatedbySystem,
		// UpdatedBy     : ,
		CreatedAt: time.Now(),
		// UpdatedAt     : ,
	}

	errCash := db.DbCon.Create(&saveCash).Error
	if errCash != nil {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[SaveTPayemntCash]-[Error : %v]", errCash))
		logrus.Println(logReq)

		// return

	}

	return

}
