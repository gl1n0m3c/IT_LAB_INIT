package database

import (
	"fmt"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/config"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
)

func GetDB() *sqlx.DB {
	connectionString := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		viper.GetString(config.DBUser),
		viper.GetString(config.DBPassword),
		viper.GetString(config.DBHost),
		viper.GetInt(config.DBPort),
		viper.GetString(config.DBName),
	)

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to DB: %s", err.Error()))
	}

	return db
}
