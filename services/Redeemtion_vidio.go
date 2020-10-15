package services

// import (
// 	"fmt"
// 	"ottopoint-purchase/constants"
// 	"ottopoint-purchase/models"
// 	ottoagmodels "ottopoint-purchase/models/ottoag"
// 	biller "ottopoint-purchase/services/ottoag"
// 	"ottopoint-purchase/utils"

// 	"github.com/opentracing/opentracing-go"
// 	"go.uber.org/zap"
// )

// type RedeemtionVidio struct {
// 	General models.GeneralModel
// }

// func (t RedeemtionVidio) RedeemtioVidioService(req models.VoucherComultaiveReq, param models.Params) models.Response {
// 	var res models.Response

// 	fmt.Println("[ Redeemtion Vidio Service ]")

// 	sugarLogger := t.General.OttoZaplog
// 	sugarLogger.Info("[RedeemtioVidioService]",
// 		zap.String("NameVoucher : ", param.NamaVoucher), zap.Int("Jumlah : ", req.Jumlah),
// 		zap.String("CampaignID : ", req.CampaignID), zap.String("CampaignID : ", req.CampaignID),
// 		zap.String("CustID2 : ", req.CustID2), zap.String("ProductCode : ", param.ProductCode),
// 		// zap.Int("Point : ", req.Point),
// 		zap.String("AccountNumber : ", param.AccountNumber), zap.String("InstitutionID : ", param.InstitutionID))

// 	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[RedeemtioVidioService]")
// 	defer span.Finish()

// 	// Inquery biller ottoag
// 	fmt.Println("[ Start Inquery Vidio ]")

// 	inqBiller := ottoagmodels.BillerInquiryDataReq{
// 		ProductCode: param.ProductCode,
// 		MemberID:    utils.MemberID,
// 		CustID:      req.CustID,
// 		Period:      req.CustID2,
// 	}

// 	inqReq := ottoagmodels.OttoAGInquiryRequest{
// 		TypeTrans:     "0003",
// 		Datetime:      utils.GetTimeFormatYYMMDDHHMMSS(),
// 		IssuerID:      param.InstitutionID,
// 		AccountNumber: param.AccountNumber,
// 		Data:          inqBiller,
// 	}

// 	reqInq := models.UseRedeemRequest{
// 		AccountNumber: param.AccountNumber,
// 		CustID:        req.CustID,
// 		CustID2:       req.CustID2,
// 		ProductCode:   param.ProductCode,
// 	}

// 	fmt.Println("[ Send Inquery OttoAg ]")
// 	dataInquery, errInquiry := biller.InquiryBiller(inqReq.Data, req, reqInq, param)

// 	paramInq := models.Params{
// 		AccountNumber: param.AccountNumber,
// 		MerchantID:    param.MerchantID,
// 		InstitutionID: param.InstitutionID,
// 		CustID:        req.CustID,
// 		TransType:     constants.CODE_TRANSTYPE_INQUERY,
// 		Reffnum:       param.Reffnum, // internal
// 		RRN:           dataInquery.Rrn,
// 		Amount:        dataInquery.Amount,
// 		NamaVoucher:   param.NamaVoucher,
// 		ProductType:   param.ProductType,
// 		ProductCode:   param.ProductCode,
// 		Category:      param.Category,
// 		Point:         param.Point,
// 		ExpDate:       param.ExpDate,
// 		SupplierID:    param.SupplierID,
// 		DataSupplier: models.Supplier{
// 			Rc: dataInquery.Rc,
// 			Rd: dataInquery.Msg,
// 		},
// 	}

// 	if dataInquery.Rc != "00" {
// 		fmt.Println("[Error-Data Inquiry]-[Redeem Vidio]")
// 		fmt.Println("[Error : %v]", errInquiry)

// 		// redeemRes = models.RedeemComuResp{
// 		// 	Code:    "01",
// 		// 	Message: "Inquiry Failed",
// 		// }

// 		go saveTransactionOttoAg(paramInq, dataInquery, reqInq, req, "01")
// 	}

// }
