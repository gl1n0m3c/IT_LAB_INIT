package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/responses"
	"net/http"
	"strings"
)

const (
	UserID = "userID"
)

func (m Middleware) Authorization(userType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		parts := strings.Split(auth, " ")

		if !strings.Contains(auth, "Bearer ") || len(parts) < 1 {
			m.logger.InfoLogger.Info().Msg(fmt.Sprintf("No jwt provided at: %v", c.Request.URL.Path))
			c.AbortWithStatusJSON(http.StatusUnauthorized, responses.NewMessageResponse("No bearer provided in authorization"))
			return
		}

		jwtToken := parts[1]

		userData, isValid, err := m.jwtUtil.Authorize(jwtToken, userType)
		if err != nil {
			m.logger.InfoLogger.Error().Msg(fmt.Sprintf("Troubles while getting user info from jwt: %v", err))
			c.AbortWithStatusJSON(http.StatusBadRequest, responses.NewMessageResponse("Bad jwt provided"))
			return
		}

		if !isValid {
			m.logger.InfoLogger.Error().Msg("Access token is expired or invalid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, responses.NewMessageResponse("Access token is expired or invalid"))
			return
		}

		c.Set(UserID, userData.ID)
	}
}
