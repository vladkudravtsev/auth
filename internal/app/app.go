package app

import (
	"log/slog"

	grpcapp "github.com/vladkudravtsev/auth/internal/app/grpc"
	httpapp "github.com/vladkudravtsev/auth/internal/app/http"
	"github.com/vladkudravtsev/auth/internal/config"
	"github.com/vladkudravtsev/auth/internal/services/auth"

	"gorm.io/gorm"
)

type App struct {
	HTTPServer *httpapp.App
	GRPCServer *grpcapp.App
}

func New(db *gorm.DB, log *slog.Logger, cfg *config.Config) *App {
	authService := auth.NewService(db, log, cfg.Auth.TokenTTL)

	httpApp := httpapp.New(authService, log, &cfg.HTTPServer)
	gRPCApp := grpcapp.New(authService, log, &cfg.GRPCServer)

	return &App{
		HTTPServer: httpApp,
		GRPCServer: gRPCApp,
	}
}
