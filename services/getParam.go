package services

import (
	"fmt"
	"ottopoint-purchase/constants"
	"ottopoint-purchase/db"
	"ottopoint-purchase/models"
)
 
func DataParameterOrder() models.ParamUV {
	res := models.ParamUV{}

	// nama := "" // nama
	// email := "UV_EMAIL_ORDER"
	// phone := "UV_PHONE_ORDER"
	// expired := "UV_EXPIRED_VOUCHER"
	// group := "UVCONFIG"

	datanama, errnama := db.ParamData(constants.CODE_CONFIG_UV_GROUP, constants.CODE_CONFIG_UV_NAME)
	if errnama != nil {
		fmt.Println("[Error get data from Db m_paramaters]")
		fmt.Println("Error : ", errnama)
		fmt.Println("Code :", constants.CODE_CONFIG_UV_NAME)
	}

	dataemail, erremail := db.ParamData(constants.CODE_CONFIG_UV_GROUP, constants.CODE_CONFIG_UV_EMAIL)
	if erremail != nil {
		fmt.Println("[Error get data from Db m_paramaters]")
		fmt.Println("Error : ", erremail)
		fmt.Println("Code :", constants.CODE_CONFIG_UV_EMAIL)
	}

	dataphone, errphone := db.ParamData(constants.CODE_CONFIG_UV_GROUP, constants.CODE_CONFIG_UV_PHONE)
	if errphone != nil {
		fmt.Println("[Error get data from Db m_paramaters]")
		fmt.Println("Error : ", errphone)
		fmt.Println("Code :", constants.CODE_CONFIG_UV_PHONE)
	}

	dataexpired, errexpired := db.ParamData(constants.CODE_CONFIG_UV_GROUP, constants.CODE_CONFIG_UV_EXPIRED)
	if errexpired != nil {
		fmt.Println("[Error get data from Db m_paramaters]")
		fmt.Println("Error : ", errexpired)
		fmt.Println("Code :", constants.CODE_CONFIG_UV_EXPIRED)
	}

	res = models.ParamUV{
		Nama:    datanama.Value,
		Email:   dataemail.Value,
		Phone:   dataphone.Value,
		Expired: dataexpired.Value,
	}

	return res
}
