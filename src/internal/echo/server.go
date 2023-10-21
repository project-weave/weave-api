package echo

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/project-weave/weave-api/src/internal/domain"
	"github.com/project-weave/weave-api/src/internal/validator"
)

type Server struct {
	server    *echo.Echo
	validator validator.Validator

	Logger       *log.Logger
	EventService domain.EventService
}

func NewServer() *Server {
	e := echo.New()
	e.Use(middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{}))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	s := &Server{
		server: e,
	}

	s.Logger = log.New(os.Stdout, "Server: ", log.Ldate|log.Lshortfile)
	//s.EventService = postgres.EventService

	return s
}

func (s *Server) Start(addr string) error {
	return s.server.Start(addr)
}

func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		s.Logger.Fatal(err)
	}
}
