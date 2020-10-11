package services

import (
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/hosts/opl/host"
	sepulsa "ottopoint-purchase/hosts/sepulsa/host"
	sepulsamodels "ottopoint-purchase/hosts/sepulsa/models"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"
	"strconv"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type UseSepulsaService struct {
	General models.GeneralModel
}

func (t UseSepulsaService) SepulsaServices(req models.VoucherComultaiveReq, param models.Params) models.Response {
	var res models.Response

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[SepulsaServices]",
		zap.String("NameVoucher : ", param.NamaVoucher), zap.Int("Jumlah : ", req.Jumlah),
		zap.String("CampaignID : ", req.CampaignID), zap.String("CampaignID : ", req.CampaignID),
		zap.String("CustID2 : ", req.CustID2), zap.String("ProductCode : ", param.ProductCode),
		zap.String("AccountNumber : ", param.AccountNumber), zap.String("InstitutionID : ", param.InstitutionID))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[SepulsaServices]")
	defer span.Finish()

	total := strconv.Itoa(req.Jumlah)
	param.CumReffnum = utils.GenTransactionId()
	param.Amount = int64(param.Point)

	redeem, errredeem := host.RedeemVoucherCumulative(req.CampaignID, param.AccountId, total, "0")
	if redeem.Message == "Invalid JWT Token" {
		fmt.Println("Error : ", errredeem)
		fmt.Println("[SepulsaVoucherService]-[RedeemVoucher]")
		fmt.Println("[Internal Server Error]-[Gagal Redeem Voucher]")

		sugarLogger.Info("[SepulsaVoucherService]-[RedeemVoucher]")
		sugarLogger.Info("[Internal Server Error]-[Gagal Redeem Voucher]")

		res := models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: models.SepulsaRes{
				Code:    "65",
				Msg:     "Token or session Expired Please Login Again",
				Success: 0,
				Failed:  req.Jumlah,
				Pending: 0,
			},
		}

		return res
	}

	if redeem.Error == "Not enough points" {
		fmt.Println("Error : ", errredeem)
		fmt.Println("[SepulsaVoucherService]-[RedeemVoucher]")
		fmt.Println("[Not enough points]-[Gagal Redeem Voucher]")

		sugarLogger.Info("[SepulsaVoucherService]-[RedeemVoucher]")
		sugarLogger.Info("[Not enough points]-[Gagal Redeem Voucher]")

		res := models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: models.SepulsaRes{
				Code:    "27",
				Msg:     "Point Tidak Mencukupi",
				Success: 0,
				Failed:  req.Jumlah,
				Pending: 0,
			},
		}

		return res
	}

	if redeem.Error == "Limit exceeded" {
		fmt.Println("Error : ", errredeem)
		fmt.Println("[SepulsaVoucherService]-[RedeemVoucher]")
		fmt.Println("[Limit exceed]-[Gagal Redeem Voucher]")

		sugarLogger.Info("[SepulsaVoucherService]-[RedeemVoucher]")
		sugarLogger.Info("[Limit exceed]-[Gagal Redeem Voucher]")

		res := models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: models.SepulsaRes{
				Code:    "65",
				Msg:     "Payment count limit exceed",
				Success: 0,
				Failed:  req.Jumlah,
				Pending: 0,
			},
		}

		return res
	}

	var coupon string
	for _, val := range redeem.Coupons {
		coupon = val.Code
	}

	if errredeem != nil || redeem.Error != "" || coupon == "" {
		fmt.Println("Error : ", errredeem)
		fmt.Println("[SepulsaVoucherService]-[RedeemVoucher]")
		fmt.Println("[Failed Redeem Voucher]-[Gagal Redeem Voucher]")

		sugarLogger.Info("[SepulsaVoucherService]-[RedeemVoucher]")
		sugarLogger.Info("[Failed Redeem Voucher]-[Gagal Redeem Voucher]")

		res = models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: models.SepulsaRes{
				Code:    "01",
				Msg:     "Gagal! Maaf transaksi Anda tidak dapat dilakukan saat ini. Silahkan dicoba lagi atau hubungi tim kami untuk informasi selengkapnya",
				Success: 0,
				Failed:  req.Jumlah,
				Pending: 0,
			},
		}

		return res
	}

	for i := req.Jumlah; i > 0; i-- {

		param.TrxID = utils.GenTransactionId()

		t := i - 1

		coupon := redeem.Coupons[t].Id
		param.CouponID = coupon

		productID, _ := strconv.Atoi(param.CouponCode)
		reqOrder := sepulsamodels.EwalletInsertTrxReq{
			CustomerNumber: param.AccountNumber,
			OrderID:        param.TrxID,
			ProductID:      productID,
		}

		// Create Transaction Ewallet
		sepulsaRes, err := sepulsa.EwalletInsertTransaction(reqOrder)
		if err != nil {
			fmt.Println("[SepulsaService]-[InsertTransaction]")
			res = models.Response{
				Meta: utils.ResponseMetaOK(),
				Data: models.SepulsaRes{
					Code:    "",
					Msg:     err.Error(),
					Success: 0,
					Failed:  req.Jumlah,
					Pending: 0,
				},
			}
			return res
		}

		param.DataSupplier.Rd = sepulsaRes.Status
		param.DataSupplier.Rc = "201"
		param.RRN = param.TrxID

		go SaveTransactionUV(param, sepulsaRes, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "09")

	}

	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: models.SepulsaRes{
			Code:    "00",
			Msg:     "Success",
			Success: 0,
			Failed:  0,
			Pending: req.Jumlah,
		},
	}

	return res

}
