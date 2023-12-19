package authhttp

import (
	"log/slog"
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/vladkudravtsev/auth/internal/config"
	"github.com/vladkudravtsev/auth/internal/http/handlers/login"
	"github.com/vladkudravtsev/auth/internal/http/handlers/register"
	"github.com/vladkudravtsev/auth/internal/http/middlewares/logger"
	"github.com/vladkudravtsev/auth/internal/services/auth"

	"github.com/gin-gonic/gin"
)

func NewServer(router *gin.Engine, auth auth.Service, cfg *config.HTTPServer, log *slog.Logger) *http.Server {
	router.Use(requestid.New())
	router.Use(logger.New(log))
	router.POST("/login", login.New(auth, log))
	router.POST("/register", register.New(auth, log))

	return &http.Server{
		Handler:      router,
		Addr:         cfg.Address,
		IdleTimeout:  cfg.IdleTimeout,
		WriteTimeout: cfg.Timeout,
		ReadTimeout:  cfg.Timeout,
	}
}
