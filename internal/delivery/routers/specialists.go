package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/delivery/handlers"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/repository"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/services"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/config"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/database"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/log"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

func InitSpecialistsRouting(group *gin.RouterGroup, db *sqlx.DB, session database.Session, logger *log.Logs) {
	CasesPerRequest := viper.GetInt(config.CasesPerRequest)

	specialistRepo := repository.InitSpecialistsRepo(db)
	caseRepo := repository.InitCaseRepo(db, CasesPerRequest)

	specialistService := services.InitSpecialistService(specialistRepo, caseRepo, logger)
	specialistHandler := handlers.InitSpecialistsHandler(specialistService, session)

	group.GET("/me", specialistHandler.GetMe)
	group.PUT("/update", specialistHandler.UpdateMe)

	group.GET("/get_cases_by_level", specialistHandler.GetCasesByLevel)

	group.POST("/create_rated", specialistHandler.CreateRated)
	group.GET("/get_rated_solved", specialistHandler.GetRatedSolved)
}
