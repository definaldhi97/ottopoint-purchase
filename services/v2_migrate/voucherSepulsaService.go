package v2_migrate

import (
	"encoding/json"
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	kafka "ottopoint-purchase/hosts/publisher/host"
	sepulsa "ottopoint-purchase/hosts/sepulsa/host"
	sepulsaModels "ottopoint-purchase/hosts/sepulsa/models"
	"ottopoint-purchase/models"
	"ottopoint-purchase/models/dbmodels"
	"ottopoint-purchase/utils"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/vjeantet/jodaTime"
	"go.uber.org/zap"
)

type VoucherSepulsaMigrateService struct {
	General models.GeneralModel
}

func (t VoucherSepulsaMigrateService) VoucherSepulsa(req models.VoucherComultaiveReq, param models.Params, header models.RequestHeader) models.Response {
	fmt.Println("[ >>>>>>>>>>>>>>>>>> Voucher Sepulsa Service <<<<<<<<<<<<<<<< ]")

	var res models.Response

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[SepulsaServices]",
		zap.String("NameVoucher : ", param.NamaVoucher), zap.Int("Jumlah : ", req.Jumlah),
		zap.String("CampaignID : ", req.CampaignID), zap.String("CampaignID : ", req.CampaignID),
		zap.String("CustID2 : ", req.CustID2), zap.String("ProductCode : ", param.ProductCode),
		zap.String("AccountNumber : ", param.AccountNumber), zap.String("InstitutionID : ", param.InstitutionID))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[SepulsaServices]")
	defer span.Finish()

	param.CumReffnum = utils.GenTransactionId()

	// total := strconv.Itoa(req.Jumlah)
	param.Amount = int64(param.Point)

	// spending point and spending usage_limit voucher
	textCommentSpending := param.CumReffnum + "#" + param.NamaVoucher
	param.Comment = textCommentSpending
	RedeemVouchSP, errRedeemVouchSP := Redeem_PointandVoucher(req.Jumlah, param, param.CumReffnum, header)

	logrus.Info("Response Spending point / Deduct point")
	logrus.Info(RedeemVouchSP)

	if RedeemVouchSP.Rc == "10" || RedeemVouchSP.Rd == "Insufficient Point" {

		logrus.Info("[SepulsaVoucherService]-[RedeemVoucher]")
		logrus.Info("[Not enough points]-[Gagal Redeem Voucher]")
		logrus.Info("[Rc] : ", RedeemVouchSP.Rc)
		logrus.Info("[Rd] : ", RedeemVouchSP.Rd)

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

	if RedeemVouchSP.Rc == "208" || RedeemVouchSP.Rd == "Voucher not available" {

		logrus.Info("[SepulsaVoucherService]-[RedeemVoucher]")
		logrus.Info("[Voucher not available]-[Gagal Redeem Voucher]")
		logrus.Info("[Rc] : ", RedeemVouchSP.Rc)
		logrus.Info("[Rd] : ", RedeemVouchSP.Rd)

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

	var c string
	for _, vall := range RedeemVouchSP.CouponseVouch {
		c = vall.CouponsCode
	}

	if errRedeemVouchSP != nil || RedeemVouchSP.Rc != "00" || c == "" {

		logrus.Info("[SepulsaVoucherService]-[RedeemVoucher]")
		logrus.Info("[Rc] : ", RedeemVouchSP.Rc)
		logrus.Info("[Rd] : ", RedeemVouchSP.Rd)

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
		couponID := RedeemVouchSP.CouponseVouch[t].CouponsID
		couponCode := RedeemVouchSP.CouponseVouch[t].CouponsCode
		param.CouponID = couponID

		productID, _ := strconv.Atoi(param.ProductCode)
		reqOrder := sepulsaModels.EwalletInsertTrxReq{
			CustomerNumber: req.CustID,
			OrderID:        param.TrxID,
			ProductID:      productID,
		}

		// Create Transaction Ewallet
		sepulsaRes, errTransaction := sepulsa.EwalletInsertTransaction(reqOrder)

		if errTransaction != nil {

			logrus.Info("[SepulsaService]-[InsertTransaction]")
			logrus.Error("ResponseDesc : ", errTransaction.Error())

			resultReversal := Adding_PointVoucher(param, param.Point, 1, header)
			fmt.Println(resultReversal)

			fmt.Println("[ >>>>>>>>>>>>>>>>>>>>>>> Send Publisher <<<<<<<<<<<<<<<<<<<< ]")
			pubreq := models.NotifPubreq{
				Type:           constants.CODE_REVERSAL_POINT,
				NotificationTo: param.AccountNumber,
				Institution:    param.InstitutionID,
				ReferenceId:    param.RRN,
				TransactionId:  param.Reffnum,
				Data: models.DataValue{
					RewardValue: "point",
					Value:       strconv.Itoa(param.Point),
				},
			}

			bytePub, _ := json.Marshal(pubreq)

			kafkaReq := kafka.PublishReq{
				Topic: utils.TopicsNotif,
				Value: bytePub,
			}

			kafkaRes, err := kafka.SendPublishKafka(kafkaReq)
			if err != nil {
				logrus.Error("Gagal Send Publisher : ", err)
			}

			logrus.Info("[ Response Publisher ] : ", kafkaRes)

			// Save Error Transaction
			go SaveTransactionSepulsa(param, errTransaction.Error(), reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "01")

			res = models.Response{
				Meta: utils.ResponseMetaOK(),
				Data: models.SepulsaRes{
					Code:    "01",
					Msg:     "Gagal! Maaf transaksi Anda tidak dapat dilakukan saat ini. Silahkan dicoba lagi atau hubungi tim kami untuk informasi selengkapnya.",
					Success: 0,
					Failed:  req.Jumlah,
					Pending: 0,
				},
			}
			return res
		}

		param.DataSupplier.Rd = sepulsaRes.Status
		param.DataSupplier.Rc = sepulsaRes.ResponseCode
		param.RRN = sepulsaRes.TransactionID

		id := utils.GenerateTokenUUID()
		go SaveDBSepulsa(id, param.InstitutionID, couponID, couponCode, param.AccountNumber, param.AccountId, req.CampaignID)
		go SaveTransactionSepulsa(param, sepulsaRes, reqOrder, req, constants.CODE_TRANSTYPE_REDEMPTION, "09")

	}

	res = models.Response{
		Meta: utils.ResponseMetaOK(),
		Data: models.SepulsaRes{
			Code:    "00",
			Msg:     fmt.Sprintf("Selamat Penukaran %s Kamu Berhasil, Silahkan Cek Saldo Kamu!", param.NamaVoucher),
			Success: req.Jumlah,
			Failed:  0,
			Pending: 0,
		},
	}

	return res

}

func SaveTransactionSepulsa(param models.Params, res interface{}, reqdata interface{}, reqOP models.VoucherComultaiveReq, transType, status string) {

	fmt.Println(fmt.Sprintf("[Start-SaveDB]-[Sepulsa]-[%v]", transType))

	var saveStatus string
	switch status {
	case "00":
		saveStatus = constants.Success
	case "09":
		saveStatus = constants.Pending
	case "01":
		saveStatus = constants.Failed
	}

	reqSepulsa, _ := json.Marshal(&reqdata)
	responseSepulsa, _ := json.Marshal(&res)
	reqdataOP, _ := json.Marshal(&reqOP)

	timeRedeem := jodaTime.Format("dd-MM-YYYY HH:mm:ss", time.Now())

	save := dbmodels.TSpending{
		ID:                utils.GenerateTokenUUID(),
		AccountNumber:     param.AccountNumber,
		Voucher:           param.NamaVoucher,
		MerchantID:        param.MerchantID,
		CustID:            reqOP.CustID,
		RRN:               param.RRN,
		TransactionId:     param.TrxID,
		ProductCode:       param.ProductCode,
		Amount:            int64(param.Amount),
		TransType:         transType,
		IsUsed:            true,
		UsedAt:            timeRedeem,
		ProductType:       param.ProductType,
		Status:            saveStatus,
		ExpDate:           param.ExpDate,
		Institution:       param.InstitutionID,
		CummulativeRef:    param.CumReffnum,
		DateTime:          utils.GetTimeFormatYYMMDDHHMMSS(),
		Point:             param.Point,
		ResponderRc:       param.DataSupplier.Rc,
		ResponderRd:       param.DataSupplier.Rd,
		RequestorData:     string(reqSepulsa),
		ResponderData:     string(responseSepulsa),
		RequestorOPData:   string(reqdataOP),
		SupplierID:        param.SupplierID,
		CouponId:          param.CouponID,
		CampaignId:        param.CampaignID,
		AccountId:         param.AccountId,
		RedeemAt:          timeRedeem,
		Comment:           param.Comment,
		RewardID:          param.RewardID,
		ProductCategoryID: param.CategoryID,
		MProductID:        param.ProductID,
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

func SaveDBSepulsa(id, institution, coupon, vouchercode, phone, custIdOPL, campaignID string) {
	fmt.Println("[SaveDB]-[SepulsaVoucherService]")
	save := dbmodels.UserMyVocuher{
		ID:            id,
		InstitutionID: institution,
		CouponID:      coupon,
		VoucherCode:   vouchercode,
		Phone:         phone,
		AccountId:     custIdOPL,
		CampaignID:    campaignID,
	}

	err := db.DbCon.Create(&save).Error
	if err != nil {
		fmt.Println("[Failed Save to DB]", err)
		fmt.Println("[Package-Service]-[SepulsaService]")
	}
}
