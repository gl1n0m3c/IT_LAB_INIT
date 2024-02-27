package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
	_ "github.com/gl1n0m3c/IT_LAB_INIT/internal/models/swagger"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/services"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/database"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/jwt"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/responses"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/validators"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type publicHandler struct {
	service services.Public
	session database.Session
	JWTUtil jwt.JWT
}

func InitPublicHandler(
	service services.Public,
	session database.Session,
	JWTUtil jwt.JWT,
) publicHandler {
	return publicHandler{
		service: service,
		session: session,
		JWTUtil: JWTUtil,
	}
}

// SpecialistRegister registers a new specialist and returns a jwt and refresh token upon successful registration.
// @Summary Specialist Registration
// @Description Registers a new specialist and returns a jwt and refresh token upon successful registration.
// @Description Automatically level=1, is_verified=false.
// @Description Login and password are required. There are some validation on password:
// @Description More than 8 symbols, contain at least one number, one uppercase and one lowercase letter.
// @Tags public
// @Accept json
// @Produce json
// @Param specialist body swagger.SpecialistCreate true "Specialist Registration"
// @Success 201 {object} responses.JWTRefresh "Successful registration, returning jwt and refresh token"
// @Failure 400 {object} responses.MessageResponse "Invalid input"
// @Failure 500 {object} responses.MessageResponse "Internal server error"
// @Router /public/specialist_register [post]
func (p publicHandler) SpecialistRegister(c *gin.Context) {
	var specialist models.SpecialistCreate

	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&specialist); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(fmt.Sprintf(responses.Response400, err)))
		return
	}

	validate := validator.New()
	_ = validate.RegisterValidation("password", validators.ValidatePassword)
	if err := validate.Struct(specialist); err != nil {
		customErrMsg := validators.CustomErrorMessage(err)
		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(customErrMsg))
		return
	}

	ID, err := p.service.SpecialistRegister(ctx, specialist)

	if err != nil {
		if errors.Is(err, utils.UniqueSpecialistErr) {
			c.JSON(http.StatusBadRequest, responses.NewMessageResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
		return
	}

	accessToken := p.JWTUtil.CreateToken(ID, jwt.Specialist)

	refreshToken, err := p.session.Set(ctx, database.SessionData{
		UserID:   ID,
		UserType: jwt.Specialist,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
		return
	}

	c.JSON(http.StatusCreated, responses.NewJWTRefreshResponse(accessToken, refreshToken))
}

// Refresh updates access and refresh tokens
// @Summary Refresh Tokens
// @Description Refreshes access and refresh tokens using a refresh token provided in the Authorization header.
// @Tags public
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Refresh header string true "Refresh Token"
// @Success 200 {object} responses.JWTRefresh "Successful token refresh, returning new jwt and refresh token"
// @Failure 400 {object} responses.MessageResponse "No refresh token provided"
// @Failure 401 {object} responses.MessageResponse "Invalid or expired refresh token"
// @Failure 500 {object} responses.MessageResponse "Internal server error"
// @Router /public/refresh [post]
func (p publicHandler) Refresh(c *gin.Context) {
	oldRefreshToken := c.GetHeader("Refresh")
	if oldRefreshToken == "" {
		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(fmt.Sprintf(responses.Response400, "Отсутсвует Refresh Token")))
		return
	}

	ctx := c.Request.Context()

	newRefreshToken, userData, err := p.session.GetAndUpdate(ctx, oldRefreshToken)
	if err != nil {
		switch err {
		case utils.NeedToAuthorize:
			c.JSON(http.StatusUnauthorized, responses.NewMessageResponse(err.Error()))
			return
		default:
			c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(err.Error()))
			return
		}
	}

	newAccessToken := p.JWTUtil.CreateToken(userData.UserID, userData.UserType)

	c.JSON(http.StatusOK, responses.NewJWTRefreshResponse(newAccessToken, newRefreshToken))
}

//{
//"JWT": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDkwNzA1MTEsIklEIjoxOCwiVXNlclR5cGUiOiJzcGVjaWFsaXN0In0.FaIRZxczoUVWWVNKTIHyYkU9qqLH26h-z7nVF3RJMgM",
//"RefreshToken": "d9ef7978-a1c5-4534-a176-f6905cf8c904"
//}
