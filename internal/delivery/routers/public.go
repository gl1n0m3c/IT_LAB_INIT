package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/delivery/handlers"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/repository"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/services"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/log"
	"github.com/jmoiron/sqlx"
)

func InitPublicRouting(group *gin.RouterGroup, db *sqlx.DB, logger *log.Logs) {

	specialistRepo := repository.InitSpecialistsRepo(db)

	publicService := services.InitPublicService(specialistRepo, logger)
	publicHandler := handlers.InitPublicHandler(publicService)

	group.POST("/specialist_register", publicHandler.SpecialistRegister)
}
