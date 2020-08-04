package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	uv "ottopoint-purchase/hosts/ultra_voucher/host"
	uvmodels "ottopoint-purchase/hosts/ultra_voucher/models"
	"ottopoint-purchase/models"
	"ottopoint-purchase/models/dbmodels"
	"ottopoint-purchase/utils"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/opentracing/opentracing-go"
	"github.com/vjeantet/jodaTime"
	"go.uber.org/zap"
)

func (t UseVoucherServices) GetVoucherUV(req models.UseVoucherReq, param models.Params) models.Response {
	var res models.Response

	logs.Info("=== GetVoucherUV ===")
	fmt.Println("=== GetVoucherUV ===")

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[GetVoucherUV-Services]",
		zap.String("AccountNumber : ", param.AccountNumber), zap.String("InstitutionID : ", param.InstitutionID),
		zap.String("category : ", param.Category), zap.String("campaignId : ", req.CampaignID),
		zap.String("cust_id : ", req.CustID), zap.String("cust_id2 : ", req.CustID2),
		zap.String("product_code : ", param.ProductCode))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[GetVoucherUV]")
	defer span.Finish()

	// // Use Voucher to Openloyalty
	// _, err2 := opl.CouponVoucherCustomer(req.CampaignID, param.CouponID, param.ProductCode, param.CustID, 1)
	// if err2 != nil {
	// 	res = utils.GetMessageResponse(res, 400, false, errors.New("Gagal Use Voucher, Harap coba lagi"))
	// 	return res
	// }

	get, errGet := db.GetVoucherUV(param.AccountNumber, param.CouponID)
	if errGet != nil || get.AccountId == "" {
		logs.Info("Internal Server Error : ", errGet)
		logs.Info("[GetVoucherUV-Servcies]-[GetVoucherUV]")
		logs.Info("[Failed get data from DB]")

		// sugarLogger.Info("Internal Server Error : ", errGet)
		sugarLogger.Info("[GetVoucherUV-Servcies]-[GetVoucherUV]")
		sugarLogger.Info("[Failed get data from DB]")

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
		// go SaveTransactionUV(param, useUV, reqUV, req, "Inquiry", "01", useUV.ResponseCode)

		res = utils.GetMessageResponse(res, 147, false, errors.New("Voucher Tidak Ditemukan"))
		// res.Data = "Transaksi Gagal"

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

func SaveTransactionUV(param models.Params, res interface{}, reqdata interface{}, reqOP interface{}, trasnType, status string) {

	fmt.Println(fmt.Sprintf("[Start-SaveDB]-[UltraVoucher]-[%v]", trasnType))

	var saveStatus string
	switch status {
	case "00":
		saveStatus = constants.Success
	case "09":
		saveStatus = constants.Pending
	case "01":
		saveStatus = constants.Failed
	}

	reqUV, _ := json.Marshal(&reqdata)   // Req UV
	responseUV, _ := json.Marshal(&res)  // Response UV
	reqdataOP, _ := json.Marshal(&reqOP) // Req Service

	timeRedeem := jodaTime.Format("dd-MM-YYYY HH:mm:ss", time.Now())

	save := dbmodels.TransaksiRedeem{
		AccountNumber: param.AccountNumber,
		Voucher:       param.NamaVoucher,
		MerchantID:    param.MerchantID,
		// CustID:          param.CustID,
		RRN:             param.RRN,
		ProductCode:     param.ProductCode,
		Amount:          int64(param.Amount),
		TransType:       trasnType,
		IsUsed:          false,
		ProductType:     param.ProductType,
		Status:          saveStatus,
		ExpDate:         param.ExpDate,
		Institution:     param.InstitutionID,
		CummulativeRef:  param.CumReffnum,
		DateTime:        utils.GetTimeFormatYYMMDDHHMMSS(),
		ResponderData:   status,
		Point:           param.Point,
		ResponderRc:     param.DataSupplier.Rc,
		ResponderRd:     param.DataSupplier.Rd,
		RequestorData:   string(reqUV),
		ResponderData2:  string(responseUV),
		RequestorOPData: string(reqdataOP),
		SupplierID:      param.SupplierID,
		CouponId:        param.CouponID,
		CampaignId:      param.CampaignID,
		AccountId:       param.AccountId,
		RedeemAt:        timeRedeem,
	}

	err := db.DbCon.Create(&save).Error
	if err != nil {
		logs.Info(fmt.Sprintf("[Error : %v]", err))
		logs.Info("[Failed Save to DB]")

		name := jodaTime.Format("dd-MM-YYYY", time.Now()) + ".csv"
		go utils.CreateCSVFile(save, name)

		// return err

	}
}
