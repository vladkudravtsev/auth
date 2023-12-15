package app

import (
	"local/gorm-example/internal/config"
	"local/gorm-example/internal/services/auth"
	grpcapp "local/gorm-example/internal/services/auth/app/grpc"
	httpapp "local/gorm-example/internal/services/auth/app/http"
	"local/gorm-example/internal/services/auth/models"
	"log/slog"

	"gorm.io/gorm"
)

type App struct {
	HTTPServer *httpapp.App
	GRPCServer *grpcapp.App
}

func New(db *gorm.DB, log *slog.Logger, cfg *config.Config) *App {
	db.AutoMigrate(&models.User{}, &models.App{})

	authService := auth.NewService(db, log, cfg.Auth.TokenTTL)

	httpApp := httpapp.New(authService, log, &cfg.HTTPServer)
	gRPCApp := grpcapp.New(authService, log, &cfg.GRPCServer)

	return &App{
		HTTPServer: httpApp,
		GRPCServer: gRPCApp,
	}
}
