package db

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego/logs"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	hcmodels "ottodigital.id/library/healthcheck/models"
	hcutils "ottodigital.id/library/healthcheck/utils"
	ODU "ottodigital.id/library/utils"
)

var (
	Dbcon   *gorm.DB
	DbErr   error
	DbUser  string
	DbPass  string
	DbName  string
	DbAddr  string
	DbPort  string
	DbDebug bool
	DbType  string
	SslMode string
)

func init() {

	DbType = ODU.GetEnv("DB_TYPE", "POSTGRES")
	fmt.Println(" DB POSTGRES")
	//orm.RegisterDataBase("default", "postgres", "postgres://postgres:cobain77@127.0.0.1:5432/ottoagsu?SslMode=disable")
	DbUser = ODU.GetEnv("DB_POSTGRES_USER", "ottoagcfg")
	DbPass = ODU.GetEnv("DB_POSTGRES_PASS", "dTj*&56$es")
	DbName = ODU.GetEnv("DB_POSTGRES_NAME", "ottopoint")
	DbAddr = ODU.GetEnv("DB_POSTGRES_HOST", "13.228.23.160")
	DbPort = ODU.GetEnv("DB_POSTGRES_PORT", "8432")
	SslMode = ODU.GetEnv("DB_POSTGRES_SSL_MODE", "disable")
	DbDebug = true //defaultValue
	if ndb := ODU.GetEnv("DB_POSTGRES_DEBUG", "true"); strings.EqualFold(ndb, "true") || strings.EqualFold(ndb, "false") {
		DbDebug, _ = strconv.ParseBool(ndb)
	}

	if DbOpen() != nil {
		//panic("DB Can't Open")
		fmt.Println("Can't open", DbName, "DB")
	}
	Dbcon = GetDbcon()
	Dbcon = Dbcon.LogMode(true)
}

// DbOpen ..
func DbOpen() error {
	args := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", DbAddr, DbPort, DbUser, DbPass, DbName, SslMode)
	Dbcon, DbErr = gorm.Open("postgres", args)

	if DbErr != nil {
		fmt.Println("Open", DbName, "DB error :", DbErr)
		return DbErr
	}

	if errping := Dbcon.DB().Ping(); errping != nil {
		return errping
	}
	return nil
}

// GetDbcon ..
func GetDbcon() *gorm.DB {
	//TODO looping try connection until timeout
	// using channel timeout
	if errping := Dbcon.DB().Ping(); errping != nil {
		logs.Error("DB not connected test ping :", errping)
		errping = nil
		if errping = DbOpen(); errping != nil {
			logs.Error("Try to connect again but error :", errping)
		}
	}
	Dbcon.LogMode(DbDebug)
	return Dbcon
}

// GetDatabaseHealthCheck ..
func GetDatabaseHealthCheck() hcmodels.DatabaseHealthCheck {
	dbCon := GetDbcon()
	return hcutils.GetDatabaseHealthCheck(&dbCon, &hcmodels.DatabaseEnv{
		Name: DbName,
		Host: DbAddr,
		Port: DbPort,
	})
}
