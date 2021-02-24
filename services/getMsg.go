package services

import (
	"fmt"
	"ottopoint-purchase/db"
)

func GetMsgCummulative(rc, msg string) string {
 
	var codeMsg string

	getmsg, errmsg := db.GetResponseCummulativeOttoAG(rc)
	if errmsg != nil || getmsg.InternalRc == "" {

		fmt.Println("[VoucherComulativeService]-[GetResponseCummulativeOttoAG]")
		fmt.Println("[Failed to Get Data Mapping Response]")
		fmt.Println(fmt.Sprintf("[Data GetResponseOttoag : ]", getmsg))
		fmt.Println(fmt.Sprintf("[Error %v]", errmsg))
		// return res, err

		codeMsg = msg

		return codeMsg
	}

	// codeRc = getmsg.InternalRc
	// codeMsg = strings.Replace(getmsg.InternalRd, "[x]", "%v", 10)
	codeMsg = getmsg.InternalRd

	return codeMsg
}
