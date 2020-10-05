package controllers

import (
	"fmt"
	"ottopoint-purchase/constants"
	hostAuth "ottopoint-purchase/hosts/auth/host"
	signature "ottopoint-purchase/hosts/signature/host"
	"ottopoint-purchase/models"
	"ottopoint-purchase/services"
	"ottopoint-purchase/utils"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
)

func ValidateRequest(ctx *gin.Context, isCheckAuth bool, reqBody interface{}, isCacheBalance bool) (models.RequestHeader, models.Response) {

	var phone string

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

	var result models.Response
	result.Meta.Status = true

	//check header request
	if header.AppsID == "" || header.ChannelID == "" || header.InstitutionID == "" || header.DeviceID == "" || header.Geolocation == "" {
		result = utils.GetMessageFailedErrorNew(result, constants.RC_ERROR_HEADER_MANDATORY, constants.RD_ERROR_HEADER_MANDATORY)
		return header, result
	}

	//check token
	if isCheckAuth {
		var RedisService = new(services.RedisService)
		auth := strings.ReplaceAll(header.Authorization, "Bearer ", "")
		keyRedis := header.InstitutionID + "-" + auth
		dataRedis := RedisService.GetData(keyRedis)

		if dataRedis.ResponseCode != "00" {
			result = utils.GetMessageFailedErrorNew(result, constants.RC_ERROR_INVALID_TOKEN, constants.RD_ERROR_INVALID_TOKEN)
			return header, result
		}

		phone = dataRedis.Value

		// result := ValidateUser(ctx, dataRedis.Value, result, header)

		// if !result.Meta.Status {
		// 	return header, result
		// }

	}

	ValidateSignature, errSignature := signature.Signature(reqBody, header)

	println("ValidateSignature === ", ValidateSignature.ResponseCode)
	if errSignature != nil || ValidateSignature.ResponseCode != "00" {
		logs.Info("[ValidateSignature]-[VoucherRedeemController]")
		logs.Info(fmt.Sprintf("Error when validation request header"))
		result = utils.GetMessageFailedErrorNew(result, constants.RC_ERROR_INVALID_SIGNATURE, constants.RD_ERROR_INVALID_SIGNATURE)
		return header, result
	}

	// Clear Cache Balance Point
	if isCacheBalance {
		go ClearCaceheBalancePoint(phone)
	}

	return header, result
}

func ValidateRequestWithoutAuth(ctx *gin.Context, reqBody interface{}) (models.RequestHeader, models.Response) {

	header := models.RequestHeader{
		DeviceID:      ctx.Request.Header.Get("DeviceId"),
		InstitutionID: ctx.Request.Header.Get("InstitutionId"),
		Geolocation:   ctx.Request.Header.Get("Geolocation"),
		ChannelID:     ctx.Request.Header.Get("ChannelId"),
		AppsID:        ctx.Request.Header.Get("AppsId"),
		Timestamp:     ctx.Request.Header.Get("Timestamp"),
		Signature:     ctx.Request.Header.Get("Signature"),
	}

	var result models.Response
	result.Meta.Status = true

	//check header request
	if header.AppsID == "" || header.ChannelID == "" || header.InstitutionID == "" || header.DeviceID == "" || header.Geolocation == "" {
		result = utils.GetMessageFailedErrorNew(result, constants.RC_ERROR_HEADER_MANDATORY, constants.RD_ERROR_HEADER_MANDATORY)
		return header, result
	}

	ValidateSignature, errSignature := signature.Signature(reqBody, header)

	println("ValidateSignature === ", ValidateSignature.ResponseCode)
	if errSignature != nil || ValidateSignature.ResponseCode != "00" {
		logs.Info("[ValidateSignature]-[VoucherRedeemController]")
		logs.Info(fmt.Sprintf("Error when validation request header"))
		result = utils.GetMessageFailedErrorNew(result, constants.RC_ERROR_INVALID_SIGNATURE, constants.RD_ERROR_INVALID_SIGNATURE)
		return header, result
	}

	return header, result
}

func ClearCaceheBalancePoint(phone string) {
	fmt.Println(">>>>>>> Clear Cache Get Balance <<<<<<")
	clearCacheBalance, err := hostAuth.ClearCacheBalance(phone)
	if err != nil {
		fmt.Println("Clear Cache Balance Error : ", err)
		return
	}
	if clearCacheBalance.ResponseCode != "00" {
		fmt.Println("Message : ", clearCacheBalance.Messages)
		fmt.Println("Response Code : ", clearCacheBalance.ResponseCode)
		return
	}
	fmt.Println("Clear Cache Get Balance: ", clearCacheBalance.Messages)
	return

}
