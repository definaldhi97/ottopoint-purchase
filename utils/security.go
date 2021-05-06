package utils

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"ottopoint-purchase/db"
	"ottopoint-purchase/models"
)

// OttoAGCreateSignature ..
func OttoAGCreateSignature(timestamp string, data interface{}, key string) string {
	jsonReq, _ := json.Marshal(data)
	bodyMsg := string(jsonReq)

	regx := regexp.MustCompile("[^a-zA-Z0-9{}:.,]")
	var bodySign = strings.ToLower(regx.ReplaceAllLiteralString(bodyMsg, "")) + "&" + timestamp + "&" + key

	signatureSystem := HashSha512(key, bodySign)

	fmt.Println("\n============= Request Identity ============= ")
	fmt.Println("Key : ", key)
	fmt.Println("Body from Requestor : ", bodyMsg)
	fmt.Println("Body to create Signature Sistem : ", bodySign)
	fmt.Println("Signature created  : ", signatureSystem)
	fmt.Println("\n ")

	return signatureSystem
}

// VoucherAggregatorSignature
func VoucherAggregatorSignature(timestamp string, data interface{}, key string) string {
	jsonReq, _ := json.Marshal(data)
	bodyMsg := string(jsonReq)

	regx := regexp.MustCompile("[^a-zA-Z0-9{}:.,]")
	var bodySign = strings.ToLower(regx.ReplaceAllLiteralString(bodyMsg, "")) + "&" + timestamp + "&" + key

	signatureSystem := HashSha512(key, bodySign)

	fmt.Println("\n============= Request Identity ============= ")
	fmt.Println("Key : ", key)
	fmt.Println("Body from Requestor : ", bodyMsg)
	fmt.Println("Body to create Signature Sistem : ", bodySign)
	fmt.Println("Signature created  : ", signatureSystem)
	fmt.Println("\n ")

	return signatureSystem
}

// CreateSignatureGeneral
func CreateSignatureGeneral(timestamp string, data interface{}, header models.RequestHeader, institutiontype int) string { // institutiontype (1 : Apikey, 2 : PubKey)

	jsonReq, _ := json.Marshal(data)
	bodyMsg := string(jsonReq)

	var key string

	keyIns, errKey := db.GetInstitutionKey(header.InstitutionID)
	if errKey != nil {

		fmt.Println("[CreateSignatureGeneral]-[GetInstitutionKey]")
		fmt.Println(fmt.Sprintf("[Failed Get Key %v]-[Error : %v]", header.InstitutionID, errKey))

		return ""

	}

	switch institutiontype {
	case 1:
		key = keyIns.Apikey
	case 2:
		key = keyIns.PubKey
	default:
		return ""
	}

	jsonRegString := SignReplaceAll(bodyMsg)
	plainText := fmt.Sprintf(jsonRegString + "&" + header.DeviceID + "&" + header.InstitutionID + "&" + header.Geolocation + "&" + header.ChannelID + "&" + header.AppsID + "&" + timestamp + "&" + key)
	fmt.Println("request data signature system : ", plainText)

	signatureSystem := HashSha512(key, plainText)

	fmt.Println("\n============= Request Identity ============= ")
	fmt.Println("Key : ", key)
	fmt.Println("Body from Requestor : ", bodyMsg)
	fmt.Println("Body to create Signature Sistem : ", plainText)
	fmt.Println("Signature created  : ", signatureSystem)
	fmt.Println("\n ")

	return signatureSystem
}

// HashSha512 ...
func HashSha512(secret, data string) string {
	hash := hmac.New(sha512.New, []byte(secret))
	hash.Write([]byte(data))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

type SignatureResp struct {
	Signature string `json:"signature"`
}

func SignReplaceAll(str string) string {
	regx := regexp.MustCompile("[^a-zA-Z0-9{}:.,]")
	output := regx.ReplaceAllLiteralString(str, "")
	output = strings.ToLower(output)
	return output
}
