package services

import (
	"fmt"
	"ottopoint-purchase/db"
	"ottopoint-purchase/models"
)

func DataParameterOrder(codeGroup, codeSupplier, codeEmail, codePhone, codeExpired string) models.ParamOrder {
	res := models.ParamOrder{}

	datanama, errnama := db.ParamData(codeGroup, codeSupplier)
	if errnama != nil {
		fmt.Println("[Error get data from Db m_paramaters]")
		fmt.Println("Error : ", errnama)
		fmt.Println("Code :", codeSupplier)
	}

	dataemail, erremail := db.ParamData(codeGroup, codeEmail)
	if erremail != nil {
		fmt.Println("[Error get data from Db m_paramaters]")
		fmt.Println("Error : ", erremail)
		fmt.Println("Code :", codeEmail)
	}

	dataphone, errphone := db.ParamData(codeGroup, codePhone)
	if errphone != nil {
		fmt.Println("[Error get data from Db m_paramaters]")
		fmt.Println("Error : ", errphone)
		fmt.Println("Code :", codePhone)
	}

	dataexpired, errexpired := db.ParamData(codeGroup, codeExpired)
	if errexpired != nil {
		fmt.Println("[Error get data from Db m_paramaters]")
		fmt.Println("Error : ", errexpired)
		fmt.Println("Code :", codeExpired)
	}

	res = models.ParamOrder{
		Nama:    datanama.Value,
		Email:   dataemail.Value,
		Phone:   dataphone.Value,
		Expired: dataexpired.Value,
	}

	return res
}
