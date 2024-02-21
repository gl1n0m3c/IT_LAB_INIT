package database

import (
	"database/sql"
	"fmt"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/config"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
)

func GetDB() *sql.DB {
	connectionString := fmt.Sprintf(
		"user=%v password=%v host=%v port=%v dbname=%v",
		viper.GetString(config.DBUser),
		viper.GetString(config.DBPassword),
		viper.GetString(config.DBHost),
		viper.GetInt(config.DBPort),
		viper.GetString(config.DBName),
	)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to DB: %s", err.Error()))
	}

	return db
}
