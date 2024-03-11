package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/services"
	customErrors "github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/custom_errors"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/responses"
	"net/http"
	"strconv"
)

type managerHandler struct {
	service services.Managers
}

func InitManagerHandler(
	service services.Managers,
) Managers {
	return managerHandler{
		service: service,
	}
}

// GetFulCaseByID @Summary Retrieve a case by ID
// @Description Retrieves a case by its ID and returns detailed information about the case.
// @Tags managers
// @Accept  json
// @Produce  json
// @Param authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param case_id path int true "ID of the case to retrieve"
// @Success 200 {object} models.CaseFul "Successfully retrieved the case details"
// @Failure 400 {object} responses.MessageResponse "Invalid query parameter or missing case_id"
// @Failure 403 {object} responses.MessageResponse "JWT is invalid or expired"
// @Failure 500 {object} responses.MessageResponse "Internal server error"
// @Router /manager/get_case/{case_id} [get]
func (m managerHandler) GetFulCaseByID(c *gin.Context) {
	ctx := c.Request.Context()

	caseIDStr := c.Param("case_id")
	caseID, err := strconv.Atoi(caseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(responses.ResponseBadQuery))
		return
	}

	caseData, err := m.service.GetFulCaseByID(ctx, caseID)
	if err != nil {
		if errors.Is(err, customErrors.NoRowsCaseErr) {
			c.JSON(http.StatusNotFound, responses.NewMessageResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
		return
	}

	c.JSON(http.StatusOK, caseData)
}
