package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/kelseyhightower/envconfig"
	hcmodels "ottodigital.id/library/healthcheck/models"
)

type EnvPostgresDb struct {
	User    string `envconfig:"POSTGRES_USER" default:"ottoagcfg"`
	Pass    string `envconfig:"POSTGRES_PASS" default:"dTj*&56$es"`
	Name    string `envconfig:"POSTGRES_NAME" default:"ottopoint"`
	Host    string `envconfig:"POSTGRES_HOST" default:"13.228.23.160"`
	Port    string `envconfig:"POSTGRES_PORT" default:"8432"`
	Debug   bool   `envconfig:"POSTGRES_DEBUG" default:"true"`
	Type    string `envconfig:"TYPE" default:"postgres"`
	SslMode string `envconfig:"POSTGRES_SSL_MODE" default:"disable"`
}

var (
	DbCon         *gorm.DB
	DbErr         error
	envPostgresDb EnvPostgresDb
)

func init() {

	fmt.Println("DB POSTGRES")

	err := envconfig.Process("OTTOPOINT", &envPostgresDb)
	if err != nil {
		fmt.Println("Failed to get OTTOPOINT env:", err)
	}

	if DbOpen() != nil {
		//panic("DB Can't Open")
		fmt.Println("Can't open", envPostgresDb.Name, "DB")
	}
	DbCon = GetDbCon()
	DbCon = DbCon.LogMode(true)

}

// DbOpen ..
func DbOpen() error {
	args := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", envPostgresDb.Host, envPostgresDb.Port, envPostgresDb.User, envPostgresDb.Pass, envPostgresDb.Name, envPostgresDb.SslMode)
	DbCon, DbErr = gorm.Open("postgres", args)

	if DbErr != nil {
		fmt.Println("Open", envPostgresDb.Name, "DB error :", DbErr)
		return DbErr
	}

	if errping := DbCon.DB().Ping(); errping != nil {
		return errping
	}
	return nil
}

// GetDbCon ..
// GetDbCon ..
func GetDbCon() *gorm.DB {
	//TODO looping try connection until timeout
	// using channel timeout
	if errping := DbCon.DB().Ping(); errping != nil {
		fmt.Println("DB not connected test ping :", errping)
		errping = nil
		if errping = DbOpen(); errping != nil {
			fmt.Println("Try to connect again but error :", errping)
		}
	}
	DbCon.LogMode(envPostgresDb.Debug)
	return DbCon
}

func GetHealthCheck() hcmodels.DatabaseHealthCheck {
	return hcmodels.DatabaseHealthCheck{
		Name:   envPostgresDb.Name,
		Host:   envPostgresDb.Host,
		Port:   envPostgresDb.Port,
		Status: "OK",
		// Description: ,
	}
}
