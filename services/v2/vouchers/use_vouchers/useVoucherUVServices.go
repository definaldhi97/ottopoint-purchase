package use_vouchers

import (
	"errors"
	"fmt"
	"ottopoint-purchase/db"
	uv "ottopoint-purchase/hosts/ultra_voucher/host"
	uvmodels "ottopoint-purchase/hosts/ultra_voucher/models"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"

	"github.com/sirupsen/logrus"
)

func UseVoucherUVServices(req models.UseVoucherReq, param models.Params) models.Response {
	// fmt.Println("[ >>>>>>>>>>>>>>>>>> V2 Migrate Use Vouhcer UV Migrate Services <<<<<<<<<<<<<<<< ]")

	var res models.Response

	nameservice := "[PackageUserVouchers]-[UseVoucherUVServices]"
	logReq := fmt.Sprintf("[CouponID : %v]", req.CouponID)

	logrus.Info(nameservice)

	get, errGet := db.GetVoucherUV(param.AccountNumber, req.CouponID)
	if errGet != nil || get.AccountId == "" {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[GetVoucherUV]-[Error : %v]", errGet))
		logrus.Println(logReq)

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

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[UseVoucherUV]-[Response : %v]", useUV))
		logrus.Println(logReq)

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

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[UseVoucherUV]-[Error : %v]", errUV))
		logrus.Println(logReq)

		res = utils.GetMessageResponse(res, 129, false, errors.New("Transaksi tidak Berhasil, Silahkan dicoba kembali."))
		// res.Data = "Transaksi Gagal"
		return res
	}

	return res

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
