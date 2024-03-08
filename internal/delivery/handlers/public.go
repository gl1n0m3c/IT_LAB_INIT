package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/decoder"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/models"
	_ "github.com/gl1n0m3c/IT_LAB_INIT/internal/models/swagger"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/services"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/database"
	customErrors "github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/custom_errors"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/jwt"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/responses"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/validators"
	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
	"net/http"
	"path/filepath"
	"strconv"
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
) Public {
	return publicHandler{
		service: service,
		session: session,
		JWTUtil: JWTUtil,
	}
}

// ManagerLogin logs in a manager and returns a JWT and refresh token upon successful login.
// @Summary Manager Login
// @Description Logs in a specialist and returns a JWT and refresh token upon successful login.
// @Tags public
// @Accept json
// @Produce json
// @Param specialist body models.ManagerBase true "Manager Login"
// @Success 201 {object} responses.JWTRefresh "Successful login, returning JWT and refresh token"
// @Failure 400 {object} responses.MessageResponse "Invalid input or incorrect password / login"
// @Failure 500 {object} responses.MessageResponse "Internal server error"
// @Router /public/manager_login [post]
func (p publicHandler) ManagerLogin(c *gin.Context) {
	var manager models.ManagerBase

	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&manager); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(fmt.Sprintf(responses.Response400, err)))
		return
	}

	validate := validator.New()
	if err := validate.Struct(manager); err != nil {
		customErrMsg := validators.CustomErrorMessage(err)
		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(customErrMsg))
		return
	}

	success, specialistData, err := p.service.ManagerLogin(ctx, manager)
	if err != nil {
		if err == customErrors.NoRowsLoginErr {
			c.JSON(http.StatusBadRequest, responses.NewMessageResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
		return
	}

	if !success {
		c.JSON(http.StatusBadRequest, responses.NewMessageResponse("Неверный пароль"))
		return
	}

	accessToken := p.JWTUtil.CreateToken(specialistData.ID, jwt.Manager)

	refreshToken, err := p.session.Set(ctx, database.SessionData{
		UserID:   specialistData.ID,
		UserType: jwt.Manager,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
		return
	}

	c.JSON(http.StatusCreated, responses.NewJWTRefreshResponse(accessToken, refreshToken))
}

// SpecialistRegister registers a new specialist, uploads their photo, and returns a jwt and refresh token upon successful registration.
// @Summary Specialist Registration with Photo Upload
// @Description Registers a new specialist, uploads their photo, and returns a jwt and refresh token upon successful registration.
// @Description Automatically level=1, is_verified=false.
// @Description Login and password are required, along with a photo upload.
// @Description There are some validation on password: More than 8 symbols, contain at least one number, one uppercase and one lowercase letter.
// @Tags public
// @Accept multipart/form-data
// @Produce json
// @Param login formData string true "Login"
// @Param password formData string true "Password"
// @Param fullname formData string false "Full Name"
// @Param photo formData file false "Photo Upload"
// @Success 201 {object} responses.JWTRefresh "Successful registration, returning jwt and refresh token"
// @Failure 400 {object} responses.MessageResponse "Invalid input"
// @Failure 500 {object} responses.MessageResponse "Internal server error"
// @Router /public/specialist_register [post]
func (p publicHandler) SpecialistRegister(c *gin.Context) {
	ctx := c.Request.Context()

	login := c.PostForm("login")
	password := c.PostForm("password")
	fullname := c.PostForm("fullname")

	specialist := models.SpecialistCreate{
		SpecialistBase: models.SpecialistBase{
			Login:    login,
			Password: password,
			Fullname: null.NewString(fullname, fullname != ""),
		},
	}

	validate := validator.New()
	_ = validate.RegisterValidation("password", validators.ValidatePassword)
	if err := validate.Struct(specialist); err != nil {
		customErrMsg := validators.CustomErrorMessage(err)
		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(customErrMsg))
		return
	}

	if file, err := c.FormFile("photo"); err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			specialist.PhotoUrl = null.NewString("", false)
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
		specialist.PhotoUrl = null.NewString(filePath, true)
	}

	ID, err := p.service.SpecialistRegister(ctx, specialist)
	if err != nil {
		if errors.Is(err, customErrors.UniqueSpecialistErr) {
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

// SpecialistLogin logs in a specialist and returns a jwt and refresh token upon successful login.
// @Summary Specialist Login
// @Description Logs in a specialist and returns a jwt and refresh token upon successful login.
// @Tags public
// @Accept json
// @Produce json
// @Param specialist body models.SpecialistLogin true "Specialist Login"
// @Success 201 {object} responses.JWTRefresh "Successful login, returning jwt and refresh token"
// @Failure 400 {object} responses.MessageResponse "Invalid input or incorrect password / login"
// @Failure 500 {object} responses.MessageResponse "Internal server error"
// @Router /public/specialist_login [post]
func (p publicHandler) SpecialistLogin(c *gin.Context) {
	var specialistLogin models.SpecialistLogin

	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&specialistLogin); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(fmt.Sprintf(responses.Response400, err)))
		return
	}

	validate := validator.New()
	if err := validate.Struct(specialistLogin); err != nil {
		customErrMsg := validators.CustomErrorMessage(err)
		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(customErrMsg))
		return
	}

	success, specialistData, err := p.service.SpecialistLogin(ctx, specialistLogin)
	if err != nil {
		if err == customErrors.NoRowsLoginErr {
			c.JSON(http.StatusBadRequest, responses.NewMessageResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
		return
	}

	if !success {
		c.JSON(http.StatusBadRequest, responses.NewMessageResponse("Неверный пароль"))
		return
	}

	accessToken := p.JWTUtil.CreateToken(specialistData.ID, jwt.Specialist)

	refreshToken, err := p.session.Set(ctx, database.SessionData{
		UserID:   specialistData.ID,
		UserType: jwt.Specialist,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
		return
	}

	c.JSON(http.StatusCreated, responses.NewJWTRefreshResponse(accessToken, refreshToken))
}

// CameraCreate creates a new camera and returns its ID upon successful creation.
// @Summary Camera Creation
// @Description Creates a new camera and returns its ID upon successful creation.
// @Tags public
// @Accept json
// @Produce json
// @Param camera body models.CameraBase true "Camera Creation"
// @Success 201 {object} responses.CreationStringResponse "Successful creation, returning camera ID"
// @Failure 400 {object} responses.MessageResponse "Invalid input"
// @Failure 500 {object} responses.MessageResponse "Internal server error"
// @Router /public/camera_create [post]
func (p publicHandler) CameraCreate(c *gin.Context) {
	var camera models.CameraBase

	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&camera); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(fmt.Sprintf(responses.Response400, err)))
		return
	}

	validate := validator.New()
	if err := validate.Struct(camera); err != nil {
		customErrMsg := validators.CustomErrorMessage(err)
		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(customErrMsg))
		return
	}

	createdCameraID, err := p.service.CameraCreate(ctx, camera)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
		return
	}

	c.JSON(http.StatusCreated, responses.CreationStringResponse{ID: createdCameraID})
}

// CameraDelete deletes an existing camera by its ID.
// @Summary Camera Deletion
// @Description Deletes an existing camera by its ID.
// @Tags public
// @Accept json
// @Produce json
// @Param id query int true "Camera ID"
// @Success 204 "Successful deletion"
// @Failure 400 {object} responses.MessageResponse "Invalid input or Camera ID not found"
// @Failure 500 {object} responses.MessageResponse "Internal server error"
// @Router /public/camera_delete [delete]
func (p publicHandler) CameraDelete(c *gin.Context) {
	ctx := c.Request.Context()

	cameraID, ok := c.GetQuery("id")
	if !ok {
		c.JSON(http.StatusBadRequest, responses.NewMessageResponse("ID камеры не указан"))
		return
	}

	err := p.service.CameraDelete(ctx, cameraID)
	if err != nil {
		if errors.Is(err, customErrors.NoRowsCameraErr) {
			c.JSON(http.StatusBadRequest, responses.NewMessageResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
		return
	}

	c.Status(http.StatusNoContent)
}

// CaseCreate creates a new case and returns its ID upon successful creation.
// @Summary Case Creation
// @Description Creates a new case with a photo (.jpeg / .jpg / .png / .svg) and case data in byte string.
// @Tags public
// @Accept multipart/form-data
// @Produce json
// @Param photo formData file true "Photo of the case"
// @Param byte_string formData string true "Case data in byte string format"
// @Success 201 {object} responses.CreationIntResponse "Successful creation, returning case ID"
// @Failure 400 {object} responses.MessageResponse "Invalid input"
// @Failure 500 {object} responses.MessageResponse "Internal server error"
// @Router /public/case_create [post]
func (p publicHandler) CaseCreate(c *gin.Context) {
	ctx := c.Request.Context()

	file, err := c.FormFile("photo")
	if err != nil {
		if err == http.ErrMissingFile {
			c.JSON(http.StatusBadRequest, responses.NewMessageResponse(responses.ResponseNoPhotoProvided))
			return
		} else {
			c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
			return
		}
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
	filePath := fmt.Sprintf("/static/img/cases/%s", uniqueFileName)
	if err := c.SaveUploadedFile(file, ".."+filePath); err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
		return
	}

	bitString := c.PostForm("byte_string")
	if bitString == "" {
		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(responses.ResponseNoByteStringProvided))
		return
	}

	var dataBytes []byte
	for i := 0; i+8 <= len(bitString); i += 8 {
		byteString := bitString[i : i+8]
		byteValue, err := strconv.ParseUint(byteString, 2, 8)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.NewMessageResponse("Ошибка при преобразовании битовой строки в байты"))
			return
		}
		dataBytes = append(dataBytes, byte(byteValue))
	}

	if len(dataBytes) <= 7 {
		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(responses.ResponseBadByteString))
		return
	}

	result, err := decoder.Decoder(bytes.NewBuffer(dataBytes[2:]))
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(responses.ResponseBadByteString))
		return
	}

	cameraType, err := decoder.MapToStruct(result)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(err.Error()))
		return
	}

	caseData, err := cameraType.CameraDataToCaseBase()
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(err.Error()))
		return
	}
	caseData.PhotoUrl = filePath

	createdCaseID, err := p.service.CaseCreate(ctx, caseData)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
		return
	}

	c.JSON(http.StatusCreated, responses.CreationIntResponse{ID: createdCaseID})
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
		case customErrors.NeedToAuthorizeErr:
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
