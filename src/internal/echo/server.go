package echo

import (
	"context"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/project-weave/weave-api/src/internal/types/event"
)

type Server struct {
	logger       *log.Logger
	server       *echo.Echo
	validate     *validator.Validate
	EventService event.EventService
}

func NewServer(logger *log.Logger, es event.EventService) *Server {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.BodyLimit("2M"))
	e.Use(middleware.CORS())

	validator := validator.New()

	validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// register event validators
	event.RegisterEventServiceValidators(validator)

	s := &Server{
		server:       e,
		logger:       logger,
		validate:     validator,
		EventService: es,
	}

	// register event routes
	s.RegisterEventRoutes()

	return s
}

func (s *Server) Start(addr string) error {
	return s.server.Start(addr)
}

func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		s.server.Logger.Fatal(err)
	}
}
