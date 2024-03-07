package middleware

import (
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/log"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/jwt"
)

type Middleware struct {
	jwtUtil jwt.JWT
	logger  *log.Logs
}

func InitMiddleware(
	JWTUtil jwt.JWT,
	logger *log.Logs,
) Middleware {
	return Middleware{
		jwtUtil: JWTUtil,
		logger:  logger,
	}
}
