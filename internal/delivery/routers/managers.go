package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/delivery/handlers"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/repository"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/services"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/config"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/log"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

func InitManagersRouting(group *gin.RouterGroup, db *sqlx.DB, logger *log.Logs) {
	CasesPerRequest := viper.GetInt(config.CasesPerRequest)

	caseRepo := repository.InitCaseRepo(db, CasesPerRequest)

	managerService := services.InitManagerService(caseRepo, logger)
	managerHandler := handlers.InitManagerHandler(managerService)

	group.GET("/get_case/:case_id", managerHandler.GetFulCaseByID)
}
