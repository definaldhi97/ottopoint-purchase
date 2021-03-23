package controllers

import (
	"fmt"
	"ottopoint-purchase/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vjeantet/jodaTime"

	"net/http"

	"ottopoint-purchase/models"
)

func CreateFileCSVController(ctx *gin.Context) {
	req := models.CreateCSV{}
	res := models.Response{}

	namectrl := "[PackageControllers]-[CreateFileCSVController]"

	if err := ctx.ShouldBindJSON(&req); err != nil {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[ShouldBindJSON]-[Error : %v]", err))
		logrus.Println("Request : ", req)

		res.Meta.Code = 03
		res.Meta.Message = "Gagal! Maaf transaksi Anda tidak dapat dilakukan saat ini. Silahkan dicoba lagi atau hubungi tim kami untuk informasi selengkapnya."
		ctx.JSON(http.StatusOK, res)
		return
	}

	name := jodaTime.Format("dd-MM-YYYY", time.Now()) + ".csv"

	go utils.CreateCSVFile(req, name)

	ctx.JSON(http.StatusOK, res)

	return

}
