package controllers

import (
	"fmt"
	"ottopoint-purchase/constants"
	signature "ottopoint-purchase/hosts/signature/host"
	"ottopoint-purchase/models"
	"ottopoint-purchase/services"
	"ottopoint-purchase/utils"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
)

func ValidateRequest(ctx *gin.Context, isCheckAuth bool, reqBody interface{}) (models.RequestHeader, models.Response) {

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

	return header, result
}

// func ValidateUser(ctx *gin.Context, phone string, result models.Response, header models.RequestHeader) models.Response {
// 	var UserService = new(services.UserService)

// 	logGen := BuildLogger("ValidateRequest", ctx, *UserService)
// 	UserService.General = logGen.General

// 	accUser, err := UserService.CheckAccount(phone)

// 	if err == nil && accUser.Status != constants.CONS_USER_STATUS_ACTIVE {
// 		result = utils.GetMessageFailedErrorNew(result, constants.RC_ERROR_USER_INACTIVE, constants.RD_ERROR_USER_INACTIVE)
// 		return result
// 	}

// 	accUserLink, err := UserService.CheckUserLink(accUser.Id, header.InstitutionID)

// 	if err == nil && accUserLink.Status != constants.CONS_USER_STATUS_ACTIVE {
// 		result = utils.GetMessageFailedErrorNew(result, constants.RC_ERROR_USER_LINKED_INACTIVE, constants.RD_ERROR_USER_LINKED_INACTIVE)
// 		return result
// 	}
// 	return result
// }

// func BuildLogger(ctrlName string, ctx *gin.Context, service services.UserService) services.UserService {
// 	sugarLogger := ottologger.GetLogger()

// 	span := TracingEmptyFirstControllerCtx(ctx, ctrlName)
// 	c := ctx.Request.Context()
// 	Context := opentracing.ContextWithSpan(c, span)
// 	defer span.Finish()

// 	spanID := utilsgo.GetSpanId(span)
// 	sugarLogger.Info(zap.String("SPANID", spanID), zap.String("CTRL", ctrlName),
// 		// zap.Any("BODY", request),
// 		zap.Any("HEADER", "header"))

// 	service.General = models.GeneralModel{
// 		ParentSpan: span,
// 		OttoZapLog: sugarLogger,
// 		SpanId:     spanID,
// 		Context:    Context,
// 	}
// 	return service
// }
