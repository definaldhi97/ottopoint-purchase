package services

import (
	"fmt"
	"ottopoint-purchase/db"
	"ottopoint-purchase/hosts/opl/host"
	"ottopoint-purchase/models"
	ottoagmodels "ottopoint-purchase/models/ottoag"
	"ottopoint-purchase/services/ottoag"
	biller "ottopoint-purchase/services/ottoag"
	"ottopoint-purchase/services/voucher"
	"ottopoint-purchase/utils"

	"github.com/astaxie/beego/logs"
)

func RedeemComulativeVoucher(req models.VoucherComultaiveReq, param models.Params, getResp chan models.RedeemComuResp, ErrRespRedeem chan error) {
	defer close(getResp)
	defer close(ErrRespRedeem)

	resRedeemComu := models.RedeemComuResp{}
	redeemRes := models.RedeemComuResp{
		Code: "00",
	}

	logs.Info("[Start]-[Package-Services]-[RedeemComulativeVoucher]")

	// validate prefix
	validate, errValidate := ValidatePrefixComulative(req.CustID, req.ProductCode)
	if validate == false {

		redeemRes = models.RedeemComuResp{
			Code:    "01",
			Message: "Invalid Prefix",
		}

		ErrRespRedeem <- errValidate

		resRedeemComu.Code = redeemRes.Code
		resRedeemComu.Message = redeemRes.Message

		return
	}

	// ==========Inquery OttoAG==========

	inqBiller := ottoagmodels.BillerInquiryDataReq{
		ProductCode: req.ProductCode,
		MemberID:    utils.MemberID,
		CustID:      req.CustID,
	}

	inqReq := ottoagmodels.OttoAGInquiryRequest{
		TypeTrans:     "0003",
		Datetime:      utils.GetTimeFormatYYMMDDHHMMSS(),
		IssuerID:      "OTTOPAY",
		AccountNumber: param.AccountNumber,
		Data:          inqBiller,
	}

	reqInq := models.UseRedeemRequest{
		AccountNumber: param.AccountNumber,
		CustID:        req.CustID,
		CustID2:       req.CustID2,
		ProductCode:   req.ProductCode,
	}

	if !ottoag.ValidateDataInq(inqReq) {
		logs.Info("[Error-DataInquiry]-[RedeemComulativeVoucher]")
		logs.Info("[Error ValidateDataInq]")
		var err error
		redeemRes = models.RedeemComuResp{
			Code:    "01",
			Message: "Inquiry PreFafix",
		}

		// go voucher.SaveTransactionPulsa(paramInq, "Inquiry", "01")

		ErrRespRedeem <- err

		resRedeemComu.Code = redeemRes.Code
		resRedeemComu.Message = redeemRes.Message

		return
	}

	logs.Info("[INQUIRY-BILLER][START]")
	dataInquery, errInquiry := biller.InquiryBiller(inqReq.Data, req, reqInq, param)

	paramInq := models.Params{
		AccountNumber: param.AccountNumber,
		MerchantID:    param.MerchantID,
		InstitutionID: param.InstitutionID,
		CustID:        req.CustID,
		Reffnum:       param.Reffnum, // internal
		RRN:           dataInquery.Rrn,
		Amount:        dataInquery.Amount,
		NamaVoucher:   param.NamaVoucher,
		ProductType:   param.ProductType,
		ProductCode:   req.ProductCode,
		Category:      param.Category,
		Point:         param.Point,
		ExpDate:       param.ExpDate,
		SupplierID:    param.SupplierID,
	}

	if dataInquery.Rc == "01" {
		logs.Info("[Error-DataInquiry]-[RedeemComulativeVoucher]")
		logs.Info("[Error : %v]", errInquiry)
		redeemRes = models.RedeemComuResp{
			Code:    "01",
			Message: "Invalid Prefix",
		}

		go voucher.SaveTransactionPulsa(paramInq, "Inquiry", "01")

		ErrRespRedeem <- errInquiry

		resRedeemComu.Code = redeemRes.Code
		resRedeemComu.Message = redeemRes.Message

		return

	}

	coupon := []models.CouponsRedeem{}

	logs.Info("============= Redeem voucher 1 =============")
	data, errx := host.RedeemVoucher(req.CampaignID, param.AccountNumber)

	if errx != nil {

		logs.Info("[ErrorRedeemVoucher]-[RedeemComulativeVoucher]")
		logs.Info(fmt.Sprintf("Error : %v", errx))

		redeemRes = models.RedeemComuResp{
			Code:    "01",
			Message: "Internal Server Error",
		}

		resRedeemComu.Code = redeemRes.Code
		resRedeemComu.Message = redeemRes.Message

		ErrRespRedeem <- errx

		return

	} else {
		logs.Info("Response Redeem 1 : ", data)
		logs.Info("Response LEN Coupons 1 : ", len(data.Coupons))
		fmt.Println("check data ", data == nil)

		//logs.Info("Response check Coupons : ", len(coupon))
		// check if no data founded
		if len(data.Coupons) == 0 {
			logs.Info("========== check coupon ", coupon)
			redeemRes = models.RedeemComuResp{
				Code:    "01",
				Message: "Anda mencapai batas maksimal pembelian voucher",
			}
		} else {
			resRedeemComu.CouponID = data.Coupons[0].Id
			resRedeemComu.CouponCode = data.Coupons[0].Code
		}
	}

	logs.Info("========== rrn ", dataInquery.Rrn)
	logs.Info("========== Amount ", dataInquery.Amount)

	ErrRespRedeem <- nil

	go voucher.SaveTransactionPulsa(paramInq, "Inquiry", "00")

	r := models.RedeemResponse{
		Rc:          dataInquery.Rc,
		Rrn:         dataInquery.Rrn,
		CustID:      dataInquery.CustID,
		ProductCode: dataInquery.ProductCode,
		Amount:      dataInquery.Amount,
		Msg:         dataInquery.Msg,
		Uimsg:       dataInquery.Uimsg,
		// Datetime:    time.Now(),
		Data: dataInquery.Data,
	}

	resRedeemComu.Code = redeemRes.Code
	resRedeemComu.Message = redeemRes.Message
	resRedeemComu.Redeem = r
	getResp <- resRedeemComu

	// return
}

func ValidatePrefixComulative(custID, productCode string) (bool, error) {

	var err error
	// get Prefix
	dataPrefix, errPrefix := db.GetOperatorCodebyPrefix(custID)
	if errPrefix != nil {

		logs.Info("[ErrorPrefix]-[RedeemComulativeVoucher]")
		logs.Info(fmt.Sprintf("dataPrefix = %v", dataPrefix))
		logs.Info(fmt.Sprintf("Prefix tidak ditemukan %v", errPrefix))

		return false, err
	}

	// check operator by OperatorCode
	prefix := utils.Operator(dataPrefix.OperatorCode)
	// check operator by ProductCode
	product := utils.ProductPulsa(productCode[0:4])

	// validate panjang nomor, Jika nomor kurang dari 4
	if len(custID) < 4 {

		logs.Info("[FAILED]-[Prefix-ottopoint]-[RedeemComulativeVoucher]")
		logs.Info(fmt.Sprintf("invalid Prefix %v", custID))

		return false, err
	}

	// validate panjang nomor, Jika nomor kurang dari 11 & lebih dari 15
	if len(custID) <= 10 || len(custID) > 15 {

		logs.Info("[FAILED]-[Prefix-ottopoint]-[RedeemComulativeVoucher]")
		logs.Info(fmt.Sprintf("invalid Prefix %v", custID))

		return false, err

	}

	// Jika Nomor tidak sesuai dengan operator
	if prefix != product {

		logs.Info("[FAILED]-[Prefix-ottopoint]-[RedeemComulativeVoucher]")
		logs.Info(fmt.Sprintf("invalid Prefix %v", prefix))

		return false, err

	}

	return true, nil
}
