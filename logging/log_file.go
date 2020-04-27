package logging

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	logging "github.com/hhkbp2/go-logging"
	ODU "ottodigital.id/library/utils"
)

var (
	log_file logging.Logger
)

func init() {

	logs.Info("log_file.go")

	config_file := ODU.GetEnv("log.path.cofig", "./log/config_log.yml")

	if err := logging.ApplyConfigFile(config_file); err != nil {
		fmt.Println("Error : ", err)
		//log.Debug("Error init config ",err, config_file)
		//log.Error("Error init config",err)
	}
	log_file = logging.GetLogger("metadata")

	logconfig := Logconfig{}
	logconfig.Filename = beego.AppConfig.DefaultString("log.filename", "./log/ottopoint-purchase.log")
	logconfig.Daily = beego.AppConfig.DefaultBool("log.daily", true)
	logconfig.MaxLines = beego.AppConfig.DefaultInt("log.maxlines", 1000000)
	logconfig.Level = beego.AppConfig.DefaultInt("log.level", 7)
	logconfig.MaxDays = beego.AppConfig.DefaultInt64("log.maxdays", 1)
	logconfig.MaxSize = beego.AppConfig.DefaultInt("log.maxsize", 1<<28)
	logconfig.Perm = beego.AppConfig.DefaultString("log.perm", "0666")
	logconfig.Rotate = beego.AppConfig.DefaultBool("log.rotate", true)
	cfglog, err := GetConfig(logconfig)
	if err != nil {
		fmt.Println("Config Log Error", err)
		logs.Error("Config Log Error ", err)
	}
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	logs.SetLogger(logs.AdapterFile, cfglog)
	logs.Async()
}

func GetLog() logging.Logger {
	return log_file
}

type Logconfig struct {
	//File name
	Filename string `json:"filename"`

	// Rotate at line
	MaxLines int `json:"maxlines"`

	// Rotate at size
	MaxSize int `json:"maxsize"`

	// Rotate daily  If log rotates by day, true by default.
	Daily bool `json:"daily"`

	// maxdays: Maximum number of days log files will be kept, 7 by default.
	MaxDays int64 `json:"maxdays"`

	//rotate: Enable logrotate or not, true by default.
	Rotate bool `json:"rotate"`

	//level: Log level, Trace by default = 7 .
	Level int `json:"level"`

	//perm: Log file permission default set 0666
	Perm string `json:"perm"`
}

func GetConfig(configlog interface{}) (string, error) {
	var configstr []byte
	var err error
	if configstr, err = json.Marshal(configlog); err != nil {
		logs.Error("Config Error [%v]", err)
		return string(configstr), err
	}
	return string(configstr), nil
}

func LogErrors(errorMessage string, errors ...error) {
	for _, err := range errors {
		if err != nil {
			logs.Warn(errorMessage, err)
		}
	}
}
