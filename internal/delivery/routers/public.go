package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/delivery/handlers"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/repository"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/services"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/database"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/log"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/jwt"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
)

func InitPublicRouting(group *gin.RouterGroup, db *sqlx.DB, session database.Session, JWTUtil jwt.JWT, logger *log.Logs, tracer trace.Tracer) {
	managerRepo := repository.InitManagerRepo(db)
	specialistRepo := repository.InitSpecialistsRepo(db)
	cameraRepo := repository.InitCameraRepo(db)
	caseRepo := repository.InitCaseRepo(db)

	publicService := services.InitPublicService(managerRepo, specialistRepo, cameraRepo, caseRepo, logger)
	publicHandler := handlers.InitPublicHandler(publicService, session, JWTUtil, tracer)

	group.POST("/manager_login", publicHandler.ManagerLogin)

	group.POST("/specialist_register", publicHandler.SpecialistRegister)
	group.POST("/specialist_login", publicHandler.SpecialistLogin)

	group.POST("/camera_create", publicHandler.CameraCreate)
	group.DELETE("/camera_delete", publicHandler.CameraDelete)

	group.POST("/case_create", publicHandler.CaseCreate)

	group.POST("/refresh", publicHandler.Refresh)
}
