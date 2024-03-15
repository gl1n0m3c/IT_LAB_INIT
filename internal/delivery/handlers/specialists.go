package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/services"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/database"
	customErrors "github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/customerr"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/responses"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/validators"
	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
	"net/http"
	"path/filepath"
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
// @Failure 400 {object} responses.MessageResponse "Invalid input data"
// @Failure 403 {object} responses.MessageResponse "JWT is invalid or expired"
// @Failure 500 {object} responses.MessageResponse "Internal server error"
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
		switch {
		case errors.Is(err, customErrors.UserUnverified):
			c.JSON(http.StatusForbidden, responses.NewMessageResponse(err.Error()))
			return
		case errors.Is(err, customErrors.NoRowsCaseErr):
			c.JSON(http.StatusBadRequest, responses.NewMessageResponse(err.Error()))
			return
		case errors.Is(err, customErrors.UniqueRatedErr):
			c.JSON(http.StatusBadRequest, responses.NewMessageResponse(err.Error()))
			return
		case errors.Is(err, customErrors.NoRowsSpecialistIDErr):
			c.JSON(http.StatusBadRequest, responses.NewMessageResponse(err.Error()))
			return
		case errors.Is(err, customErrors.CaseAlreadySolved):
			c.JSON(http.StatusBadRequest, responses.NewMessageResponse(err.Error()))
			return
		case errors.Is(err, customErrors.UserBadLevel):
			c.JSON(http.StatusBadRequest, responses.NewMessageResponse(err.Error()))
			return
		default:
			c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
			return
		}
	}

	c.JSON(http.StatusCreated, responses.CreationIntResponse{ID: createdRatedID})
}

// GetCasesByLevel @Summary Retrieves cases by level
// @Description Retrieves cases based on the provided cursor ID and the user's ID. It returns cases that match the level of difficulty or rating specified for the user.
// @Description Returned cursor can be only int or null. It depends on existence of cases.
// @Tags specialists
// @Accept  json
// @Produce  json
// @Param cursor_id query int true "Cursor ID for pagination"
// @Param authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} models.CaseCursor "Successfully retrieved the cases by level"
// @Failure 400 {object} responses.MessageResponse "Invalid input data or bad query parameter"
// @Failure 401 {object} responses.MessageResponse "JWT is invalid or expired"
// @Failure 403 {object} responses.MessageResponse "User is unverified"
// @Failure 500 {object} responses.MessageResponse "Internal server error"
// @Router /specialist/get_cases_by_level [get]
func (s specialistsHandler) GetCasesByLevel(c *gin.Context) {
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

	userID := c.GetInt("userID")

	cases, err := s.service.GetCasesByLevel(ctx, userID, cursor)
	if err != nil {
		switch {
		case errors.Is(err, customErrors.UserUnverified):
			c.JSON(http.StatusForbidden, responses.NewMessageResponse(err.Error()))
			return
		case errors.Is(err, customErrors.NoRowsSpecialistIDErr):
			c.JSON(http.StatusBadRequest, responses.NewMessageResponse(err.Error()))
			return
		default:
			c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
			return
		}
	}

	c.JSON(http.StatusOK, cases)
}

// GetRatedSolved @Summary Get rated solved
// @Description Retrieves a rated solved entry based on the provided cursor ID.
// @Tags specialists
// @Accept  json
// @Produce  json
// @Param cursor_id query int true "Cursor ID for pagination"
// @Param authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} models.RatedCursor "Successfully retrieved the rated solved"
// @Failure 400 {object} responses.MessageResponse "Invalid input data"
// @Failure 401 {object} responses.MessageResponse "JWT is invalid or expired"
// @Failure 500 {object} responses.MessageResponse "Internal server error"
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

	userID := c.GetInt("userID")

	ratedCursor, err := s.service.GetRatedSolved(ctx, userID, cursor)
	if err != nil {
		switch {
		case errors.Is(err, customErrors.UserUnverified):
			c.JSON(http.StatusForbidden, responses.NewMessageResponse(err.Error()))
			return
		case errors.Is(err, customErrors.NoRowsSpecialistIDErr):
			c.JSON(http.StatusBadRequest, responses.NewMessageResponse(err.Error()))
			return
		default:
			c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
			return
		}
	}

	c.JSON(http.StatusOK, ratedCursor)
}

// GetMe @Summary Get specialist info
// @Description Retrieves information about the current specialist based on their user ID.
// @Tags specialists
// @Accept  json
// @Produce  json
// @Param authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Success 200 {object} models.Specialist "Successfully retrieved the specialist info"
// @Failure 400 {object} responses.MessageResponse "Invalid input data"
// @Failure 401 {object} responses.MessageResponse "JWT is invalid or expired"
// @Failure 500 {object} responses.MessageResponse "Internal server error"
// @Router /specialist/me [get]
func (s specialistsHandler) GetMe(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.GetInt("userID")

	specialist, err := s.service.GetMe(ctx, userID)
	if err != nil {
		if errors.Is(err, customErrors.NoRowsSpecialistIDErr) {
			c.JSON(http.StatusBadRequest, responses.NewMessageResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
		return
	}

	c.JSON(http.StatusOK, specialist)
}

// UpdateMe updates specialist information if it provides.
// @Summary UpdateMain Specialist Information with Photo Upload
// @Description Updates an existing specialist's information including their password, full name, and photo.
// @Description The password must be more than 8 symbols and contain at least one number, one uppercase, and one lowercase letter.
// @Description The photo upload is optional but must be a valid image file if provided.
// @Tags specialists
// @Accept multipart/form-data
// @Produce json
// @Param authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param password formData string false "Password"
// @Param fullname formData string false "Full Name"
// @Param photo formData file false "Photo Upload"
// @Success 204 "Successful update, no content returned"
// @Failure 400 {object} responses.MessageResponse "Invalid input data"
// @Failure 401 {object} responses.MessageResponse "JWT is invalid or expired"
// @Failure 500 {object} responses.MessageResponse "Internal server error, could not process the request"
// @Router /specialist/update [put]
func (s specialistsHandler) UpdateMe(c *gin.Context) {
	if c.Request.ContentLength == 0 {
		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(fmt.Sprintf(responses.Response400, "no data provided")))
		return
	}

	var updateSpecialistData models.SpecialistUpdate

	ctx := c.Request.Context()

	updateSpecialistData.ID = c.GetInt("userID")
	updateSpecialistData.Password = c.PostForm("password")
	updateSpecialistData.FullName = c.PostForm("fullname")

	if updateSpecialistData.Password != "" {
		validate := validator.New()
		_ = validate.RegisterValidation("password", validators.ValidatePassword)
		if err := validate.Struct(updateSpecialistData); err != nil {
			customErrMsg := validators.CustomErrorMessage(err)
			c.JSON(http.StatusBadRequest, responses.NewMessageResponse(customErrMsg))
			return
		}
	}

	if file, err := c.FormFile("photo"); err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			updateSpecialistData.PhotoUrl = ""
		} else {
			c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
			return
		}
	} else {
		// Проверка на ограничение по размеру файла на 2МБ
		err := c.Request.ParseMultipartForm(2 << 20)
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.NewMessageResponse(responses.ResponseBadFileSize))
			return
		}

		// Проверка на допустимый тип `Content-Type` и расширение
		if !validators.ValidateFileTypeExtension(file) {
			c.JSON(http.StatusBadRequest, responses.NewMessageResponse(responses.ResponseBadPhotoFile))
			return
		}

		uuidBytes, err := uuid.NewV4()
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
			return
		}
		uniqueFileName := uuidBytes.String() + filepath.Ext(file.Filename)
		filePath := fmt.Sprintf("/static/img/specialists/%s", uniqueFileName)
		if err := c.SaveUploadedFile(file, ".."+filePath); err != nil {
			c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
			return
		}
		updateSpecialistData.PhotoUrl = filePath
	}

	err := s.service.UpdateMe(ctx, updateSpecialistData)
	if err != nil {
		if errors.Is(err, customErrors.NoRowsSpecialistIDErr) {
			c.JSON(http.StatusBadRequest, responses.NewMessageResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
		return
	}

	c.Status(http.StatusNoContent)
}
