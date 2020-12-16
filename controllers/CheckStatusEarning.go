package controllers

import (
	"fmt"
	services "ottopoint-purchase/services/earnings"
	"ottopoint-purchase/utils"
	"time"

	"github.com/gin-gonic/gin"
	zaplog "github.com/opentracing-contrib/go-zap/log"
	"go.uber.org/zap"

	"net/http"

	"ottopoint-purchase/models"
)

func CheckStatusEarningController(ctx *gin.Context) {
	req := models.CheckStatusEarningReq{}
	res := models.Response{}

	namectrl := "[CheckStatusEarningController]"

	if err := ctx.ShouldBindJSON(&req); err != nil {

		fmt.Println("[Error, body Request]-[CheckStatusEarningController]")
		fmt.Println(fmt.Sprintf("[Error : %v]", err))

		res.Meta.Code = 03
		res.Meta.Message = "Gagal! Maaf transaksi Anda tidak dapat dilakukan saat ini. Silahkan dicoba lagi atau hubungi tim kami untuk informasi selengkapnya."
		ctx.JSON(http.StatusOK, res)

		return
	}

	span := TracingFirstControllerCtx(ctx, req, namectrl)

	// validate request
	header, resultValidate := ValidateRequestWithoutAuth(ctx, req)
	if !resultValidate.Meta.Status {
		ctx.JSON(http.StatusOK, resultValidate)
		return
	}

	// dataToken, _ := token.CheckToken(header)

	checkStatusEarning := new(services.CheckStatusEarningService)

	res = checkStatusEarning.CheckStatusEarningServices(req.ReferenceId, header.InstitutionID)

	datalog := utils.LogSpanMax(res)
	zaplog.InfoWithSpan(span, namectrl,
		zap.Any("RESP", datalog),
		zap.Duration("backoff", time.Second))

	defer span.Finish()
	ctx.JSON(http.StatusOK, res)

}
