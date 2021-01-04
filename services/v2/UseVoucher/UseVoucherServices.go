package UseVoucher

import (
	"errors"
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	uv "ottopoint-purchase/hosts/ultra_voucher/host"
	uvmodels "ottopoint-purchase/hosts/ultra_voucher/models"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

type V2_UseVoucherServices struct {
	General models.GeneralModel
}

func (service V2_UseVoucherServices) UseVoucherUV(req models.UseVoucherReq, param models.Params) models.Response {
	fmt.Println("[ >>>>>>>>>>>>>>>>>> V2 Migrate Use Vouhcer UV Migrate Services <<<<<<<<<<<<<<<< ]")

	var res models.Response

	sugarLogger := service.General.OttoZaplog
	sugarLogger.Info("[GetVoucherUV-Services]",
		zap.String("AccountNumber : ", param.AccountNumber), zap.String("InstitutionID : ", param.InstitutionID),
		zap.String("category : ", param.Category), zap.String("campaignId : ", req.CampaignID),
		zap.String("cust_id : ", req.CustID), zap.String("cust_id2 : ", req.CustID2),
		zap.String("product_code : ", param.ProductCode))

	span, _ := opentracing.StartSpanFromContext(service.General.Context, "[GetVoucherUV]")
	defer span.Finish()

	get, errGet := db.GetVoucherUV(param.AccountNumber, req.CouponID)
	if errGet != nil || get.AccountId == "" {
		logrus.Error("Failed Get Voucher UV : ", errGet)
		res = utils.GetMessageResponse(res, 422, false, errors.New("Voucher Tidak Ditemukan"))
		return res
	}

	comulative_ref := utils.GenTransactionId()
	param.Reffnum = comulative_ref
	param.Amount = int64(param.Point)

	reqUV := uvmodels.UseVoucherUVReq{
		Account:     get.AccountId,
		VoucherCode: get.VoucherCode,
	}

	// get to UV
	useUV, errUV := uv.UseVoucherUV(reqUV)

	if useUV.ResponseCode == "10" {

		fmt.Println(fmt.Sprintf("[Response UV : %v]", useUV.ResponseCode))

		fmt.Println(">>> Voucher Tidak Ditemukan <<<")

		res = utils.GetMessageResponse(res, 147, false, errors.New("Voucher Tidak Ditemukan"))

		return res
	}

	if useUV.ResponseCode == "14" || useUV.ResponseCode == "00" {

		fmt.Println(">>> Success <<<")

		fmt.Println(fmt.Sprintf("[Response UV : %v]", useUV.ResponseCode))
		res = models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: models.GetVoucherUVResp{
				Voucher:     param.NamaVoucher,
				VoucherCode: get.VoucherCode,
				Link:        useUV.Data.Link,
			},
		}
		return res
	}

	if errUV != nil || useUV.ResponseCode == "" || useUV.ResponseCode != "00" {

		fmt.Println(">>> Time Out / Gagal <<<")
		fmt.Println(fmt.Sprintf("[Response UV : %v]", useUV.ResponseCode))
		logs.Info("Internal Server Error : ", errUV)
		logs.Info("[GetVoucherUV-Servcies]-[UseVoucherUV]")
		logs.Info("[Failed Use Voucher UV]-[Gagal Use Voucher UV]")

		// sugarLogger.Info("Internal Server Error : ", errUV)
		sugarLogger.Info("[GetVoucherUV-Servcies]-[UseVoucherUV]")
		sugarLogger.Info("[Failed Use Voucher UV]-[Gagal Use Voucher UV]")

		res = utils.GetMessageResponse(res, 129, false, errors.New("Transaksi tidak Berhasil, Silahkan dicoba kembali."))
		// res.Data = "Transaksi Gagal"
		return res
	}

	return res

}

func (service V2_UseVoucherServices) UseVoucherAggregator(req models.UseVoucherReq, param models.Params) models.Response {
	fmt.Println("[ >>>>>>>>>>>>>>>>>> V2 Migrate Use Vouhcer Aggregator Migrate Services <<<<<<<<<<<<<<<< ]")

	var res models.Response

	sugarLogger := service.General.OttoZaplog
	sugarLogger.Info("[UseVoucherAggregator-Services]",
		zap.String("AccountNumber : ", param.AccountNumber), zap.String("InstitutionID : ", param.InstitutionID),
		zap.String("category : ", param.Category), zap.String("campaignId : ", req.CampaignID),
		zap.String("cust_id : ", req.CustID), zap.String("cust_id2 : ", req.CustID2),
		zap.String("product_code : ", param.ProductCode))

	span, _ := opentracing.StartSpanFromContext(service.General.Context, "[GetVoucherUV]")
	defer span.Finish()

	get, err := db.GetVoucherAg(param.AccountId, req.CouponID)
	if err != nil {
		logs.Info("Internal Server Error : ", err)
		logs.Info("[UseVoucherAggregator]-[GetVoucherUV]")
		logs.Info("[Failed get data from DB]")

		sugarLogger.Info("[UseVoucherAggregator]-[GetVoucherUV]")
		sugarLogger.Info("[Failed get data from DB]")

		res = utils.GetMessageResponse(res, 422, false, errors.New("Voucher Tidak Ditemukan"))
		return res
	}

	spend, err := db.GetVoucherSpending(get.AccountId, get.CouponID)
	if err != nil {
		logs.Info("Internal Server Error : ", err)
		logs.Info("[UseVoucherAggregator]-[GetVoucherUV]")
		logs.Info("[Failed get data from DB]")

		sugarLogger.Info("[UseVoucherAggregator]-[GetVoucherUV]")
		sugarLogger.Info("[Failed get data from DB]")

		res = utils.GetMessageResponse(res, 422, false, errors.New("Terjadi Kesalahan"))
		return res
	}

	// // Update Status Voucher
	// // timeUse := jodaTime.Format("dd-MM-YYYY HH:mm:ss", time.Now())
	// timeUse := time.Now()
	// go db.UpdateVoucher(timeUse, spend.CouponId)

	// codeVoucher := decryptVoucherCode(spend.VoucherCode, spend.CouponId)

	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: models.GetVoucherAgResp{
			Voucher:     param.NamaVoucher,
			VoucherCode: spend.VoucherCode,
			Link:        spend.VoucherLink,
		},
	}
	return res

}

func (service V2_UseVoucherServices) UseVoucherVidio(couponId string) models.Response {
	fmt.Println("[ >>>>>>>>>>>>>>>>>> V2 Migrate Use Vouhcer Vidio Migrate Services <<<<<<<<<<<<<<<< ]")

	resp := models.Response{}

	sugarLogger := service.General.OttoZaplog
	sugarLogger.Info("[ViewVoucher-Services]",
		zap.String("Coupon Id : ", couponId))

	span, _ := opentracing.StartSpanFromContext(service.General.Context, "[ViewVoucher]")
	defer span.Finish()

	// get voucher
	getVouc, errGetVouc := db.GetUseVoucher(couponId)
	if errGetVouc != nil || getVouc.AccountNumber == "" {

		logrus.Error("Failed GetVoucher : ", errGetVouc)
		resp = utils.GetMessageFailedErrorNew(resp, constants.RC_VOUCHER_NOTFOUND, constants.RD_VOUCHER_NOTFOUND)
		return resp
	}

	// update transaction redeem into use
	// timeUse := jodaTime.Format("dd-MM-YYYY HH:mm:ss", time.Now())
	timeUse := time.Now()
	_, errUpdate := db.UpdateVoucher(timeUse, getVouc.CouponId)
	if errUpdate != nil {
		logrus.Error("Failed Update Status Voucher : ", errUpdate)
		logrus.Info("[Gagal Update Voucher]")
		logrus.Info("[UseVoucherVidio]-[Package-Services]")
	}

	respVouch := models.RespUseVoucher{}
	respVouch.Code = getVouc.ProductCode
	respVouch.CouponID = getVouc.CouponId
	respVouch.Used = true
	respVouch.CampaignID = getVouc.MRewardID
	respVouch.CustomerID = getVouc.AccountId

	resp = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: respVouch,
	}

	return resp
}

func decryptVoucherCode(voucherCode, couponID string) string {

	var codeVoucher string
	if voucherCode == "" {
		return voucherCode
	}

	a := []rune(couponID)
	key32 := string(a[0:32])
	secretKey := []byte(key32)
	codeByte := []byte(voucherCode)
	chiperText, _ := utils.DecryptAES(codeByte, secretKey)
	codeVoucher = string(chiperText)

	return codeVoucher

}
