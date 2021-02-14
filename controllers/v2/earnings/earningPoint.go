package earnings

import (
	"encoding/json"
	"errors"
	"fmt"
	"ottopoint-purchase/constants"
	ctrl "ottopoint-purchase/controllers"
	kafka "ottopoint-purchase/hosts/publisher/host"
	worker "ottopoint-purchase/hosts/worker/host"
	modelsworker "ottopoint-purchase/hosts/worker/models"
	"ottopoint-purchase/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"net/http"

	"ottopoint-purchase/models"
)

func EarningsPointController(ctx *gin.Context) {
	req := models.EarningReq{}
	res := models.Response{}

	namectrl := "[PackageEarnings]-[EarningsPointController]"

	logrus.Info(namectrl)

	if err := ctx.ShouldBindJSON(&req); err != nil {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[ShouldBindJSON]-[Error : %v]", err))
		logrus.Println("Request : ", req)

		res.Meta.Code = 03
		res.Meta.Message = "Gagal! Maaf transaksi Anda tidak dapat dilakukan saat ini. Silahkan dicoba lagi atau hubungi tim kami untuk informasi selengkapnya."
		ctx.JSON(http.StatusOK, res)
		return
	}

	// validate request
	header, resultValidate := ctrl.ValidateRequestWithoutAuth(ctx, req)
	if !resultValidate.Meta.Status {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[ValidateRequestWithoutAuth]-[Error : %v]", resultValidate))
		logrus.Println("Request : ", req)

		ctx.JSON(http.StatusOK, resultValidate)
		return
	}

	if req.Earning == "" || req.TransactionTime == "" || req.AccountNumber1 == "" || len(req.TransactionTime) != 19 {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[ValidateRequestMandatory]-[Invalid Mandatory]"))
		logrus.Println("Request : ", req)

		res = utils.GetMessageResponse(res, 61, false, errors.New("Invalid Mandatory"))

		ctx.JSON(http.StatusOK, res)

		return
	}

	logrus.Println("[Request]")

	logrus.Info(
		"Earning : ", req.Earning, "ReferenceId : ", req.ReferenceId, "ProductCode : ", req.ProductCode, "ProductName : ", req.ProductName, "AccountNumber1 : ", req.AccountNumber1,
		"AccountNumber2 : ", req.AccountNumber2, "Amount : ", req.Amount, "Remark : ", req.Remark, "TransactionTime : ", req.TransactionTime)

	res = utils.GetMessageResponse(res, 200, true, errors.New("Transaksi sedang di proses"))

	code := req.Earning[:3]
	switch code {
	case constants.GeneralSpending:
		fmt.Println("===== GeneralSpending =====")
		go publishEarning(req, header)
	case constants.InstantReward:
		fmt.Println("===== InstantReward =====")
		go publishEarning(req, header)
	case constants.EventRule:
		fmt.Println("===== EventRule =====")
		go publishEarning(req, header)
	case constants.CustomerReferral:
		fmt.Println("===== CustomerReferral =====")
		go publishEarning(req, header)
	case constants.CustomeEventRule:
		fmt.Println("===== CustomeEventRule =====")
		go publishEarning(req, header)
	default:
		fmt.Println("===== Invalid Code =====")

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[ValidateCodeEarning]-[Invalid Code]"))
		logrus.Println("Request : ", req)

		res = utils.GetMessageResponse(res, 178, false, errors.New("Earning Rule not found"))
	}

	ctx.JSON(http.StatusOK, res)

}

func publishEarning(req models.EarningReq, header models.RequestHeader) {
	fmt.Println(">>>>> Publisher Earning <<<<<")

	layout := "2006-01-02 15:04:05"

	t, _ := time.Parse(layout, req.TransactionTime)

	pubReq := models.PublishEarningReq{
		Header:          header,
		Earning:         req.Earning,
		ReferenceId:     req.ReferenceId,
		ProductCode:     req.ProductCode,
		ProductName:     req.ProductName,
		AccountNumber1:  req.AccountNumber1,
		AccountNumber2:  req.AccountNumber2,
		Amount:          req.Amount,
		Remark:          req.Remark,
		TransactionTime: t,
	}

	bytePub, _ := json.Marshal(pubReq)

	kafkaReq := kafka.PublishReq{
		Topic: utils.TopicsEarning,
		Value: bytePub,
	}

	kafkaRes, errKafka := kafka.SendPublishKafka(kafkaReq)
	fmt.Println("Response Publisher : ", kafkaRes)
	if errKafka != nil || kafkaRes.ResponseCode != "00" {

		logrus.Warn("[PackageEarnings]-[EarningsPointController]")
		logrus.Warn(fmt.Sprintf("[SendPublishKafka-PublishEarning]-[Error : %v]", errKafka))
		logrus.Println("Request : ", kafkaReq)

		reqWorker := modelsworker.WorkerEarningReq{
			InstitutionId:   header.InstitutionID,
			AccountNumber1:  req.AccountNumber1,
			AccountNumber2:  req.AccountNumber2,
			Earning:         req.Earning,
			ReferenceId:     req.ReferenceId,
			ProductCode:     req.ProductCode,
			ProductName:     req.ProductName,
			Amount:          req.Amount,
			Remark:          req.Remark,
			TransactionTime: req.TransactionTime,
		}

		fmt.Println(">> Send to API-Worker <<")
		workerApi, errworker := worker.WorkerEarning(reqWorker)
		if errworker != nil || workerApi.Code != 200 {

			logrus.Error("[PackageEarnings]-[EarningsPointController]")
			logrus.Error(fmt.Sprintf("[WorkerEarning-PublishEarning]-[Error : %v]", errworker))
			logrus.Println("Request : ", reqWorker)

		}

	}
}
