package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

const (
	K               = "K"
	J               = "J"
	ReportingPeriod = "REPORTING_PERIOD"

	DBName     = "DB_NAME"
	DBUser     = "DB_USER"
	DBPassword = "DB_PASSWORD"
	DBPort     = "DB_PORT"
	DBHost     = "DB_HOST"

	DBResponseTime = "DB_RESPONSE_TIME"

	SessionPassword = "SESSION_PASSWORD"
	SessionHost     = "SESSION_HOST"
	SessionPort     = "SESSION_PORT"
	SessionSaveTime = "SESSION_SAVE_TIME"

	JWTExpire = "JWT_EXPIRE"
	JWTSecret = "JWT_SECRET"

	EntitiesPerRequest = "ENTITIES_PER_REQUEST"

	Mail         = "MAIL"
	MailPassword = "MAIL_PASSWORD"
	MailHost     = "MAIL_HOST"
	MailPort     = "MAIL_PORT"

	JaegerHost = "JAEGER_HOST"
	JaegerPort = "JAEGER_PORT"
)

func InitConfig() {
	envPath, _ := os.Getwd()
	envPath = filepath.Join(envPath, "..") // workdir is cmd

	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(envPath)
	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Sprintf("Failed to init config: %v", err.Error()))
	}
}
