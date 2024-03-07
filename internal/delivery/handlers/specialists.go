package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/services"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/database"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/responses"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/validators"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
	"time"
)

type specialistsHandler struct {
	service services.Specialists
	session database.Session
}

func InitSpecialistsHandler(
	service services.Specialists,
	session database.Session,
) Specialists {
	return specialistsHandler{
		service: service,
		session: session,
	}
}

// CreateRated @Summary Create a new rating
// @Description Creates a new rating entry based on the provided data.
// @Tags specialists
// @Accept  json
// @Produce  json
// @Param rated_data body models.RatedCreate true "Rated data"
// @Param authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 201 {object} responses.CreationIntResponse "Successfully created the rating"
// @Failure 400 {object} map[string]string "Invalid input data"
// @Failure 403 {object} map[string]string "JWT is invalid or expired"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /specialist/create_rated [post]
func (s specialistsHandler) CreateRated(c *gin.Context) {
	var (
		ratedReq models.RatedCreate
		rated    models.RatedBase
	)

	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&ratedReq); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(fmt.Sprintf(responses.Response400, err)))
		return
	}

	validate := validator.New()
	if err := validate.Struct(ratedReq); err != nil {
		customErrMsg := validators.CustomErrorMessage(err)
		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(customErrMsg))
		return
	}

	rated.RatedCreate = ratedReq
	rated.SpecialistID = c.GetInt("userID")
	rated.Date = time.Now().UTC()
	rated.Status = "Unknown"

	createdRatedID, err := s.service.CreateRated(ctx, rated)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
		return
	}

	c.JSON(http.StatusCreated, responses.CreationIntResponse{ID: createdRatedID})
}

// GetRatedSolved @Summary Get rated solved
// @Description Retrieves a rated solved entry based on the provided cursor ID.
// @Tags specialists
// @Accept  json
// @Produce  json
// @Param cursor_id query int true "Cursor ID for pagination"
// @Param authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} models.RatedCursor "Successfully retrieved the rated solved"
// @Failure 400 {object} map[string]string "Invalid input data"
// @Failure 401 {object} map[string]string "JWT is invalid or expired"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /specialist/get_rated_solved [get]
func (s specialistsHandler) GetRatedSolved(c *gin.Context) {
	ctx := c.Request.Context()

	cursorStr, ok := c.GetQuery("cursor_id")
	if !ok {
		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(responses.ResponseBadQuery))
		return
	}

	cursor, err := strconv.Atoi(cursorStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(responses.ResponseBadQuery))
		return
	}

	ratedCursor, err := s.service.GetRatedSolved(ctx, cursor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
		return
	}

	c.JSON(http.StatusCreated, ratedCursor)
}

// UpdateRatedStatus @Summary Update rated status
// @Description Updates the status of a rated entry based on the provided Rated ID and sets it to the new status.
// @Tags specialists
// @Accept  json
// @Produce  json
// @Param rated_update body models.RatedUpdate true "Rated update information"
// @Param authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 "Successfully updated the rated status"
// @Failure 400 {object} map[string]string "Invalid input data"
// @Failure 401 {object} map[string]string "JWT is invalid or expired"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /specialist/update_rated_status [put]
func (s specialistsHandler) UpdateRatedStatus(c *gin.Context) {
	var ratedUpdate models.RatedUpdate

	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&ratedUpdate); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(fmt.Sprintf(responses.Response400, err)))
		return
	}

	validate := validator.New()
	if err := validate.Struct(ratedUpdate); err != nil {
		customErrMsg := validators.CustomErrorMessage(err)
		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(customErrMsg))
		return
	}

	validStatuses := map[string]bool{
		"Correct":   true,
		"Incorrect": true,
		"Unknown":   true,
	}

	_, valid := validStatuses[ratedUpdate.Status]
	fmt.Println(valid)
	if _, valid := validStatuses[ratedUpdate.Status]; !valid {
		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(responses.ResponseBadBody))
		return
	}

	err := s.service.UpdateRatedStatus(ctx, ratedUpdate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
		return
	}

	c.Status(http.StatusOK)
}
