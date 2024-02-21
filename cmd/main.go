package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/config"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/database"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/log"
)

func main() {
	router := gin.Default()

	logger, loggerInfoFile, loggerErrorFile := log.InitLoggers()
	defer loggerInfoFile.Close()
	defer loggerErrorFile.Close()
	logger.InfoLogger.Info().Msg("Logger Initialized")

	config.InitConfig()
	logger.InfoLogger.Info().Msg("Config Initialized")

	db := database.GetDB()
	logger.InfoLogger.Info().Msg("Database Initialized")

	_ = db // убрать

	if err := router.Run("0.0.0.0:8080"); err != nil {
		panic(fmt.Sprintf("Failed to run client: %s", err.Error()))
	}
}
