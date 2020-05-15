package services

import (
	"fmt"
	"ottopoint-purchase/constants"
	db "ottopoint-purchase/db"
	opl "ottopoint-purchase/hosts/opl/host"
	"ottopoint-purchase/models"
	"ottopoint-purchase/services/voucher"
	"strings"
	"sync"
)

// func UseVoucherService1(header models.RequestHeader, req models.VoucherComultaiveReq, dataToken ottomartmodels.ResponseToken, amount int64, rrn string) models.Response {
func UseVoucherComulative(req models.VoucherComultaiveReq, redeemComu models.RedeemComuResp, param models.Params, getRespChan chan models.RedeemResponse, ErrRespUseVouc chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(getRespChan)
	defer close(ErrRespUseVouc)
	// resRedeemComu := models.RedeemComuResp{}
	fmt.Println("[UseVoucherComulative]-[Package-Services]")

	// get CustID
	dataUser, errUser := db.CheckUser(param.AccountNumber)
	if errUser != nil {
		fmt.Println("User Belum Eligible, Error : ", errUser)
	} else {
		// Use Voucher to Openloyalty
		fmt.Println("campaignId : ", req.CampaignID)
		fmt.Println("couponId : ", redeemComu.CouponID)
		fmt.Println("code : ", redeemComu.CouponCode)
		fmt.Println("used : ", 1)
		fmt.Println("customerId : ", dataUser.CustID)

		_, err2 := opl.CouponVoucherCustomer(req.CampaignID, redeemComu.CouponID, redeemComu.CouponCode, dataUser.CustID, 1)
		fmt.Println("================ doing use voucher couponId : ", redeemComu.CouponID)
		if err2 != nil {
			fmt.Println("================ doing use voucher couponId Error: ", redeemComu.CouponID)

		} else {

			// Reedem Use Voucher
			param.Amount = redeemComu.Redeem.Amount
			param.RRN = redeemComu.Redeem.Rrn
			resRedeem := RedeemUseVoucherComulative(req, param)

			getRespChan <- resRedeem
		}
	}

}

// function reedem use voucher
func RedeemUseVoucherComulative(req models.VoucherComultaiveReq, param models.Params) models.RedeemResponse {
	res := models.RedeemResponse{}

	fmt.Println("[RedeemUseVoucherComulative]-[Package-Services]")

	category := strings.ToLower(param.Category)

	reqRedeem := models.UseRedeemRequest{
		AccountNumber: param.AccountNumber,
		CustID:        req.CustID,
		CustID2:       req.CustID2,
		ProductCode:   param.ProductCode,
		Jumlah:        param.Total,
	}

	resRedeem := models.UseRedeemResponse{}

	switch category {
	case constants.CategoryPulsa:
		resRedeem = voucher.RedeemPulsaComulative(reqRedeem, req, param)
	case constants.CategoryPLN:
		resRedeem = voucher.RedeemPLNComulative(reqRedeem, req, param)
	case constants.CategoryMobileLegend, constants.CategoryFreeFire:
		resRedeem = voucher.RedeemGameComulative(reqRedeem, req, param)
	}

	res = models.RedeemResponse{
		Rc:          resRedeem.Rc,
		Rrn:         resRedeem.Rrn,
		CustID:      resRedeem.CustID,
		ProductCode: resRedeem.ProductCode,
		Amount:      resRedeem.Amount,
		Msg:         resRedeem.Msg,
		Uimsg:       resRedeem.Uimsg,
		Datetime:    resRedeem.Datetime,
		Data:        resRedeem.Data,
	}

	return res

}
