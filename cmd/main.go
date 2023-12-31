package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/vladkudravtsev/auth/internal/app"
	"github.com/vladkudravtsev/auth/internal/config"
	"github.com/vladkudravtsev/auth/internal/database"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	db, err := database.New(cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName)
	if err != nil {
		panic("failed to connect database")
	}

	authApp := app.New(db, log, cfg)

	go authApp.HTTPServer.MustRun()
	go authApp.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("stopping service", slog.String("signal", sign.String()))

	// TODO close db connection

	authApp.GRPCServer.Stop()
	authApp.HTTPServer.Stop()
	log.Info("gracefully stopped")
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case "local":
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log.Fatalf("cannot setup logger with env: %s", env)
	}

	return logger
}
