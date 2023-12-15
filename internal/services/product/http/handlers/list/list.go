package list

import (
	"local/gorm-example/internal/services/product"
	"local/gorm-example/internal/services/product/http/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

func New(product product.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := product.List()

		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": err})
			return
		}

		var response []*api.ProductResponse

		for _, product := range products {
			response = append(response, api.ToProductResponse(product))
		}

		c.JSON(http.StatusOK, response)
	}
}
