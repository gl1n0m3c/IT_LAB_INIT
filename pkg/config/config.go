package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

const (
	DBName     = "DB_NAME"
	DBUser     = "DB_USER"
	DBPassword = "DB_PASSWORD"
	DBPort     = "DB_PORT"
	DBHost     = "DB_HOST"
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
