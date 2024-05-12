package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/project-weave/weave-api/src/internal/echo"
	"github.com/project-weave/weave-api/src/internal/postgres"
)

type config struct {
	port int
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 8080, "server port")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// viper.SetConfigFile(".env")
	// err := viper.ReadInConfig()
	// if err != nil {
	// 	logger.Fatal(err)
	// }

	// dsn := viper.GetString("POSTGRES_DSN")

	dsn := os.Getenv("DATABASE_URL")
	fmt.Println(dsn)

	m, err := migrate.New(
		"file:///usr/src/app/migrations",
		dsn)
	if err != nil {
		fmt.Println("hello")
		logger.Fatal(err)
	}

	logger.Println("Applying migration")

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			logger.Println("No changes were applied to database")
			logger.Fatal(err)
		}
	}

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
