package echo

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/project-weave/weave-api/src/internal/domain"
)

func (s *Server) RegisterEventRoutes() {
	s.server.POST("/event", s.AddEvent)
	s.server.GET("/event/:id", s.GetEvent)

	s.server.GET("/event/:id/availability", s.GetEventAvailability)
	s.server.POST("/event/:id/availability", s.UpdateEventAvailability)
}

func (s *Server) AddEvent(ctx echo.Context) error {
	var event *domain.Event
	if err := ctx.Bind(event); err != nil {
		return s.badRequestResponse(ctx, err)
	}

	validationErrors := domain.ValidateEvent(event)
	if len(validationErrors) != 0 {
		return s.validationErrorResponse(ctx, validationErrors)
	}

	if err := s.EventService.AddEvent(event); err != nil {
		return s.serverErrorResponse(ctx, err)
	}
	return ctx.JSON(http.StatusOK, event)
}

func (s *Server) GetEvent(ctx echo.Context) error {
	idString := ctx.Param("id")
	eventUUID, err := uuid.Parse(idString)
	if err != nil {
		return s.notFoundResponse(ctx, err)
	}

	event, err := s.EventService.GetEvent(eventUUID)
	if err != nil {
		return s.notFoundResponse(ctx, err)
	}
	return ctx.JSON(http.StatusOK, event)
}

func (s *Server) GetEventAvailability(ctx echo.Context) error {
	return nil
}

func (s *Server) UpdateEventAvailability(ctx echo.Context) error {
	return nil
}
