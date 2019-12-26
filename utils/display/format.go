package display

import (
	"fmt"
)

// GetFormattedToken ...
func GetFormattedToken(token string) string {

	tokenLen := len(token)

	if tokenLen == 0 {
		return token
	}

	formattedToken := ""

	for index := 0; index < tokenLen; index++ {
		if (index % 4) == 0 && index != 0 && index != tokenLen -1 {
			formattedToken = fmt.Sprintf("%s %s", formattedToken, token[index:index+1])

			continue
		}

		formattedToken = fmt.Sprintf("%s%s", formattedToken, token[index:index+1])
	}


	return formattedToken
}