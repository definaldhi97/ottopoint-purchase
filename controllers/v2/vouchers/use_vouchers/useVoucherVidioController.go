package use_vouchers

import (
	"fmt"
	"net/http"
	"ottopoint-purchase/constants"
	redishost "ottopoint-purchase/hosts/redis_token/host"
	"ottopoint-purchase/models"
	usevoucher "ottopoint-purchase/services/v2/vouchers/use_vouchers"
	"ottopoint-purchase/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func UseVoucherVidioController(ctx *gin.Context) {
	fmt.Println("[ >>>>>>>>>>>>>>>>>> Use Vouhcer Vidio Controller <<<<<<<<<<<<<<<< ]")

	var resp models.Response

	namectrl := "[PackageUserVoucher]-[UseVoucherVidioController]"

	logrus.Info(namectrl)

	couponId := ctx.Request.URL.Query().Get("couponId")

	logReq := fmt.Sprintf("[CouponID : %v ]", couponId)

	// header
	header := models.RequestHeader{
		DeviceID:      ctx.Request.Header.Get("DeviceId"),
		InstitutionID: ctx.Request.Header.Get("InstitutionId"),
		Geolocation:   ctx.Request.Header.Get("Geolocation"),
		ChannelID:     ctx.Request.Header.Get("ChannelId"),
		AppsID:        ctx.Request.Header.Get("AppsId"),
		Timestamp:     ctx.Request.Header.Get("Timestamp"),
		Authorization: ctx.Request.Header.Get("Authorization"),
		Signature:     ctx.Request.Header.Get("Signature"),
	}

	//check header request
	if header.AppsID == "" || header.ChannelID == "" || header.InstitutionID == "" || header.DeviceID == "" || header.Geolocation == "" {

		logrus.Error(namectrl)
		logrus.Error("[Invalid Mandatory]")
		logrus.Println(logReq)

		resp = utils.GetMessageFailedErrorNew(resp, constants.RC_ERROR_HEADER_MANDATORY, constants.RD_ERROR_HEADER_MANDATORY)
		ctx.JSON(http.StatusOK, resp)
		return
	}

	// Validate Token user
	auth := strings.ReplaceAll(header.Authorization, "Bearer ", "")
	keyRedis := header.InstitutionID + "-" + auth
	dataRedis, _ := redishost.GetToken(keyRedis)

	if dataRedis.ResponseCode != "00" {

		logrus.Error(namectrl)
		logrus.Error(fmt.Sprintf("[GetDatafromRedis]-[Error : %v]", dataRedis))
		logrus.Println(logReq)

		resp = utils.GetMessageFailedErrorNew(resp, constants.RC_ERROR_INVALID_TOKEN, constants.RD_ERROR_INVALID_TOKEN)
		ctx.JSON(http.StatusOK, resp)
		return
	}

	logrus.Println("[Request]")
	logrus.Info("CouponID : ", couponId)

	resp = usevoucher.UseVoucherVidioServices(couponId)

	ctx.JSON(http.StatusOK, resp)

	return

}
