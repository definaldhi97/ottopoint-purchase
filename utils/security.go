package utils

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
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

// HashSha512 ...
func HashSha512(secret, data string) string {
	hash := hmac.New(sha512.New, []byte(secret))
	hash.Write([]byte(data))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
