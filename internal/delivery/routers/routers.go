package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/database"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/log"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/jwt"
	"github.com/jmoiron/sqlx"
)

func InitRouting(r *gin.Engine, db *sqlx.DB, session database.Session, JWTUtil jwt.JWT, logger *log.Logs) {
	publicGroup := r.Group("/public")

	InitPublicRouting(publicGroup, db, session, JWTUtil, logger)
}
