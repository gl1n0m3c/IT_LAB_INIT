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
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/tracing"
	customErrors "github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/customerr"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/jwt"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/responses"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/validators"
	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"path/filepath"
	"strconv"
)

type publicHandler struct {
	service services.Public
	session database.Session
	JWTUtil jwt.JWT
	tracer  trace.Tracer
}

func InitPublicHandler(
	service services.Public,
	session database.Session,
	JWTUtil jwt.JWT,
	tracer trace.Tracer,
) Public {
	return publicHandler{
		service: service,
		session: session,
		JWTUtil: JWTUtil,
		tracer:  tracer,
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

	ctx, span := p.tracer.Start(c.Request.Context(), tracing.ManagerLogin)
	defer span.End()

	if err := c.ShouldBindJSON(&manager); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(fmt.Sprintf(responses.Response400, err)))
		return
	}

	validate := validator.New()
	if err := validate.Struct(manager); err != nil {
		customErrMsg := validators.CustomErrorMessage(err)

		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.ValidationType, customErrMsg)),
		)
		span.SetStatus(codes.Error, customErrMsg)

		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(customErrMsg))
		return
	}

	span.AddEvent(tracing.CallToService)
	success, specialistData, err := p.service.ManagerLogin(ctx, manager)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.ManagerLoginType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		if errors.Is(err, customErrors.NoRowsSpecialistLoginErr) {
			c.JSON(http.StatusBadRequest, responses.NewMessageResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
		return
	}

	if !success {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.AccessType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusBadRequest, responses.NewMessageResponse("Неверный пароль"))
		return
	}

	accessToken := p.JWTUtil.CreateToken(specialistData.ID, jwt.Manager)

	refreshToken, err := p.session.Set(ctx, database.SessionData{
		UserID:   specialistData.ID,
		UserType: jwt.Manager,
	})

	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.SessionType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
		return
	}

	span.SetStatus(codes.Ok, tracing.SuccessfulCompleting)

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
	ctx, span := p.tracer.Start(c.Request.Context(), tracing.SpecialistRegister)
	defer span.End()

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

		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.ValidationType, customErrMsg)),
		)
		span.SetStatus(codes.Error, customErrMsg)

		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(customErrMsg))
		return
	}

	if file, err := c.FormFile("photo"); err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			specialist.PhotoUrl = null.NewString("", false)
		} else {
			span.RecordError(err, trace.WithAttributes(
				attribute.String(tracing.FileType, err.Error())),
			)
			span.SetStatus(codes.Error, err.Error())

			c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
			return
		}
	} else {
		// Проверка на ограничение по размеру файла на 2МБ
		err := c.Request.ParseMultipartForm(2 << 20)
		if err != nil {
			span.RecordError(err, trace.WithAttributes(
				attribute.String(tracing.FileType, err.Error())),
			)
			span.SetStatus(codes.Error, err.Error())

			c.JSON(http.StatusBadRequest, responses.NewMessageResponse(responses.ResponseBadFileSize))
			return
		}

		// Проверка на допустимый тип `Content-Type` и расширение
		if !validators.ValidateFileTypeExtension(file) {
			span.RecordError(err, trace.WithAttributes(
				attribute.String(tracing.FileType, err.Error())),
			)
			span.SetStatus(codes.Error, err.Error())

			c.JSON(http.StatusBadRequest, responses.NewMessageResponse(responses.ResponseBadPhotoFile))
			return
		}

		uuidBytes, err := uuid.NewV4()
		if err != nil {
			span.RecordError(err, trace.WithAttributes(
				attribute.String(tracing.InternalErr, err.Error())),
			)
			span.SetStatus(codes.Error, err.Error())

			c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
			return
		}
		uniqueFileName := uuidBytes.String() + filepath.Ext(file.Filename)
		filePath := fmt.Sprintf("/static/img/specialists/%s", uniqueFileName)
		if err := c.SaveUploadedFile(file, ".."+filePath); err != nil {
			span.RecordError(err, trace.WithAttributes(
				attribute.String(tracing.InternalErr, err.Error())),
			)
			span.SetStatus(codes.Error, err.Error())

			c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
			return
		}
		specialist.PhotoUrl = null.NewString(filePath, true)
	}

	span.AddEvent(tracing.CallToService)
	ID, err := p.service.SpecialistRegister(ctx, specialist)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.SpecialistRegisterType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

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
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.SessionType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
		return
	}

	span.SetStatus(codes.Ok, tracing.SuccessfulCompleting)

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

	ctx, span := p.tracer.Start(c.Request.Context(), tracing.SpecialistLogin)
	defer span.End()

	if err := c.ShouldBindJSON(&specialistLogin); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(fmt.Sprintf(responses.Response400, err)))
		return
	}

	validate := validator.New()
	if err := validate.Struct(specialistLogin); err != nil {
		customErrMsg := validators.CustomErrorMessage(err)

		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.ValidationType, customErrMsg)),
		)
		span.SetStatus(codes.Error, customErrMsg)

		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(customErrMsg))
		return
	}

	span.AddEvent(tracing.CallToService)
	success, specialistData, err := p.service.SpecialistLogin(ctx, specialistLogin)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.SpecialistLoginType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		if errors.Is(err, customErrors.NoRowsSpecialistLoginErr) {
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
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.SessionType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
		return
	}

	span.SetStatus(codes.Ok, tracing.SuccessfulCompleting)

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

	ctx, span := p.tracer.Start(c.Request.Context(), tracing.CameraCreate)
	defer span.End()

	if err := c.ShouldBindJSON(&camera); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.BindType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(fmt.Sprintf(responses.Response400, err)))
		return
	}

	validate := validator.New()
	if err := validate.Struct(camera); err != nil {
		customErrMsg := validators.CustomErrorMessage(err)

		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.ValidationType, customErrMsg)),
		)
		span.SetStatus(codes.Error, customErrMsg)

		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(customErrMsg))
		return
	}

	span.AddEvent(tracing.CallToService)
	createdCameraID, err := p.service.CameraCreate(ctx, camera)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.CameraCreateType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
		return
	}

	span.SetStatus(codes.Ok, tracing.SuccessfulCompleting)

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
	ctx, span := p.tracer.Start(c.Request.Context(), tracing.CameraDelete)
	defer span.End()

	cameraID, ok := c.GetQuery("id")
	if !ok {
		er := fmt.Errorf("bad `id` query provided")
		span.RecordError(er, trace.WithAttributes(
			attribute.String(tracing.QueryType, er.Error())),
		)
		span.SetStatus(codes.Error, er.Error())

		c.JSON(http.StatusBadRequest, responses.NewMessageResponse("ID камеры не указан"))
		return
	}

	span.AddEvent(tracing.CallToService)
	err := p.service.CameraDelete(ctx, cameraID)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.CameraDeleteType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		if errors.Is(err, customErrors.NoRowsCameraErr) {
			c.JSON(http.StatusBadRequest, responses.NewMessageResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
		return
	}

	span.SetStatus(codes.Ok, tracing.SuccessfulCompleting)

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
	ctx, span := p.tracer.Start(c.Request.Context(), tracing.CaseCreate)
	defer span.End()

	file, err := c.FormFile("photo")
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.FileType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		if errors.Is(err, http.ErrMissingFile) {
			c.JSON(http.StatusBadRequest, responses.NewMessageResponse(responses.ResponseNoPhotoProvided))
			return
		} else {
			c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
			return
		}
	}

	// Проверка на допустимый тип `Content-Type` и расширение
	if !validators.ValidateFileTypeExtension(file) {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.FileType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(responses.ResponseBadPhotoFile))
		return
	}

	uuidBytes, err := uuid.NewV4()
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.InternalErr, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
		return
	}
	uniqueFileName := uuidBytes.String() + filepath.Ext(file.Filename)
	filePath := fmt.Sprintf("/static/img/cases/%s", uniqueFileName)
	if err := c.SaveUploadedFile(file, ".."+filePath); err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.InternalErr, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

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
			span.RecordError(err, trace.WithAttributes(
				attribute.String(tracing.InternalErr, err.Error())),
			)
			span.SetStatus(codes.Error, err.Error())

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
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.DecoderType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(responses.ResponseBadByteString))
		return
	}

	cameraType, err := decoder.MapToStruct(result)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.DecoderType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(err.Error()))
		return
	}

	caseData, err := cameraType.CameraDataToCaseBase()
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.DecoderType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(err.Error()))
		return
	}
	caseData.PhotoUrl = filePath

	span.AddEvent(tracing.CallToService)
	createdCaseID, err := p.service.CaseCreate(ctx, caseData)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.CaseCreateType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
		return
	}

	span.SetStatus(codes.Ok, tracing.SuccessfulCompleting)

	c.JSON(http.StatusCreated, responses.CreationIntResponse{ID: createdCaseID})
}

// Refresh updates access and refresh tokens
// @Summary Refresh Tokens
// @Description Refreshes access and refresh tokens using a refresh token provided in the Authorization header.
// @Tags public
// @Accept json
// @Produce json
// @Param refresh header string true "Refresh Token"
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

	ctx, span := p.tracer.Start(c.Request.Context(), tracing.Refresh)
	defer span.End()

	span.AddEvent(tracing.CallToService)
	newRefreshToken, userData, err := p.session.GetAndUpdate(ctx, oldRefreshToken)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.RefreshType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		switch {
		case errors.Is(err, customErrors.NeedToAuthorizeErr):
			c.JSON(http.StatusUnauthorized, responses.NewMessageResponse(err.Error()))
			return
		default:
			c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(err.Error()))
			return
		}
	}

	newAccessToken := p.JWTUtil.CreateToken(userData.UserID, userData.UserType)

	span.SetStatus(codes.Ok, tracing.SuccessfulCompleting)

	c.JSON(http.StatusOK, responses.NewJWTRefreshResponse(newAccessToken, newRefreshToken))
}
