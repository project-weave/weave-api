package echo

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/project-weave/weave-api/src/internal/types/event"
)

func (s *Server) RegisterEventRoutes() {
	s.server.POST("/v1/event", s.AddEvent)
	s.server.GET("/v1/event/:event_id", s.GetEvent)
	s.server.POST("/v1/event/:event_id/availability", s.UpsertUserEventAvailability)
}

func (s *Server) AddEvent(ctx echo.Context) error {
	var event event.Event
	if err := ctx.Bind(&event); err != nil {
		return s.badRequestResponse(ctx, err)
	}

	if err := s.validate.Struct(event); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return s.validationErrorResponse(ctx, validationErrors)
	}
	id, err := s.EventService.AddEvent(ctx.Request().Context(), &event)
	if err != nil {
		return s.serverErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, envelope{"event_id": id})
}

func (s *Server) GetEvent(ctx echo.Context) error {
	idParam := ctx.Param("event_id")

	eventID, err := event.ToEventUUID(idParam)
	if err != nil {
		return s.notFoundResponse(ctx, err)
	}

	event, responses, err := s.EventService.GetEvent(ctx.Request().Context(), eventID)
	if err != nil {
		return s.notFoundResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, envelope{"event": event, "responses": responses})
}

func (s *Server) UpsertUserEventAvailability(ctx echo.Context) error {
	eventResponse := event.EventResponse{}
	if err := ctx.Bind(&eventResponse); err != nil {
		return s.badRequestResponse(ctx, err)
	}

	if err := s.validate.Struct(eventResponse); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return s.validationErrorResponse(ctx, validationErrors)
	}

	err := s.EventService.UpsertUserEventAvailability(ctx.Request().Context(), eventResponse)
	if err != nil {
		return s.serverErrorResponse(ctx, err)
	}

	return ctx.JSON(http.StatusOK, envelope{"success": true})
}
