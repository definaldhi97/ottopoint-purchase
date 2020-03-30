package services

import (
	"errors"
	"fmt"
	db "ottopoint-purchase/db"
	opl "ottopoint-purchase/hosts/opl/host"
	"ottopoint-purchase/models"
	"ottopoint-purchase/models/dbmodels"
	"ottopoint-purchase/utils"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

type DeductSplitBillServices struct {
	General models.GeneralModel
}

func (t DeductSplitBillServices) DeductSplitBill(req models.DeductPointReq, accountNumber, institution string) models.Response {
	var res models.Response

	resMeta := models.MetaData{
		Code:    200,
		Status:  true,
		Message: "Succesful",
	}

	resData := models.ResponseData{
		ResponseCode: "05",
		ResponseDesc: "Error",
	}

	sugarLogger := t.General.OttoZaplog
	sugarLogger.Info("[DeductSplitBill-Services]",
		zap.String("AccountNumber : ", accountNumber), zap.Int("Point : ", req.Point),
		zap.Int("deductType : ", req.DeductType), zap.String("trxID : ", req.TrxID),
		zap.Int("amount : ", req.Amount), zap.String("productCode : ", req.ProductCode),
		zap.String("productName : ", req.ProductName), zap.String("InstitutionID : ", institution))

	span, _ := opentracing.StartSpanFromContext(t.General.Context, "[DeductSplitBill-Services]")
	defer span.Finish()

	// Get CustID OPL from DB
	dataDB, errDB := db.CheckUser(accountNumber)
	if errDB != nil || dataDB.CustID == "" {
		logs.Info("Internal Server Error : ", errDB)
		logs.Info("[DeductSplitBill-Services]")
		logs.Info("[Get CustId OPL to DB]")

		// sugarLogger.Info("Internal Server Error : ", errDB)
		// sugarLogger.Info("[DeductSplitBill-Services]")
		// sugarLogger.Info("[Get CustId OPL to DB]")

		res = utils.GetMessageResponseData(res, resData, 422, false, errors.New("Nomor belum eligible"))
		return res
	}

	logs.Info("CustID OPL : ", dataDB.CustID)
	// sugarLogger.Info("CustID OPL : ", dataDB.CustID)

	// Cek Balance
	dataBalance, errBalance := opl.GetBalance(dataDB.CustID)
	if errBalance != nil {
		logs.Info("Internal Server Error : ", errDB)
		logs.Info("[DeductSplitBill-Services]")
		logs.Info("[Get CustId OPL to DB]")

		// sugarLogger.Info("Internal Server Error : ", errDB)
		sugarLogger.Info("[DeductSplitBill-Services]")
		sugarLogger.Info("[Get CustId OPL to DB]")

		res = utils.GetMessageResponseData(res, resData, 422, false, errors.New("Failed to GetBalance"))
		return res
	}

	// Validate TRX ID
	dataTrx, errTrx := db.GetData(req.TrxID, institution)
	if errTrx == nil {
		logs.Info("Internal Server Error : ", errDB)
		logs.Info("[DeductSplitBill-Services]")
		logs.Info("[Get Data TRXID to DB]")

		// sugarLogger.Info("Internal Server Error : ", errDB)
		sugarLogger.Info("[DeductSplitBill-Services]")
		sugarLogger.Info("[Get Data TRXID to DB]")

		// res = utils.GetMessageResponseData(res, resData, 422, false, errors.New("Failed to GetTRxID"))
		// return res
	}

	if req.TrxID == dataTrx.TrxID {
		// logs.Info("Internal Server Error : ", errDB)
		logs.Info("[DeductSplitBill-Services]")
		logs.Info("[Validate TrxID]")

		// sugarLogger.Info("Internal Server Error : ", errDB)
		sugarLogger.Info("[DeductSplitBill-Services]")
		sugarLogger.Info("Validate TrxID]")

		res = utils.GetMessageResponseData(res, resData, 422, false, errors.New("Duplicate TrxID"))
		return res
	}

	point := int(dataBalance.Points)

	// Validate Balance Point
	if req.Point > point {
		logs.Info("Request Point :", req.Point)
		logs.Info("Balance Point :", point)
		logs.Info("[DeductSplitBill-Services]-[Point Tidak Mencukupi]")
		sugarLogger.Info("[DeductSplitBill-Services]-[Point Tidak Mencukupi]")

		res = models.Response{
			Meta: resMeta,
			Data: models.ResponseData{
				ResponseCode: "27",
				ResponseDesc: "Point Anda Tidak Mencukupi",
			},
		}

		return res
	}

	text := fmt.Sprintf("Deduct Point from %v", req.ProductName)
	// Hit to Openloyalty
	data, err := opl.SpendPoint(dataDB.CustID, strconv.Itoa(req.Point), text)
	if err != nil || data.PointsTransferId == "" {

		logs.Info("Internal Server Error : ", err)
		logs.Info("[DeductSplitBill-Services]")
		logs.Info("[Hit Transfer API to OPL]")

		// sugarLogger.Info("Internal Server Error : ", err)
		sugarLogger.Info("[DeductSplitBill-Services]")
		sugarLogger.Info("[Hit Spend API to OPL]")

		res = utils.GetMessageResponseData(res, resData, 422, false, errors.New("Gagal Transfer Point"))

		return res
	}

	// save to DB
	save := dbmodels.DeductTransaction{
		ID:            utils.GenerateTokenUUID(),
		TrxID:         req.TrxID,
		AccountID:     accountNumber,
		CustomerID:    req.CustID,
		InstitutionID: institution,
		DeductType:    req.DeductType,
		ProductCode:   req.ProductCode,
		ProductName:   req.ProductName,
		Amount:        req.Amount,
		Point:         req.Point,
		Status:        "00",
	}

	errSave := db.DbCon.Create(&save).Error
	if errSave != nil {
		logs.Info("[DeductSplitBill-Services]")
		logs.Info("Failed Save to database", errSave)

		sugarLogger.Info("[DeductSplitBill-Services]")
		// sugarLogger.Info("Failed Save to database", errSave)
		// return errSave
	}

	res = models.Response{
		Data: models.ResponseData{
			ResponseCode: "00",
			ResponseDesc: "Transaction Success",
		},
		Meta: resMeta,
	}

	return res
}
