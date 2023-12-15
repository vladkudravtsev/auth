package find

import (
	"local/gorm-example/internal/services/product"
	"local/gorm-example/internal/services/product/http/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

func New(product product.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params api.ProductParams

		if err := c.ShouldBindUri(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		}

		product, err := product.GetOne(params.ID)

		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, api.ToProductResponse(*product))
	}
}
