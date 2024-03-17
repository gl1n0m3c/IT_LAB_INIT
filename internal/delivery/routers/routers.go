package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/delivery/middleware"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/database"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/log"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/jwt"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/trace"
)

func InitRouting(r *gin.Engine, db *sqlx.DB, session database.Session, JWTUtil jwt.JWT, middleware middleware.Middleware, logger *log.Logs, tracer trace.Tracer) {
	managerGroup := r.Group("/manager")
	publicGroup := r.Group("/public")
	specialistsGroup := r.Group("/specialist")

	specialistsGroup.Use(middleware.Authorization(jwt.Specialist))
	managerGroup.Use(middleware.Authorization(jwt.Manager))

	InitManagersRouting(managerGroup, db, logger, tracer)
	InitPublicRouting(publicGroup, db, session, JWTUtil, logger, tracer)
	InitSpecialistsRouting(specialistsGroup, db, session, logger, tracer)
}
