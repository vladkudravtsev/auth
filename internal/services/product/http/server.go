package producthttp

import (
	"local/gorm-example/internal/config"
	"local/gorm-example/internal/services/product"
	"local/gorm-example/internal/services/product/http/handlers/create"
	"local/gorm-example/internal/services/product/http/handlers/delete"
	"local/gorm-example/internal/services/product/http/handlers/find"
	"local/gorm-example/internal/services/product/http/handlers/list"
	"local/gorm-example/internal/services/product/http/handlers/update"
	"local/gorm-example/internal/services/product/http/middlewares/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func New(router *gin.Engine, product product.Service, cfg *config.Config) *http.Server {
	router.Use(auth.New(cfg.AppSecret))
	router.GET("/products", list.New(product))
	router.GET("/products/:id", find.New(product))
	router.POST("/products", create.New(product))
	router.DELETE("/products/:id", delete.New(product))
	router.PATCH("/products/:id", update.New(product))

	return &http.Server{
		Handler:      router,
		Addr:         cfg.HTTPServer.Address,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
}
