package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/delivery/docs"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/delivery/middleware"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/delivery/routers"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/reporting_period"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/config"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/database"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/log"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/tracing"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/jwt"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const serviceName = "admin-panel"

func main() {
	router := gin.Default()

	router.Static("/static", "../static")

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

	middleWarrior := middleware.InitMiddleware(JWTUtil, logger)
	logger.InfoLogger.Info().Msg("Middleware Initialized")

	jaegerURL := fmt.Sprintf("http://%v:%v/api/traces", viper.GetString(config.JaegerHost), viper.GetString(config.JaegerPort))
	tracer := tracing.InitTracer(jaegerURL, serviceName)

	routers.InitRouting(router, db, session, JWTUtil, middleWarrior, logger, tracer)
	logger.InfoLogger.Info().Msg("Routing Initialized")

	// Для загрузки тестовых данных
	//utils.LoadFixtures(db)
	//utils.ClearDatabase(db)

	go reporting_period.StartReporting(db, logger)

	if err := router.Run("0.0.0.0:8080"); err != nil {
		panic(fmt.Sprintf("Failed to run client: %s", err.Error()))
	}
}
