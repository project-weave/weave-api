package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/project-weave/weave-api/src/internal/echo"
	"github.com/project-weave/weave-api/src/internal/postgres"
	"github.com/spf13/viper"
)

type config struct {
	port int
}

const version = "1"

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 8080, "server port")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatal(err)
	}

	dsn := viper.GetString("POSTGRES_DSN")

	db := postgres.NewDB(dsn)
	err = db.Open()
	if err != nil {
		logger.Fatal(err)
	}

	eventService := postgres.NewEventService(db)

	server := echo.NewServer(logger, eventService)

	go func() {
		if err := server.Start(fmt.Sprintf(":%d", cfg.port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	server.Shutdown()
}
