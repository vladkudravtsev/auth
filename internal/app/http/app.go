package httpapp

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/vladkudravtsev/auth/internal/config"
	authhttp "github.com/vladkudravtsev/auth/internal/http"
	"github.com/vladkudravtsev/auth/internal/services/auth"

	"github.com/gin-gonic/gin"
)

type App struct {
	httpServer *http.Server
	log        slog.Logger
	address    string
}

func New(authService auth.Service, log *slog.Logger, cfg *config.HTTPServer) *App {
	router := gin.Default()

	httpServer := authhttp.NewServer(router, authService, cfg)

	return &App{
		httpServer: httpServer,
		address:    cfg.Address,
		log:        *log,
	}
}

func (a *App) MustRun() {
	const fn = "httpapp.MustRun"
	if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(fmt.Errorf("%s: %w", fn, err))
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
