package main

import (
	"local/gorm-example/internal/config"
	"local/gorm-example/internal/database"
	"local/gorm-example/internal/services/product/app"
	"log"
	"log/slog"
	"os"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	db, err := database.New(cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName)
	if err != nil {
		panic("failed to connect database")
	}

	productApp := app.New(db, log, cfg)

	// productApp.HTTPServer.MustRun()
	productApp.GRPCServer.MustRun()
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
