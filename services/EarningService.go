package services

import (
	"errors"
	"fmt"
	"math"
	"ottopoint-purchase/db"
	opl "ottopoint-purchase/hosts/opl/host"
	redismodels "ottopoint-purchase/hosts/redis_token/models"
	"ottopoint-purchase/models"
	"ottopoint-purchase/utils"
	"strconv"

	"ottopoint-purchase/constants"

	"github.com/astaxie/beego/logs"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type EarningServices struct {
	General models.GeneralModel
}

// func (t EarningServices) EarningPoint(req models.RulePointReq, dataToken redismodels.TokenResp, header models.RequestHeader) models.Response {
// 	var res models.Response

// 	sugarLogger := t.General.OttoZaplog
// 	sugarLogger.Info("[EarningPoint-Services]",
// 		zap.String("rule : ", req.EventName), zap.Int("amount : ", req.Amount), zap.String("institution", header.InstitutionID))

// 	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[EarningPoint-Services]")
// 	defer span.Finish()

// 	// Get CustID OPL from DB
// 	dataDB, errDB := db.CheckUser(dataToken.Data)
// 	if errDB != nil || dataDB.CustID == "" {
// 		logs.Info("Internal Server Error : ", errDB)
// 		logs.Info("[EarningPoint-Services]")
// 		logs.Info("[Get CustId OPL to DB]")

// 		//sugarLogger.Info("Internal Server Error : ", errDB)
// 		sugarLogger.Info("[EarningPoint-Services]")
// 		sugarLogger.Info("[Get CustId OPL to DB]")
// 		res = utils.GetMessageFailedErrorNew(res, constants.RC_ERROR_ACC_NOT_ELIGIBLE, constants.RD_ERROR_ACC_NOT_ELIGIBLE)
// 		//utils.GetMessageResponse(res, 422, false, errors.New("Nomor belum eligible"))
// 		return res
// 	}

// 	// Get ALL RulePoint
// 	getRule, errgetRule := opl.ListRulePoint(dataToken.Data)
// 	if errgetRule != nil || len(getRule.EarningRules) == 0 {
// 		sugarLogger.Info("[Error-dataPoint :")
// 		logs.Info("[Error-dataPoint :", errgetRule)
// 		res = utils.GetMessageFailedErrorNew(res, constants.RC_ERROR_FAILED_GET_POINT, constants.RD_ERROR_FAILED_GET_POINT)
// 		//res = utils.GetMessageResponse(res, 422, false, errors.New("Gagal Dapat Point"))
// 		return res
// 	}

// 	a := []models.GetEarningRulesResp{}
// 	for _, val := range getRule.EarningRules {
// 		b := models.GetEarningRulesResp{
// 			Name:         val.Name,
// 			EventName:    val.EventName,
// 			PointsAmount: val.PointsAmount,
// 		}
// 		a = append(a, b)
// 	}

// 	var product string
// 	// var amountpoint int
// 	for _, value := range a {
// 		if value.EventName == req.EventName {
// 			product = value.Name
// 			// name = value.Name
// 			// amountpoint = value.PointsAmount
// 		}
// 	}

// 	if product == "" {
// 		logs.Info("[EarningPoint-Services]")
// 		logs.Info("[Product Kosong]")

// 		sugarLogger.Info("[EarningPoint-Services]")
// 		sugarLogger.Info("[Product Kosong]")
// 		res = utils.GetMessageFailedErrorNew(res, constants.RC_ERROR_FAILED_GET_POINT, constants.RD_ERROR_FAILED_GET_POINT)
// 		//res = utils.GetMessageResponse(res, 422, false, errors.New("Gagal Dapat Point"))
// 		return res
// 	}

// 	// Earning by Denom / Product
// 	if req.Amount == 0 && req.EventName != "" {

// 		dataPoint, errPoint := opl.RulePoint(req.EventName, dataToken.Data)
// 		if errPoint != nil || dataPoint.Point == 0 {
// 			sugarLogger.Info("[Error-dataPoint :")
// 			logs.Info("[Error-dataPoint :", errPoint)
// 			res = utils.GetMessageFailedErrorNew(res, constants.RC_ERROR_FAILED_GET_POINT, constants.RD_ERROR_FAILED_GET_POINT)
// 			//res = utils.GetMessageResponse(res, 422, false, errors.New("Gagal Dapat Point"))
// 			return res
// 		}

// 		res = models.Response{
// 			Data: models.RulePointResp{
// 				Point:       dataPoint.Point,
// 				Product:     product,
// 				Institution: header.InstitutionID,
// 			},
// 			Meta: utils.ResponseMetaOK(),
// 		}

// 		return res
// 	}

// 	//Get Config in DB
// 	logs.Info("=== Get Config in DB ===")
// 	dataConf, errConf := db.GetConfig()
// 	if errConf != nil {
// 		sugarLogger.Info("[EEROR-DATABASE]-[EarningPoint-Services]-[GetConfig-DB]")
// 		sugarLogger.Info(fmt.Sprintf("Failed Get Data from DB getConfig %v", errConf))

// 		logs.Info("[EEROR-DATABASE]-[EarningPoint-Services]-[GetConfig-DB]")
// 		logs.Info(fmt.Sprintf("Failed Get Data from DB getConfig %v", errConf))

// 		res = utils.GetMessageResponse(res, 422, false, errors.New("Gagal! Maaf transaksi Anda tidak dapat dilakukan saat ini. Silahkan dicoba lagi atau hubungi tim kami untuk informasi selengkapnya."))

// 		return res
// 	}

// 	// perhitungan point
// 	point := float64(req.Amount) * dataConf.TransaksiPPOB
// 	// Ngambil angka di depan koma
// 	point = math.Floor(point)
// 	totalPoint := fmt.Sprintf("%.f", point)

// 	total, _ := strconv.Atoi(totalPoint)

// 	// validate jika melebihi limit trx
// 	if total > dataConf.LimitTransaksi {

// 		total = dataConf.LimitTransaksi
// 	}

// 	text := fmt.Sprintf("Transaksi %v", product)

// 	// Hit to Openloyalty
// 	dataTfPoint, errTfPoint := opl.TransferPoint(dataDB.CustID, strconv.Itoa(total), text)
// 	if errTfPoint != nil || dataTfPoint.PointsTransferId == "" {

// 		logs.Info("Internal Server Error : ", errTfPoint)
// 		logs.Info("[EarningPoint-Services]")
// 		logs.Info("[Hit Transfer API to OPL]")

// 		//sugarLogger.Info("Internal Server Error : ", errTfPoint)
// 		sugarLogger.Info("[EarningPoint-Services]")
// 		sugarLogger.Info("[Hit Transfer API to OPL]")
// 		res = utils.GetMessageFailedErrorNew(res, constants.RC_ERROR_FAILED_GET_POINT, constants.RD_ERROR_FAILED_GET_POINT)
// 		//res = utils.GetMessageResponse(res, 422, false, errors.New("Gagal Transfer Point"))
// 		return res
// 	}

// 	res = models.Response{
// 		Data: models.RulePointResp{
// 			Point:       total,
// 			Product:     product,
// 			Institution: header.InstitutionID,
// 		},
// 		Meta: utils.ResponseMetaOK(),
// 	}

// 	return res

// }

func (t EarningServices) EarningPoint(req models.RulePointReq, dataToken redismodels.TokenResp, header models.RequestHeader) models.Response {
	res := models.Response{}

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[EarningPoint-Services]",
		zap.String("rule : ", req.EventName), zap.Int("amount : ", req.Amount), zap.String("institution", header.InstitutionID))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[EarningPoint-Services]")
	defer span.Finish()

	// Get CustID OPL from DB
	dataDB, errDB := db.CheckUser(dataToken.Data)
	if errDB != nil || dataDB.CustID == "" {
		logs.Info("Internal Server Error : ", errDB)
		logs.Info("[EarningPoint-Services]")
		logs.Info("[Get CustId OPL to DB]")

		sugarLogger.Info("Internal Server Error : ")
		sugarLogger.Info("[EarningPoint-Services]")
		sugarLogger.Info("[Get CustId OPL to DB]")
		res = utils.GetMessageFailedErrorNew(res, constants.RC_ERROR_ACC_NOT_ELIGIBLE, constants.RD_ERROR_ACC_NOT_ELIGIBLE)
		//utils.GetMessageResponse(res, 422, false, errors.New("Nomor belum eligible"))
		return res
	}

	// Get ALL RulePoint
	getRule, errgetRule := opl.ListRulePoint(dataToken.Data)
	if errgetRule != nil || len(getRule.EarningRules) == 0 {
		sugarLogger.Info("[Error-dataPoint :")
		logs.Info("[Error-dataPoint :", errgetRule)
		res = utils.GetMessageFailedErrorNew(res, constants.RC_ERROR_FAILED_GET_POINT, constants.RD_ERROR_FAILED_GET_POINT)
		//res = utils.GetMessageResponse(res, 422, false, errors.New("Gagal Dapat Point"))
		return res
	}

	a := []models.GetEarningRulesResp{}
	for _, val := range getRule.EarningRules {
		b := models.GetEarningRulesResp{
			Name:         val.Name,
			EventName:    val.EventName,
			PointsAmount: val.PointsAmount,
		}
		a = append(a, b)
	}

	var product string
	// var amountpoint int
	for _, value := range a {
		if value.EventName == req.EventName {
			product = value.EventName
			// name = value.Name
			// amountpoint = value.PointsAmount
		}
	}

	if product == "" {
		logs.Info("[EarningPoint-Services]")
		logs.Info("[Product Kosong]")

		sugarLogger.Info("[EarningPoint-Services]")
		sugarLogger.Info("[Product Kosong]")
		res = utils.GetMessageFailedErrorNew(res, constants.RC_ERROR_FAILED_GET_POINT, constants.RD_ERROR_FAILED_GET_POINT)
		//res = utils.GetMessageResponse(res, 422, false, errors.New("Gagal Dapat Point"))
		return res
	}

	// Earning by Denom / Product
	if req.Amount == 0 && req.EventName != "" {

		dataPoint, errPoint := opl.RulePoint(req.EventName, dataToken.Data)
		if errPoint != nil || dataPoint.Point == 0 {
			sugarLogger.Info("[Error-dataPoint :")
			logs.Info("[Error-dataPoint :", errPoint)
			res = utils.GetMessageFailedErrorNew(res, constants.RC_ERROR_FAILED_GET_POINT, constants.RD_ERROR_FAILED_GET_POINT)
			//res = utils.GetMessageResponse(res, 422, false, errors.New("Gagal Dapat Point"))
			return res
		}

		res = models.Response{
			Data: models.RulePointResp{
				Point:       dataPoint.Point,
				Product:     product,
				Institution: header.InstitutionID,
			},
			Meta: utils.ResponseMetaOK(),
		}

		return res
	}

	//Get Config in DB
	logs.Info("=== Get Config in DB ===")
	dataConf, errConf := db.GetConfig()
	if errConf != nil {
		sugarLogger.Info("[EEROR-DATABASE]-[EarningPoint-Services]-[GetConfig-DB]")
		sugarLogger.Info(fmt.Sprintf("Failed Get Data from DB getConfig %v", errConf))

		logs.Info("[EEROR-DATABASE]-[EarningPoint-Services]-[GetConfig-DB]")
		logs.Info(fmt.Sprintf("Failed Get Data from DB getConfig %v", errConf))

		res = utils.GetMessageResponse(res, 422, false, errors.New("Gagal! Maaf transaksi Anda tidak dapat dilakukan saat ini. Silahkan dicoba lagi atau hubungi tim kami untuk informasi selengkapnya."))

		return res
	}

	var percentage float64
	var limit int
	if product == "earning indomarco" {
		percentage = dataConf.TransaksiMerchant
		limit = dataConf.LimitTransaksi

	} else {
		percentage = dataConf.TransaksiMerchant
		limit = dataConf.LimitTransaksi

	}

	// perhitungan point
	point := float64(req.Amount) * percentage
	// Ngambil angka di depan koma
	point = math.Floor(point)
	totalPoint := fmt.Sprintf("%.f", point)

	total, _ := strconv.Atoi(totalPoint)

	// validate jika melebihi limit trx
	if total > limit {

		total = limit
	}

	text := fmt.Sprintf("Transaksi %v", product)

	// Hit to Openloyalty
	dataTfPoint, errTfPoint := opl.TransferPoint(dataDB.CustID, strconv.Itoa(total), text)
	if errTfPoint != nil || dataTfPoint.PointsTransferId == "" {

		logs.Info("Internal Server Error : ", errTfPoint)
		logs.Info("[EarningPoint-Services]")
		logs.Info("[Hit Transfer API to OPL]")

		sugarLogger.Info("Internal Server Error : ")
		sugarLogger.Info("[EarningPoint-Services]")
		sugarLogger.Info("[Hit Transfer API to OPL]")
		res = utils.GetMessageFailedErrorNew(res, constants.RC_ERROR_FAILED_GET_POINT, constants.RD_ERROR_FAILED_GET_POINT)
		//res = utils.GetMessageResponse(res, 422, false, errors.New("Gagal Transfer Point"))
		return res
	}

	res = models.Response{
		Data: models.RulePointResp{
			Product:     product,
			Point:       total,
			Institution: header.InstitutionID,
		},
		Meta: utils.ResponseMetaOK(),
	}

	return res

}
