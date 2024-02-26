package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/delivery/docs"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/delivery/routers"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/config"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/database"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/log"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/jwt"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	session := database.InitRedisSession()
	logger.InfoLogger.Info().Msg("Session Initialized")

	JWTUtil := jwt.InitJWTUtil()
	logger.InfoLogger.Info().Msg("JWTUtil Initialized")

	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	logger.InfoLogger.Info().Msg("Swagger Initialized")

	routers.InitRouting(router, db, session, JWTUtil, logger)
	logger.InfoLogger.Info().Msg("Routing Initialized")

	if err := router.Run("0.0.0.0:8080"); err != nil {
		panic(fmt.Sprintf("Failed to run client: %s", err.Error()))
	}
}
