package echo

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) errorResponse(ctx echo.Context, status int, message any) error {
	env := envelope{"error": message}
	return ctx.JSON(status, env)
}

func (s *Server) serverErrorResponse(ctx echo.Context, err error) error {
	s.Logger.Println(err.Error())

	message := "your request cannot be processed due to a server error"
	return s.errorResponse(ctx, http.StatusInternalServerError, message)
}

func (s *Server) notFoundResponse(ctx echo.Context, err error) error {
	message := "the requested resource could not be found"
	return s.errorResponse(ctx, http.StatusNotFound, message)
}

func (s *Server) badRequestResponse(ctx echo.Context, err error) error {
	return s.errorResponse(ctx, http.StatusBadRequest, err.Error())
}

func (s *Server) validationErrorResponse(ctx echo.Context, errors map[string]string) error {
	return s.errorResponse(ctx, http.StatusUnprocessableEntity, errors)
}
