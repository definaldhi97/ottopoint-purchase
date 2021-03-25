package services

import (
	"ottopoint-purchase/db"
	"ottopoint-purchase/models"

	"github.com/sirupsen/logrus"
)

func DataParameterOrder(codeGroup, codeSupplier, codeEmail, codePhone, codeExpired string) models.ParamOrder {
	res := models.ParamOrder{}

	datanama, errnama := db.ParamData(codeGroup, codeSupplier)
	if errnama != nil {
		logrus.Error("[Error get data from Db m_paramaters]")
		logrus.Error("Error : ", errnama)
		logrus.Error("Code :", codeSupplier)
	}

	dataemail, erremail := db.ParamData(codeGroup, codeEmail)
	if erremail != nil {
		logrus.Error("[Error get data from Db m_paramaters]")
		logrus.Error("Error : ", erremail)
		logrus.Error("Code :", codeEmail)
	}

	dataphone, errphone := db.ParamData(codeGroup, codePhone)
	if errphone != nil {
		logrus.Error("[Error get data from Db m_paramaters]")
		logrus.Error("Error : ", errphone)
		logrus.Error("Code :", codePhone)
	}

	dataexpired, errexpired := db.ParamData(codeGroup, codeExpired)
	if errexpired != nil {
		logrus.Error("[Error get data from Db m_paramaters]")
		logrus.Error("Error : ", errexpired)
		logrus.Error("Code :", codeExpired)
	}

	res = models.ParamOrder{
		Nama:    datanama.Value,
		Email:   dataemail.Value,
		Phone:   dataphone.Value,
		Expired: dataexpired.Value,
	}

	return res
}
