package authhttp

import (
	"local/gorm-example/internal/config"
	"local/gorm-example/internal/http/handlers/login"
	"local/gorm-example/internal/http/handlers/register"
	"local/gorm-example/internal/services/auth"
	"net/http"

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
