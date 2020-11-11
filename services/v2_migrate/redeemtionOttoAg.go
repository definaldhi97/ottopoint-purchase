package v2_migrate

import (
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/models"
	ottoagmodels "ottopoint-purchase/models/ottoag"
	"ottopoint-purchase/services"
	biller "ottopoint-purchase/services/ottoag"
	"ottopoint-purchase/utils"
)

func RedeemVoucherOttoAg(req models.VoucherComultaiveReq, param models.Params, getResp chan models.RedeemComuResp, ErrRespRedeem chan error) {
	fmt.Println("[ >>>>>>>>>>>>>>>>>> Redeemtion Comulative Voucher Otto AG <<<<<<<<<<<<<<<< ]")
	fmt.Println("[ Inquery OttoAG ] - [ Deduct point OPL & Deduct Voucher ]")

	defer close(getResp)
	defer close(ErrRespRedeem)

	resRedeemComu := models.RedeemComuResp{}
	redeemRes := models.RedeemComuResp{
		Code: "00",
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

	fmt.Println("[INQUIRY-BILLER][START]")
	dataInquery, errInquiry := biller.InquiryBiller(inqReq.Data, req, reqInq, param)

	paramInq := models.Params{
		AccountNumber: param.AccountNumber,
		MerchantID:    param.MerchantID,
		InstitutionID: param.InstitutionID,
		CustID:        req.CustID,
		TransType:     constants.CODE_TRANSTYPE_INQUERY,
		Reffnum:       param.Reffnum, // internal
		RRN:           dataInquery.Rrn,
		TrxID:         param.TrxID,
		Amount:        dataInquery.Amount,
		NamaVoucher:   param.NamaVoucher,
		ProductType:   param.ProductType,
		ProductCode:   param.ProductCode,
		Category:      param.Category,
		Point:         param.Point,
		ExpDate:       param.ExpDate,
		SupplierID:    param.SupplierID,
		CategoryID:    param.CategoryID,
		DataSupplier: models.Supplier{
			Rc: dataInquery.Rc,
			Rd: dataInquery.Msg,
		},
	}

	if dataInquery.Rc != constants.CODE_SUCCESS {
		fmt.Println("[Error-DataInquiry]-[Redeem Comulative Voucher Otto AG]")
		fmt.Println("[Error : %v]", errInquiry)

		redeemRes = models.RedeemComuResp{
			Code:    "01",
			Message: "Inquiry Failed",
		}

		go services.SaveTransactionOttoAg(paramInq, dataInquery, reqInq, req, constants.CODE_FAILED)

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

	// Time Out
	if dataInquery.Rc == "" {
		fmt.Println("[Error-DataInquiry]-[Redeem Comulative Voucher Otto AG]")
		fmt.Println("[Error : %v]", errInquiry)
		redeemRes = models.RedeemComuResp{
			Code:    "01",
			Message: "Inquiry Failed",
		}

		go services.SaveTransactionOttoAg(paramInq, dataInquery, reqInq, req, constants.CODE_FAILED)

		ErrRespRedeem <- errInquiry

		resRedeemComu.Redeem.Rc = "01"
		resRedeemComu.Redeem.Rc = "Time Out"
		resRedeemComu.Code = redeemRes.Code
		resRedeemComu.Message = redeemRes.Message

		getResp <- resRedeemComu

		return

	}

	go services.SaveTransactionOttoAg(paramInq, dataInquery, reqInq, req, constants.CODE_SUCCESS)

	// deduct point and deduct usage_limit voucher
	resultRedeemVouch, errRedeemVouch := Redeem_PointandVoucher(req.Jumlah, param)
	fmt.Println("Response Deduct point dan voucher")
	fmt.Println(resultRedeemVouch)

	if resultRedeemVouch.Rc != "00" {
		fmt.Println("[ Error Redeem_PointandVoucher]")
		fmt.Println(resultRedeemVouch.Rd)

		redeemRes = models.RedeemComuResp{
			Code:    "01",
			Message: "Gagal Redeem",
		}

		ErrRespRedeem <- errRedeemVouch

		resRedeemComu.Redeem.Rc = "01"
		resRedeemComu.Redeem.Msg = "Maaf transaksi Anda tidak dapat dilakukan saat ini. Silahkan dicoba lagi atau hubungi tim kami untuk informasi selengkapnya"
		resRedeemComu.Code = redeemRes.Code
		resRedeemComu.Message = redeemRes.Message

		getResp <- resRedeemComu

		return
	}

	resRedeemComu.CouponCode = resultRedeemVouch.CouponsCode

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

}
