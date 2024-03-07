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

	caseRepo := repository.InitCaseRepo(db, CasesPerRequest)

	specialistService := services.InitSpecialistService(caseRepo, logger)
	specialistHandler := handlers.InitSpecialistsHandler(specialistService, session)

	group.POST("/create_rated", specialistHandler.CreateRated)
	group.GET("/get_rated_solved", specialistHandler.GetRatedSolved)
	group.PUT("/update_rated_status", specialistHandler.UpdateRatedStatus)
}
