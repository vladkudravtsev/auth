package httpapp

import (
	"context"
	"fmt"
	"local/gorm-example/internal/config"
	"local/gorm-example/internal/services/product"
	producthttp "local/gorm-example/internal/services/product/http"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type App struct {
	httpServer *http.Server
	address    string
	secret     string
	log        *slog.Logger
}

func New(productService product.Service, cfg *config.Config, log *slog.Logger) *App {
	router := gin.Default()

	httpServer := producthttp.New(router, productService, cfg)

	return &App{
		httpServer: httpServer,
		address:    cfg.HTTPServer.Address,
		secret:     cfg.AppSecret,
		log:        log,
	}
}

func (a *App) MustRun() {
	const fn = "httpapp.MustRun"

	log := a.log.With("fn", fn)
	log.Info("http server is running", slog.String("address", a.address))

	if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(fmt.Errorf("%s: %w", fn, err))
	}

}

func (a *App) Stop() {
	const fn = "httpapp.Stop"

	log := a.log.With(slog.String("fn", fn))
	log.Info("stopping http server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		log.Error("failed to stop http server", slog.Any("err", err))
		return
	}

	log.Info("http server stopped")
}
