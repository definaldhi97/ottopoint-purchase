package CallbackVoucher

import (
	"net/http"
	"ottopoint-purchase/models"

	sepulsaModels "ottopoint-purchase/hosts/sepulsa/models"
	"ottopoint-purchase/services/v2/Redeemtion"

	"ottopoint-purchase/controllers"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	ottologer "ottodigital.id/library/logger"
	utilsgo "ottodigital.id/library/utils"
)

type V2_CallbackSepulsaController struct{}

func (controller *V2_CallbackSepulsaController) VoucherCallbackSepulsaController(ctx *gin.Context) {

	logrus.Info("[ >>>>>>>>>>>>>>>>>>>>> Callbakc Sepulsa COntroller <<<<<<<<<<<<<<<<<< ]")

	req := sepulsaModels.CallbackTrxReq{}
	res := models.Response{}

	sugarLogger := ottologer.GetLogger()
	namectrl := "[CallbackSepulsaController]"

	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.Meta.Code = 03
		res.Meta.Message = err.Error()
		ctx.JSON(http.StatusOK, res)
		go sugarLogger.Error("Error, body request", zap.Error(err))
		return
	}

	span := controllers.TracingFirstControllerCtx(ctx, req, namectrl)
	c := ctx.Request.Context()
	context := opentracing.ContextWithSpan(c, span)
	spanId := utilsgo.GetSpanId(span)

	VoucherSepulsaMigrateService := Redeemtion.V2_VoucherSepulsaService{
		General: models.GeneralModel{
			ParentSpan: span,
			OttoZaplog: sugarLogger,
			SpanId:     spanId,
			Context:    context,
		},
	}

	res = VoucherSepulsaMigrateService.CallbackVoucherSepulsa(req)

	// sugarLogger.Info("RESPONSE : ", zap.String("SPANID", spanId), zap.String("CTRL", namectrl),
	// 	zap.Any("BODY : ", res))

	// datalog := utils.LogSpanMax(res)
	// zaplog.InfoWithSpan(span, namectrl,
	// 	zap.Any("RESP : ", datalog),
	// 	zap.Duration("backoff : ", time.Second))

	defer span.Finish()

	ctx.JSON(http.StatusOK, res)

}