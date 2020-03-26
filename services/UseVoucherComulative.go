package services

import (
	"errors"
	"fmt"
	"ottopoint-purchase/constants"
	db "ottopoint-purchase/db"
	opl "ottopoint-purchase/hosts/opl/host"
	"ottopoint-purchase/models"
	"ottopoint-purchase/redis"
	"ottopoint-purchase/services/voucher"
	"ottopoint-purchase/utils"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/vjeantet/jodaTime"
)

// func UseVoucherService1(header models.RequestHeader, req models.VoucherComultaiveReq, dataToken ottomartmodels.ResponseToken, amount int64, rrn string) models.Response {
func UseVoucherComulative(req models.VoucherComultaiveReq, redeemComu models.RedeemComuResp, param models.Params, getRespChan chan models.RedeemResponse, ErrRespUseVouc chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(getRespChan)
	defer close(ErrRespUseVouc)
	// resRedeemComu := models.RedeemComuResp{}
	logs.Info("=============== in Function_Usevoucher1")

	// get CustID
	dataUser, errUser := db.CheckUser(param.AccountNumber)
	if errUser != nil {
		logs.Info("=============== User belum Eligible")
	} else {
		// Use Voucher to Openloyalty
		_, err2 := opl.CouponVoucherCustomer(req.CampaignID, redeemComu.CouponID, redeemComu.CouponCode, dataUser.CustID, 1)
		logs.Info("================ doing use voucher couponId : ", redeemComu.CouponID)
		if err2 != nil {
			logs.Info("================ doing use voucher couponId Error: ", redeemComu.CouponID)

		} else {
			// save Timer voucher to redis usedAt

			date := jodaTime.Format("dd-MM-YYYY hh:mm:ss", time.Now())
			go SaveTimeVoucher(date, redeemComu.CouponID, param.AccountNumber)

			// Reedem Use Voucher
			param.Amount = redeemComu.Redeem.Amount
			param.RRN = redeemComu.Redeem.Rrn
			resRedeem := RedeemUseVoucherComulative(req, param)
			// validate response reedem
			go ValidateRespRedeem(resRedeem.Msg, req.CampaignID, redeemComu.CouponID, redeemComu.CouponCode, dataUser.CustID)
			getRespChan <- resRedeem
		}
	}

}

// save Timer voucher to redis usedAt
func SaveTimeVoucher(date, couponId, accountNumber string) {
	date = (time.Now().Local().Add(time.Hour * time.Duration(7))).Format("2006-01-02T15:04:05-0700")
	keyTimeVoucher := fmt.Sprintf("usedVoucherAt-%s-%s", couponId, accountNumber)
	go redis.SaveRedis(keyTimeVoucher, date)
}

// function reedem use voucher
func RedeemUseVoucherComulative(req models.VoucherComultaiveReq, param models.Params) models.RedeemResponse {
	res := models.RedeemResponse{}

	category := strings.ToLower(req.Category)

	reqRedeem := models.UseRedeemRequest{
		AccountNumber: param.AccountNumber,
		CustID:        req.CustID,
		CustID2:       req.CustID2,
		ProductCode:   param.ProductCode,
	}

	resRedeem := models.UseRedeemResponse{}

	switch category {
	case constants.CategoryPulsa, constants.CategoryPaketData:
		resRedeem = voucher.RedeemPulsaComulative(reqRedeem, req, param)
	case constants.CategoryToken:
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

// function validate response reedem
func ValidateRespRedeem(resRedeem_Msg, campaign, couponId, couponCode, custID string) models.Response {
	res := models.Response{}

	if resRedeem_Msg == "Prefix Failed" {
		logs.Info("[Prefix Failed]")
		logs.Info("[Services-Voucher-UserVoucher]")

		res = utils.GetMessageFailedError(res, 500, errors.New("Transaksi Gagal"))
		return res
	}

	if resRedeem_Msg == "Inquiry Failed" {
		logs.Info("[Inquiry Failed]")
		logs.Info("[Services-Voucher-UserVoucher]")

		res = utils.GetMessageFailedError(res, 500, errors.New("Transaksi Gagal"))
		return res
	}

	if resRedeem_Msg == "Payment Failed" {
		logs.Info("[Payment Failed]")
		logs.Info("[Services-Voucher-UserVoucher]")

		res = utils.GetMessageFailedError(res, 500, errors.New("Transaksi Gagal"))
		return res
	}

	if resRedeem_Msg == "Request in progress" {
		logs.Info("[Prefix Failed]")
		logs.Info("[Services-Voucher-UserVoucher]")

		res = utils.GetMessageFailedError(res, 500, errors.New("Transaksi Pending"))
		return res
	}
	return res
}
