package services

import (
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	"ottopoint-purchase/hosts/opl/host"
	"ottopoint-purchase/models"
	ottoagmodels "ottopoint-purchase/models/ottoag"
	biller "ottopoint-purchase/services/ottoag"
	"ottopoint-purchase/services/voucher"
	"ottopoint-purchase/utils"
	"strings"
)

func RedeemComulativeVoucher(req models.VoucherComultaiveReq, param models.Params, getResp chan models.RedeemComuResp, ErrRespRedeem chan error) {
	defer close(getResp)
	defer close(ErrRespRedeem)

	resRedeemComu := models.RedeemComuResp{}
	redeemRes := models.RedeemComuResp{
		Code: "00",
	}

	category := strings.ToLower(param.Category)

	fmt.Println("[Start][Inquiry]-[Package-Services]-[RedeemComulativeVoucher]")

	if category == constants.CategoryPulsa {
		// validate prefix
		validate, errValidate := ValidatePrefixComulative(req.CustID, param.ProductCode)
		if validate == false {

			fmt.Println("Invalid Prefix")

			redeemRes = models.RedeemComuResp{
				Code:    "01",
				Message: "Invalid Prefix",
			}

			ErrRespRedeem <- errValidate

			resRedeemComu.Code = redeemRes.Code
			resRedeemComu.Message = redeemRes.Message
			resRedeemComu.Redeem.Rc = "16"
			resRedeemComu.Redeem.Msg = "Nomor yang kamu masukkan salah"

			getResp <- resRedeemComu

			return
		}
	}

	// ==========Inquery OttoAG==========

	inqBiller := ottoagmodels.BillerInquiryDataReq{
		ProductCode: param.ProductCode,
		MemberID:    utils.MemberID,
		CustID:      req.CustID,
		Period:      req.CustID2,
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
		ProductCode:   param.ProductCode,
	}

	// if !ottoag.ValidateDataInq(inqReq) {
	// 	fmt.Println("[Error-DataInquiry]-[RedeemComulativeVoucher]")
	// 	fmt.Println("[Error ValidateDataInq]")
	// 	var err error
	// 	redeemRes = models.RedeemComuResp{
	// 		Code:    "01",
	// 		Message: "Inquiry Gagal",
	// 	}

	// 	// go voucher.SaveTransactionPulsa(paramInq, "Inquiry", "01")

	// 	ErrRespRedeem <- err

	// 	resRedeemComu.Code = redeemRes.Code
	// 	resRedeemComu.Message = redeemRes.Message
	// 	resRedeemComu.Redeem.Rc =
	// 	resRedeemComu.Redeem.Msg = "Invalid "

	// 	getResp <- resRedeemComu

	// 	return
	// }

	fmt.Println("[INQUIRY-BILLER][START]")
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
		ProductCode:   param.ProductCode,
		Category:      param.Category,
		Point:         param.Point,
		ExpDate:       param.ExpDate,
		SupplierID:    param.SupplierID,
	}

	if dataInquery.Rc != "00" {
		fmt.Println("[Error-DataInquiry]-[RedeemComulativeVoucher]")
		fmt.Println("[Error : %v]", errInquiry)
		redeemRes = models.RedeemComuResp{
			Code:    "01",
			Message: "Inquiry Failed",
		}

		go voucher.SaveTransactionPulsa(paramInq, dataInquery, req, inqBiller, "Inquiry", "01", dataInquery.Rc)

		ErrRespRedeem <- errInquiry

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

		resRedeemComu.Redeem = r
		resRedeemComu.Code = redeemRes.Code
		resRedeemComu.Message = redeemRes.Message

		getResp <- resRedeemComu

		return

	}

	go voucher.SaveTransactionPulsa(paramInq, dataInquery, req, inqBiller, "Inquiry", "00", dataInquery.Rc)

	// coupon := []models.CouponsRedeem{}

	fmt.Println("[Start][Redeem]-[Package-Services]-[RedeemComulativeVoucher]")
	data, errx := host.RedeemVoucher(req.CampaignID, param.AccountNumber)

	if errx != nil || data.Message != "" || len(data.Coupons) == 0 {

		fmt.Println("[ErrorRedeemVoucher]-[RedeemComulativeVoucher]")
		fmt.Println(fmt.Sprintf("Error : %v", errx))

		redeemRes = models.RedeemComuResp{
			Code:    "01",
			Message: "Gagal Redeem",
		}

		ErrRespRedeem <- errx

		resRedeemComu.Redeem.Rc = "01"
		resRedeemComu.Redeem.Msg = "Maaf transaksi Anda tidak dapat dilakukan saat ini. Silahkan dicoba lagi atau hubungi tim kami untuk informasi selengkapnya"
		resRedeemComu.Code = redeemRes.Code
		resRedeemComu.Message = redeemRes.Message

		getResp <- resRedeemComu

		return

	}
	// else {
	// 	fmt.Println("Response Redeem 1 : ", data)
	// 	fmt.Println("Response LEN Coupons 1 : ", len(data.Coupons))
	// 	fmt.Println("check data ", data == nil)

	// 	//fmt.Println("Response check Coupons : ", len(coupon))
	// 	// check if no data founded
	// 	if len(data.Coupons) == 0 {
	// 		fmt.Println("========== check coupon ", coupon)
	// 		redeemRes = models.RedeemComuResp{
	// 			Code:    "01",
	// 			Message: "Anda mencapai batas maksimal pembelian voucher",
	// 		}
	// 	} else {
	// resRedeemComu.CouponID = data.Coupons[0].Id
	// resRedeemComu.CouponCode = data.Coupons[0].Code
	// 	}
	// }

	resRedeemComu.CouponID = data.Coupons[0].Id
	resRedeemComu.CouponCode = data.Coupons[0].Code

	ErrRespRedeem <- nil

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

	// validate panjang nomor, Jika nomor kurang dari 4
	if len(custID) < 4 {

		fmt.Println("[Kurang dari 4]-[Prefix-ottopoint]-[RedeemComulativeVoucher]")
		fmt.Println(fmt.Sprintf("invalid Prefix %v", custID))

		return false, err
	}

	// validate panjang nomor, Jika nomor kurang dari 11 & lebih dari 15
	if len(custID) <= 10 || len(custID) > 15 {

		fmt.Println("[Kurang dari 10 atau lebih dari 15]-[Prefix-ottopoint]-[RedeemComulativeVoucher]")
		fmt.Println(fmt.Sprintf("invalid Prefix %v", custID))

		return false, err

	}

	// get Prefix
	dataPrefix, errPrefix := db.GetOperatorCodebyPrefix(custID)
	if errPrefix != nil {

		fmt.Println("[ErrorPrefix]-[RedeemComulativeVoucher]")
		fmt.Println(fmt.Sprintf("dataPrefix = %v", dataPrefix))
		fmt.Println(fmt.Sprintf("Prefix tidak ditemukan %v", errPrefix))

		return false, err
	}

	// check operator by OperatorCode
	prefix := utils.Operator(dataPrefix.OperatorCode)
	// check operator by ProductCode
	product := utils.ProductPulsa(productCode[0:4])

	// Jika Nomor tidak sesuai dengan operator
	if prefix != product {

		fmt.Println("[Operator]-[Prefix-ottopoint]-[RedeemComulativeVoucher]")
		fmt.Println(fmt.Sprintf("invalid Prefix %v", prefix))

		return false, err

	}

	return true, nil
}
