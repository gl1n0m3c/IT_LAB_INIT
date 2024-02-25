package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
	_ "github.com/gl1n0m3c/IT_LAB_INIT/internal/models/swagger"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/services"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/responses"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type publicHandler struct {
	service services.Public
}

func InitPublicHandler(
	service services.Public,
) publicHandler {
	return publicHandler{
		service: service,
	}
}

// SpecialistRegister registers a new specialist and returns a JWT and refresh token upon successful registration.
// @Summary Specialist Registration
// @Description Registers a new specialist and returns a JWT and refresh token upon successful registration.
// @Description Automatically level=1, is_verified=false.
// @Description Login and password are required. There are no validation on password, but they could be.
// @Tags public
// @Accept json
// @Produce json
// @Param specialist body swagger.SpecialistCreate true "Specialist Registration"
// @Success 201 {object} responses.MessageDataResponse "Successful registration, returning JWT and refresh token"
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
	if err := validate.Struct(specialist); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(fmt.Sprintf(responses.Response400, err)))
		return
	}

	ID, err := p.service.SpecialistRegister(ctx, specialist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
		return
	}

	c.JSON(http.StatusCreated, responses.NewMessageDataResponse(fmt.Sprintf(responses.Response201, "specialist", ID), ID))
}
