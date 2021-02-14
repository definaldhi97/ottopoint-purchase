package controllers

import (
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	hostAuth "ottopoint-purchase/hosts/auth/host"
	redisService "ottopoint-purchase/hosts/redis_token/host"
	signature "ottopoint-purchase/hosts/signature/host"
	"ottopoint-purchase/models"
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
		// var RedisService = new(services.RedisService)
		auth := strings.ReplaceAll(header.Authorization, "Bearer ", "")
		keyRedis := header.InstitutionID + "-" + auth
		dataRedis, _ := redisService.GetToken(keyRedis)

		if dataRedis.ResponseCode != "00" {
			result = utils.GetMessageFailedErrorNew(result, constants.RC_ERROR_INVALID_TOKEN, constants.RD_ERROR_INVALID_TOKEN)
			return header, result
		}

		phone = dataRedis.Data

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

func ValidatePerfix(CustID, ProductCode, category string) bool {
	// res := models.Response{Meta: utils.ResponseMetaOK()}
	fmt.Println("[Category : " + category + " ]")
	category1 := strings.ToLower(category)
	if category1 == constants.CategoryPulsa || category1 == constants.CategoryPaketData {
		// validate prefix
		fmt.Println("Process validasi verfix : ", category1)
		validate, _ := ValidatePrefixComulative(CustID, ProductCode, category1)
		if validate == false {

			fmt.Println("Invalid Prefix")
			// res = utils.GetMessageResponse(res, 500, false, errors.New("Nomor akun ini tidak terdafatr"))
			return false
		}

	}

	return true
}

func ValidatePrefixComulative(custID, productCode, category string) (bool, error) {

	var err error
	var product string
	var prefix string

	// validate panjang nomor, Jika nomor kurang dari 4
	if len(custID) < 4 {

		fmt.Println("[Kurang dari 4]-[Prefix-ottopoint]-[RedeemComulativeVoucher]")
		fmt.Println(fmt.Sprintf("invalid Prefix %v", custID))

		return false, err
	}

	// validate panjang nomor, Jika nomor kurang dari 11 & lebih dari 15
	if len(custID) <= 10 || len(custID) > 15 {

		fmt.Println("[Kurang dari 10 atau lebih dari 15]-[Prefix-ottopoint]-[RedeemComulativeVoucher]")
		fmt.Println(fmt.Sprintf("invalid Prefix %v", custID))

		return false, err

	}

	// get Prefix
	dataPrefix, errPrefix := db.GetOperatorCodebyPrefix(custID)
	if errPrefix != nil {

		fmt.Println("[ErrorPrefix]-[RedeemComulativeVoucher]")
		fmt.Println(fmt.Sprintf("dataPrefix = %v", dataPrefix))
		fmt.Println(fmt.Sprintf("Prefix tidak ditemukan %v", errPrefix))

		return false, err
	}

	// check operator by OperatorCode
	prefix = utils.Operator(dataPrefix.OperatorCode)
	// check operator by ProductCode
	// product = utils.ProductPulsa(productCode[0:4])

	if category == constants.CategoryPulsa {
		product = utils.ProductPulsa(productCode[0:4])
	}
	if category == constants.CategoryPaketData {
		product = utils.ProductPaketData(productCode[0:5])
	}

	// Jika Nomor tidak sesuai dengan operator
	if prefix != product {

		fmt.Println("[Operator]-[Prefix-ottopoint]-[RedeemComulativeVoucher]")
		fmt.Println(fmt.Sprintf("invalid Prefix %v", prefix))

		return false, err

	}

	return true, nil
}
