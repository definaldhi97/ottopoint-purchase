package v2_migrate

import (
	"fmt"
	"net/http"
	"ottopoint-purchase/controllers"
	"ottopoint-purchase/models"
	"ottopoint-purchase/services/v2_migrate"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	ottologer "ottodigital.id/library/logger"
	utilsgo "ottodigital.id/library/utils"
)

func CallbackVoucherAggController(ctx *gin.Context) {

	fmt.Println("[ >>>>>>>>>>>>>>>>>>>>> Callbakc Voucher Agg Controller <<<<<<<<<<<<<<<<<< ]")
	var (
		req models.CallbackRequestVoucherAg
		res models.Response
	)

	sugarLogger := ottologer.GetLogger()
	namectrl := "[HandleCallbackRequestVoucherAg]"

	if err := ctx.ShouldBindJSON(&req); err != nil {
		res.Meta.Code = 03
		res.Meta.Message = "Gagal! Maaf transaksi Anda tidak dapat dilakukan saat ini. Silahkan dicoba lagi atau hubungi tim kami untuk informasi selengkapnya."
		ctx.JSON(http.StatusOK, res)
		go sugarLogger.Error("Error, body Request", zap.Error(err))
		return
	}

	span := controllers.TracingFirstControllerCtx(ctx, req, namectrl)
	c := ctx.Request.Context()
	context := opentracing.ContextWithSpan(c, span)

	spanid := utilsgo.GetSpanId(span)
	sugarLogger.Info("REQUEST:", zap.String("SPANID", spanid), zap.String("CTRL", namectrl),
		zap.Any("BODY", req),
		zap.Any("HEADER", ctx.Request.Header))

	VoucherAgMigrateServices := v2_migrate.VoucherAgMigrateServices{
		General: models.GeneralModel{
			ParentSpan: span,
			OttoZaplog: sugarLogger,
			SpanId:     spanid,
			Context:    context,
		},
	}

	res = VoucherAgMigrateServices.CallbackVoucherAgg(req)

	ctx.JSON(http.StatusOK, res)
}
