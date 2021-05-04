package controllers

import (
	"net/http"
	"ottopoint-purchase/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func PatchingInvoiceNumberController(ctx *gin.Context) {

	namectrl := "[PackageController]-[PatchingInvoiceNumberController]"

	logrus.Info(namectrl)

	res := services.PatchingInvoiceNumberService()
	ctx.JSON(http.StatusOK, res)
	return

}
