package authhttp

import (
	"net/http"

	"github.com/vladkudravtsev/auth/internal/config"
	"github.com/vladkudravtsev/auth/internal/http/handlers/login"
	"github.com/vladkudravtsev/auth/internal/http/handlers/register"
	"github.com/vladkudravtsev/auth/internal/services/auth"

	"github.com/gin-gonic/gin"
)

func NewServer(router *gin.Engine, auth auth.Service, cfg *config.HTTPServer) *http.Server {
	router.POST("/login", login.New(auth))
	router.POST("/register", register.New(auth))

	return &http.Server{
		Handler:      router,
		Addr:         cfg.Address,
		IdleTimeout:  cfg.IdleTimeout,
		WriteTimeout: cfg.Timeout,
		ReadTimeout:  cfg.Timeout,
	}
}
