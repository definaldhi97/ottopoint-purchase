package services

import (
	"errors"
	db "ottopoint-purchase/db"
	opl "ottopoint-purchase/hosts/opl/host"
	redismodels "ottopoint-purchase/hosts/redis_token/models"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type SpendPointServices struct {
	General models.GeneralModel
}

func (t SpendPointServices) NewSpendPointServices(req models.PointReq, dataToken redismodels.TokenResp, header models.RequestHeader) models.Response {
	var res models.Response

	resMeta := models.MetaData{
		Code:    200,
		Status:  true,
		Message: "Succesful",
	}

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[SpendPoint-Services]",
		zap.String("AccountNumber : ", req.AccountNumber), zap.Int("Point : ", req.Point),
		zap.String("Text : ", req.Text), zap.String("Type : ", req.Type))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[SpendPoint-Service]")
	defer span.Finish()

	// Get CustID OPL from DB
	dataDB, errDB := db.CheckUser(req.AccountNumber)
	if errDB != nil || dataDB.CustID == "" {
		logs.Info("Internal Server Error : ", errDB)
		logs.Info("[SpendPoint-Services]")
		logs.Info("[Get CustId OPL to DB]")

		sugarLogger.Info("Internal Server Error : ", errDB)
		sugarLogger.Info("[SpendPoint-Services]")
		sugarLogger.Info("[Get CustId OPL to DB]")

		utils.GetMessageResponse(res, 422, false, errors.New("Nomor belum eligible"))
		return res
	}

	logs.Info("CustID OPL : ", dataDB.CustID)
	sugarLogger.Info("CustID OPL : ", dataDB.CustID)

	// Hit to Openloyalty
	data, err := opl.SpendPoint(dataDB.CustID, strconv.Itoa(req.Point), req.Text)
	if err != nil || data.PointsTransferId == "" {

		logs.Info("Internal Server Error : ", err)
		logs.Info("[SpendPoint-Services]")
		logs.Info("[Hit Transfer API to OPL]")

		sugarLogger.Info("Internal Server Error : ", err)
		sugarLogger.Info("[SpendPoint-Services]")
		sugarLogger.Info("[Hit Transfer API to OPL]")

		res = utils.GetMessageResponse(res, 422, false, errors.New("Gagal Transfer Point"))
		return res
	}

	// logs.Info("[Send Notif]")
	// notif := models.NotifReq{
	// 	CollapseKey: "type_ottopoint",
	// 	Title:       "Selamat! Poin kamu bertambah",
	// 	Body:        fmt.Sprintf("Kamu mendapatkan %v poin dari OttoPay, makin sering transaksi makin untung.", int64(req.Point)),
	// 	Target:      "earning point",
	// 	Phone:       req.Phone,
	// 	Rc:          "00",
	// }

	// logs.Info("[Request Send Notif : ]", notif)

	// dataNotif, errNotif := ottomart.NotifInboxOttomart(notif, header)
	// if errNotif != nil {
	// 	res = utils.GetMessageFailedError(res, 422, errors.New("Error to send Notif & Inbox"))
	// 	return res
	// }

	res = models.Response{
		Data: models.PointResp{
			Nama:          dataDB.Nama,
			AccountNumber: dataDB.Phone,
			Point:         req.Point,
			Text:          req.Text,
		},
		Meta: resMeta,
	}

	return res
}
