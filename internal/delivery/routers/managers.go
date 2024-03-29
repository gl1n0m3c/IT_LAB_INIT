package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/delivery/handlers"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/repository"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/services"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/log"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
)

func InitManagersRouting(group *gin.RouterGroup, db *sqlx.DB, logger *log.Logs, tracer trace.Tracer) {
	specialistsRepo := repository.InitSpecialistsRepo(db)
	caseRepo := repository.InitCaseRepo(db)

	managerService := services.InitManagerService(caseRepo, specialistsRepo, logger)
	managerHandler := handlers.InitManagerHandler(managerService, tracer)

	group.GET("/get_case", managerHandler.GetFulCaseByID)
	group.GET("/get_specialists_rating", managerHandler.GetSpecialistRating)
}
