package tests

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

const (
	TestDBName     = "TEST_DB_NAME"
	TestDBUser     = "TEST_DB_USER"
	TestDBPassword = "TEST_DB_PASSWORD"
	TestDBPort     = "TEST_DB_PORT"
	TestDBHost     = "TEST_DB_HOST"
)

func InitTestConfig() {
	envPath, _ := os.Getwd()
	envPath = filepath.Join(envPath, "..")

	viper.SetConfigName("tests")
	viper.SetConfigType("env")
	viper.AddConfigPath(envPath)
	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Sprintf("Failed to init config: %v", err.Error()))
	}
}
