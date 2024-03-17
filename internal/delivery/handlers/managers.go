package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gl1n0m3c/IT_LAB_INIT/internal/services"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/tracing"
	customErrors "github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/customerr"
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/utils/responses"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strconv"
	"time"
)

type managerHandler struct {
	service services.Managers
	tracer  trace.Tracer
}

func InitManagerHandler(
	service services.Managers,
	tracer trace.Tracer,
) Managers {
	return managerHandler{
		service: service,
		tracer:  tracer,
	}
}

// GetFulCaseByID @Summary Retrieve a case by ID
// @Description Retrieves a case by its ID and returns detailed information about the case.
// @Description Field `rated_covers` could be null if there are no ratings
// @Tags managers
// @Accept  json
// @Produce  json
// @Param authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param case_id query int true "ID of the case to retrieve"
// @Success 200 {object} models.CaseFul "Successfully retrieved the case details"
// @Failure 400 {object} responses.MessageResponse "Invalid query parameter or missing case_id"
// @Failure 403 {object} responses.MessageResponse "JWT is invalid or expired"
// @Failure 500 {object} responses.MessageResponse "Internal server error"
// @Router /manager/get_case [get]
func (m managerHandler) GetFulCaseByID(c *gin.Context) {
	ctx, span := m.tracer.Start(c.Request.Context(), tracing.GetFulCaseByID)
	defer span.End()

	caseIDStr, ok := c.GetQuery("case_id")
	if !ok {
		er := fmt.Errorf("bad `case_id` query provided")
		span.RecordError(er, trace.WithAttributes(
			attribute.String(tracing.QueryType, er.Error())),
		)
		span.SetStatus(codes.Error, er.Error())

		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(responses.ResponseBadQuery))
		return
	}

	caseID, err := strconv.Atoi(caseIDStr)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.QueryType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(responses.ResponseBadQuery))
		return
	}

	caseData, err := m.service.GetFulCaseByID(ctx, caseID)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.GetFulCaseByIDType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		if errors.Is(err, customErrors.NoRowsCaseErr) {
			c.JSON(http.StatusNotFound, responses.NewMessageResponse(err.Error()))
			return
		}
		c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
		return
	}

	span.SetStatus(codes.Ok, tracing.SuccessfulCompleting)

	c.JSON(http.StatusOK, caseData)
}

// GetSpecialistRating @Summary Retrieve specialists' ratings
// @Description Retrieves a list of specialists ratings within a specified time range, paginated by a cursor.
// @Description Time example (2023-04-12T15:04:05Z - without time zone / 2023-04-12T15:04:05+07:00 - with time zone)
// @Tags managers
// @Accept  json
// @Produce  json
// @Param authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param cursor query int false "Cursor for pagination"
// @Param time_from query string true "Start time for filtering ratings (inclusive), in RFC3339 format"
// @Param time_to query string true "End time for filtering ratings (inclusive), in RFC3339 format"
// @Success 200 {object} models.RatingSpecialistCountCursor "Successfully retrieved the specialists' ratings"
// @Failure 400 {object} responses.MessageResponse "Invalid query parameter or missing required fields"
// @Failure 403 {object} responses.MessageResponse "JWT is invalid or expired"
// @Failure 500 {object} responses.MessageResponse "Internal server error"
// @Router /manager/get_specialists_rating [get]
func (m managerHandler) GetSpecialistRating(c *gin.Context) {
	ctx, span := m.tracer.Start(c.Request.Context(), tracing.GetSpecialistRating)
	defer span.End()

	cursorStr, ok := c.GetQuery("cursor")
	if !ok {
		er := fmt.Errorf("bad `cursor` query provided")
		span.RecordError(er, trace.WithAttributes(
			attribute.String(tracing.QueryType, er.Error())),
		)
		span.SetStatus(codes.Error, er.Error())

		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(responses.ResponseBadQuery))
		return
	}

	cursor, err := strconv.Atoi(cursorStr)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.QueryType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(responses.ResponseBadQuery))
		return
	}

	timeFromStr, ok := c.GetQuery("time_from")
	if !ok {
		er := fmt.Errorf("bad `time_from` query provided")
		span.RecordError(er, trace.WithAttributes(
			attribute.String(tracing.QueryType, er.Error())),
		)
		span.SetStatus(codes.Error, er.Error())

		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(responses.ResponseBadQuery))
		return
	}

	timeToStr, ok := c.GetQuery("time_to")
	if !ok {
		er := fmt.Errorf("bad `time_to` query provided")
		span.RecordError(er, trace.WithAttributes(
			attribute.String(tracing.QueryType, er.Error())),
		)
		span.SetStatus(codes.Error, er.Error())

		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(responses.ResponseBadQuery))
		return
	}

	timeFrom, err := time.Parse(time.RFC3339, timeFromStr)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.TimeFormatType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(responses.ResponseBadTime))
		return
	}

	timeTo, err := time.Parse(time.RFC3339, timeToStr)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.TimeFormatType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(responses.ResponseBadTime))
		return
	}

	if timeFrom.After(timeTo) || timeFrom.After(time.Now()) {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.TimeFormatType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusBadRequest, responses.NewMessageResponse(responses.ResponseBadTime))
		return
	}

	specialists, err := m.service.GetSpecialistRating(ctx, timeFrom, timeTo, cursor)
	if err != nil {
		span.RecordError(err, trace.WithAttributes(
			attribute.String(tracing.GetSpecialistRatingType, err.Error())),
		)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusInternalServerError, responses.NewMessageResponse(responses.Response500))
		return
	}

	span.SetStatus(codes.Ok, tracing.SuccessfulCompleting)

	c.JSON(http.StatusOK, specialists)
}
