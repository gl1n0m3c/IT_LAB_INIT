package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/log"
	"github.com/jmoiron/sqlx"
)

func InitRouting(r *gin.Engine, db *sqlx.DB, logger *log.Logs) {
	publicGroup := r.Group("/public")

	InitPublicRouting(publicGroup, db, logger)
}
