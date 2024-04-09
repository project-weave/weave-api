package echo

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func (s *Server) errorResponse(ctx echo.Context, status int, message any) error {
	env := envelope{"error": message}
	return ctx.JSON(status, env)
}

func (s *Server) serverErrorResponse(ctx echo.Context, err error) error {
	s.server.Logger.Print(err.Error())

	message := "your request cannot be processed due to a server error"
	return s.errorResponse(ctx, http.StatusInternalServerError, message)
}

func (s *Server) notFoundResponse(ctx echo.Context, err error) error {
	s.server.Logger.Print(err.Error())

	message := "the requested resource could not be found"
	return s.errorResponse(ctx, http.StatusNotFound, message)
}

func (s *Server) badRequestResponse(ctx echo.Context, err error) error {
	return s.errorResponse(ctx, http.StatusBadRequest, err.Error())
}

type validationError struct {
	Namespace       string `json:"namespace"` // can differ when a custom TagNameFunc is registered or
	Field           string `json:"field"`     // by passing alt name to ReportError like below
	StructNamespace string `json:"structNamespace"`
	StructField     string `json:"structField"`
	Tag             string `json:"tag"`
	ActualTag       string `json:"actualTag"`
	Kind            string `json:"kind"`
	Type            string `json:"type"`
	Value           string `json:"value"`
	Param           string `json:"param"`
	Message         string `json:"message"`
}

func (s *Server) validationErrorResponse(ctx echo.Context, errors validator.ValidationErrors) error {
	output := make([]validationError, len(errors))

	for i, err := range errors {
		e := validationError{
			Namespace:       err.Namespace(),
			Field:           err.Field(),
			StructNamespace: err.StructNamespace(),
			StructField:     err.StructField(),
			Tag:             err.Tag(),
			ActualTag:       err.ActualTag(),
			Kind:            fmt.Sprintf("%v", err.Kind()),
			Type:            fmt.Sprintf("%v", err.Type()),
			Value:           fmt.Sprintf("%v", err.Value()),
			Param:           err.Param(),
			Message:         err.Error(),
		}
		output[i] = e
	}

	return s.errorResponse(ctx, http.StatusBadRequest, output)
}
