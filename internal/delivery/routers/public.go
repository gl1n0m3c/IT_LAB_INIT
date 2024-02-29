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
)

func InitPublicRouting(group *gin.RouterGroup, db *sqlx.DB, session database.Session, JWTUtil jwt.JWT, logger *log.Logs) {

	specialistRepo := repository.InitSpecialistsRepo(db)

	publicService := services.InitPublicService(specialistRepo, logger)
	publicHandler := handlers.InitPublicHandler(publicService, session, JWTUtil)

	group.POST("/specialist_register", publicHandler.SpecialistRegister)
	group.POST("/specialist_login", publicHandler.SpecialistLogin)
	group.POST("/refresh", publicHandler.Refresh)
}
