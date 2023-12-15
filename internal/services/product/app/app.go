package app

import (
	"local/gorm-example/internal/config"
	"local/gorm-example/internal/services/product"
	grpcapp "local/gorm-example/internal/services/product/app/grpc"
	httpapp "local/gorm-example/internal/services/product/app/http"
	"local/gorm-example/internal/services/product/models"
	"log/slog"

	"gorm.io/gorm"
)

type App struct {
	HTTPServer *httpapp.App
	GRPCServer *grpcapp.App
}

func New(db *gorm.DB, log *slog.Logger, cfg *config.Config) *App {
	db.AutoMigrate(&models.Product{})

	productService := product.NewService(db, log)

	httpServer := httpapp.New(productService, cfg, log)
	grpcServer := grpcapp.New(productService, log, &cfg.GRPCServer)

	return &App{
		HTTPServer: httpServer,
		GRPCServer: grpcServer,
	}
}
