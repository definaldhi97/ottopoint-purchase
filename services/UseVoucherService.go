package services

import (
	"errors"
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	opl "ottopoint-purchase/hosts/opl/host"
	redismodels "ottopoint-purchase/hosts/redis_token/models"
	"ottopoint-purchase/models"
	"ottopoint-purchase/redis"
	redeem "ottopoint-purchase/services/voucher"
	"ottopoint-purchase/utils"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type UseVoucherServices struct {
	General models.GeneralModel
}

func (t UseVoucherServices) UseVoucher(req models.UseVoucherReq, dataToken redismodels.TokenResp) models.Response {
	var res models.Response

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[UseVoucher-Services]",
		zap.String("category", req.Category), zap.String("campaignId", req.CampaignID),
		zap.String("cust_id", req.CustID), zap.String("cust_id2", req.CustID2),
		zap.String("product_code", req.ProductCode), zap.String("date", req.Date))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[RedeemVoucher]")
	defer span.Finish()

	// get CustID
	dataUser, errUser := db.CheckUser(dataToken.Data)
	if errUser != nil {
		res = utils.GetMessageResponse(res, 422, false, errors.New("User belum Eligible"))
		return res
	}

	data, err := opl.HistoryVoucherCustomer(dataToken.Data, "")
	if err != nil {
		res = utils.GetMessageResponse(res, 422, false, errors.New("Gagal Get History Voucher Customer"))
		return res
	}

	// var campaign, couponId, couponCode string
	resp := []models.CampaignsDetail{}
	for _, val := range data.Campaigns {
		if val.CampaignID == req.CampaignID && val.CanBeUsed == true {
			a := models.CampaignsDetail{
				Name:       val.Campaign.Name,
				CampaignID: val.CampaignID,
				ActiveTo:   val.ActiveTo,
				Coupon: models.CouponDetail{
					Code: val.Coupon.Code,
					ID:   val.Coupon.ID,
				},
			}

			resp = append(resp, a)
		}
	}

	var campaign, couponId, couponCode, nama, expDate string
	for _, valco := range resp {
		nama = valco.Name
		campaign = valco.CampaignID
		couponId = valco.Coupon.ID
		couponCode = valco.Coupon.Code
		expDate = valco.ActiveTo
	}

	// Use Voucher to Openloyalty
	_, err2 := opl.CouponVoucherCustomer(campaign, couponId, couponCode, dataUser.CustID, 1)
	if err2 != nil {
		res = utils.GetMessageResponse(res, 400, false, errors.New("Gagal Redeem Voucher, Harap coba lagi"))
		return res
	}

	// save to redis usedAt
	req.Date = (time.Now().Local().Add(time.Hour * time.Duration(7))).Format("2006-01-02T15:04:05-0700")
	keyTimeVoucher := fmt.Sprintf("usedVoucherAt-%s-%s", couponId, dataToken.Data)
	go redis.SaveRedis(keyTimeVoucher, req.Date)

	cekData, errDB := db.CheckUser(dataToken.Data)
	if errDB != nil || cekData.Phone == "" {
		res = utils.GetMessageResponse(res, 422, false, errors.New("Internal Server Error"))
		return res
	}

	//get config
	memberid, errmember := db.GetConfig()
	if errmember != nil {
		fmt.Println("[EEROR-DATABASE]")
		res = utils.GetMessageResponse(res, 422, false, errors.New("Internal Server Error"))
		return res
	}

	category := strings.ToLower(req.Category)

	logs.Info("===== nama : %v =====", resp[0].Name)
	logs.Info("===== Category : %v =====", category)
	logs.Info("===== couponId : %v =====", couponId)
	logs.Info("===== couponCode : %v =====", couponCode)
	logs.Info("===== expDate : %v =====", expDate)

	resRedeem := models.UseRedeemResponse{}

	reqRedeem := models.UseRedeemRequest{
		AccountNumber: dataToken.Data,
		CustID:        req.CustID,
		CustID2:       req.CustID2,
		ProductCode:   req.ProductCode,
	}
	switch category {
	case constants.CategoryPulsa, constants.CategoryPaketData:
		resRedeem = redeem.RedeemPulsa(reqRedeem, dataToken, memberid.MemberID, nama, expDate, category)
	case constants.CategoryToken:
		resRedeem = redeem.RedeemPLN(reqRedeem, dataToken, memberid.MemberID, nama, expDate, category)
	case constants.CategoryMobileLegend, constants.CategoryFreeFire:
		resRedeem = redeem.RedeemGame(reqRedeem, dataToken, memberid.MemberID, nama, expDate, category)
	}

	if resRedeem.Msg == "Prefix Failed" {
		logs.Info("[Prefix Failed]")
		logs.Info("[Services-Voucher-UserVoucher]")

		res = utils.GetMessageResponse(res, 500, false, errors.New("Transaksi Gagal"))
		return res
	}

	if resRedeem.Msg == "Inquiry Failed" {
		logs.Info("[Inquiry Failed]")
		logs.Info("[Services-Voucher-UserVoucher]")

		logs.Info("[Reversal Voucher")
		_, erv := opl.CouponVoucherCustomer(campaign, couponId, couponCode, dataUser.CustID, 0)
		if erv != nil {
			res = utils.GetMessageResponse(res, 500, false, errors.New("Gagal Reversal Voucher"))
			return res
		}

		res = utils.GetMessageResponse(res, 500, false, errors.New("Transaksi Gagal"))
		return res
	}

	if resRedeem.Msg == "Payment Failed" {
		logs.Info("[Payment Failed]")
		logs.Info("[Services-Voucher-UserVoucher]")

		logs.Info("[Reversal Voucher")
		_, erv := opl.CouponVoucherCustomer(campaign, couponId, couponCode, dataUser.CustID, 0)
		if erv != nil {
			res = utils.GetMessageResponse(res, 500, false, errors.New("Gagal Reversal Voucher"))
			return res
		}

		res = utils.GetMessageResponse(res, 500, false, errors.New("Transaksi Gagal"))
		return res
	}

	if resRedeem.Msg == "Request in progress" {
		logs.Info("[Prefix Failed]")
		logs.Info("[Services-Voucher-UserVoucher]")

		res = utils.GetMessageResponse(res, 500, false, errors.New("Transaksi Pending"))
		return res
	}

	if resRedeem.Msg == "SUCCESS" {
		if resRedeem.Category == "PLN" {
			res = models.Response{
				Data: models.ResponseUseVoucherPLN{
					Voucher:     nama,
					CustID:      resRedeem.CustID,
					CustID2:     resRedeem.CustID2,
					ProductCode: resRedeem.ProductCode,
					Amount:      resRedeem.Amount,
					Token:       resRedeem.Data.Tokenno,
				},
				Meta: utils.ResponseMetaOK(),
			}
			return res
		}

		res = models.Response{
			Data: models.ResponseUseVoucher{
				Voucher:     nama,
				CustID:      resRedeem.CustID,
				CustID2:     resRedeem.CustID2,
				ProductCode: resRedeem.ProductCode,
				Amount:      resRedeem.Amount,
			},
			Meta: utils.ResponseMetaOK(),
		}
		return res
	}

	return res
}
