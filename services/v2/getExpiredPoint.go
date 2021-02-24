package services

import (
	"fmt"
	"ottopoint-purchase/hosts/opl/host"
	"ottopoint-purchase/utils"
	"time"
)

func ExpiredPointService() string {

	fmt.Println(">>> ExpiredPointService <<<")

	get, err := host.SettingsOPL()
	if err != nil || get.Settings.ProgramName == "" {

		fmt.Println(fmt.Sprintf("[Error : %v]", err))
		fmt.Println("[PackageBulkService]-[ExpiredPointService]")

	}

	data := get.Settings.PointsDaysActiveCount + 1

	expired := utils.FormatTimeString(time.Now(), 0, 0, data)

	return expired

} 
