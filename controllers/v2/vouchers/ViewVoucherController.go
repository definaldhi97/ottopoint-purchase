package vouchers

import (
	"net/http"
	"ottopoint-purchase/constants"
	redisService "ottopoint-purchase/hosts/redis_token/host"
	"ottopoint-purchase/models"
	services "ottopoint-purchase/services/v2/vouchers"
	"ottopoint-purchase/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ViewVoucherController(ctx *gin.Context) {

	namectrl := "[PackageVouchers]-[ViewVoucherController]"

	logrus.Info(namectrl)

	var resp models.Response

	couponId := ctx.Request.URL.Query().Get("couponId")

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
		resp = utils.GetMessageFailedErrorNew(resp, constants.RC_ERROR_HEADER_MANDATORY, constants.RD_ERROR_HEADER_MANDATORY)
		ctx.JSON(http.StatusOK, resp)
		return
	}

	// Validate Token user
	auth := strings.ReplaceAll(header.Authorization, "Bearer ", "")
	keyRedis := header.InstitutionID + "-" + auth
	dataRedis, _ := redisService.GetToken(keyRedis)

	if dataRedis.ResponseCode != "00" {
		resp = utils.GetMessageFailedErrorNew(resp, constants.RC_ERROR_INVALID_TOKEN, constants.RD_ERROR_INVALID_TOKEN)
		ctx.JSON(http.StatusOK, resp)
		return
	}

	logrus.Println("[Request]")
	logrus.Info("CouponId : ", couponId)

	// service
	resp = services.ViewVoucherServices(dataRedis.Data, couponId)
	ctx.JSON(http.StatusOK, resp)

	return

}
