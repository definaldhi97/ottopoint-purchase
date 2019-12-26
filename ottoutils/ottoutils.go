package ottoutils

import (
	"fmt"
	"os"
	"plugin"
	"crypto/sha256"
	"encoding/hex"
)

func GetSymbol(namefileso string, symbolname string) (plugin.Symbol,error)   {
	p, err := plugin.Open(namefileso)
	if err != nil {
		fmt.Printf("Error Open: %s\n", err)
		return nil ,err
	}
	h, err := p.Lookup(symbolname)
	if err != nil {
		fmt.Printf("Error Lookup: %s\n", err)
		return nil ,err
	}

	return h , err
}


// Simple helper function to read an environment or return a default value
func GetEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}


func SHA(str string) string {

	bytes := []byte(str)

	// Converts string to sha2
	h := sha256.New()                   // new sha256 object
	h.Write(bytes)                      // data is now converted to hex
	code := h.Sum(nil)                  // code is now the hex sum
	codestr := hex.EncodeToString(code) // converts hex to string

	return codestr
}
