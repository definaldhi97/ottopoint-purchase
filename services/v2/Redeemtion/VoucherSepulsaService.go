package redeemtion

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
	"ottopoint-purchase/services/v2/Trx"
	"ottopoint-purchase/utils"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/sirupsen/logrus"
	"github.com/vjeantet/jodaTime"
)

// func (t V2_VoucherSepulsaService) VoucherSepulsa(req models.VoucherComultaiveReq, param models.Params) models.Response {
func RedeemtionSepulsaServices(req models.VoucherComultaiveReq, param models.Params) models.Response {
	// fmt.Println("[ >>>>>>>>>>>>>>>>>> V2 Migrate Voucher Sepulsa Service <<<<<<<<<<<<<<<< ]")

	nameservice := "[PackageRedeemtion]-[RedeemtionSepulsaServices]"
	logReq := fmt.Sprintf("[AccountNumber : %v, RewardID : %v]", param.AccountNumber, param.RewardID)

	logrus.Info(nameservice)

	var res models.Response

	param.CumReffnum = utils.GenTransactionId()

	// total := strconv.Itoa(req.Jumlah)
	param.Amount = int64(param.Point)

	// spending point and spending usage_limit voucher
	textCommentSpending := param.CumReffnum + "#" + param.NamaVoucher
	param.Comment = textCommentSpending
	RedeemVouchSP, errRedeemVouchSP := Trx.V2_Redeem_PointandVoucher(req.Jumlah, param)
	param.PointTransferID = RedeemVouchSP.PointTransferID

	if RedeemVouchSP.Rd == "Invalid JWT Token" {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[V2_Redeem_PointandVoucher]-[Response : %v]", RedeemVouchSP))
		logrus.Info("[ ResponseCode ] : ", RedeemVouchSP.Rc)
		logrus.Info("[ ResponseDesc ] : ", RedeemVouchSP.Rd)
		logrus.Println(logReq)

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

	if RedeemVouchSP.Rd == "not enough points" {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[V2_Redeem_PointandVoucher]-[Response : %v]", RedeemVouchSP))
		logrus.Info("[ ResponseCode ] : ", RedeemVouchSP.Rc)
		logrus.Info("[ ResponseDesc ] : ", RedeemVouchSP.Rd)
		logrus.Println(logReq)

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

	if RedeemVouchSP.Rc == "208" {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[V2_Redeem_PointandVoucher]-[Response : %v]", RedeemVouchSP))
		logrus.Info("[ ResponseCode ] : ", RedeemVouchSP.Rc)
		logrus.Info("[ ResponseDesc ] : ", RedeemVouchSP.Rd)
		logrus.Println(logReq)

		res := models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: models.SepulsaRes{
				Code:    "65",
				Msg:     "Voucher not available",
				Success: 0,
				Failed:  req.Jumlah,
				Pending: 0,
			},
		}

		return res
	}

	if RedeemVouchSP.Rc == "209" {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[V2_Redeem_PointandVoucher]-[Response : %v]", RedeemVouchSP))
		logrus.Info("[ ResponseCode ] : ", RedeemVouchSP.Rc)
		logrus.Info("[ ResponseDesc ] : ", RedeemVouchSP.Rd)
		logrus.Println(logReq)

		res = models.Response{
			Meta: utils.ResponseMetaOK(),
			Data: models.SepulsaRes{
				Code:    "66",
				Msg:     "Payment count limit exceeded",
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
	fmt.Println("Value CouponCode : ", c)

	if errRedeemVouchSP != nil || RedeemVouchSP.Rc != "00" {

		logrus.Error(nameservice)
		logrus.Error(fmt.Sprintf("[V2_Redeem_PointandVoucher]-[Error : %v]", errRedeemVouchSP))
		logrus.Info("[ ResponseCode ] : ", RedeemVouchSP.Rc)
		logrus.Info("[ ResponseDesc ] : ", RedeemVouchSP.Rd)
		logrus.Println(logReq)

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

			logrus.Error(nameservice)
			logrus.Error(fmt.Sprintf("[EwalletInsertTransaction]-[Error : %v]", errTransaction))
			logrus.Println(logReq)

			resultReversal := Trx.V2_Adding_PointVoucher(param, param.Point, 1)
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

			_, errKafka := kafka.SendPublishKafka(kafkaReq)
			if errKafka != nil {

				logrus.Error(nameservice)
				logrus.Error(fmt.Sprintf("[SendPublishKafka]-[Error : %v]", errKafka))
				logrus.Println(logReq)

			}

			// logrus.Info("[ Response Publisher ] : ", kafkaRes)

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
		go SaveTSchedulerRetry(param.RRN, constants.CodeSchedulerSepulsa)
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

	var ExpireDate time.Time
	var redeemDate time.Time

	var saveStatus string
	switch status {
	case "00":
		saveStatus = constants.Success
	case "09":
		saveStatus = constants.Pending
	case "01":
		saveStatus = constants.Failed
	}

	if transType == constants.CODE_TRANSTYPE_REDEMPTION {
		ExpireDate = utils.ExpireDateVoucherAGt(constants.EXPDATE_VOUCHER)
		redeemDate = time.Now()
	}

	reqSepulsa, _ := json.Marshal(&reqdata)
	responseSepulsa, _ := json.Marshal(&res)
	reqdataOP, _ := json.Marshal(&reqOP)

	// timeRedeem := jodaTime.Format("dd-MM-YYYY HH:mm:ss", time.Now())

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
		UsedAt:            &redeemDate,
		ProductType:       param.ProductType,
		Status:            saveStatus,
		ExpDate:           utils.DefaultNulTime(ExpireDate),
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
		RedeemAt:          utils.DefaultNulTime(redeemDate),
		Comment:           param.Comment,
		MRewardID:         param.RewardID,
		ProductCategoryID: param.CategoryID,
		MProductID:        param.ProductID,
		PointsTransferID:  param.PointTransferID,
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

func SaveTSchedulerRetry(trxID, code string) {

	fmt.Println(fmt.Sprintf("[Start-SaveTSchedulerRetry]-[Sepulsa]"))

	schedulerData := dbmodels.TSchedulerRetry{
		Code:          code,
		TransactionID: trxID,
		Count:         0,
		IsDone:        false,
		CreatedAT:     time.Now(),
	}

	errSaveScheduler := db.DbCon.Create(&schedulerData).Error
	if errSaveScheduler != nil {

		fmt.Println("===== Gagal SaveScheduler ke DB =====")
		fmt.Println(fmt.Sprintf("Error : %v", errSaveScheduler))

	}

}
