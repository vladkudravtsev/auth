package update

import (
	"local/gorm-example/internal/services/product"
	"local/gorm-example/internal/services/product/http/api"
	"local/gorm-example/internal/services/product/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type productUpdateRequest struct {
	Code  string `json:"code"`
	Price uint   `json:"price"`
}

func New(product product.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params api.ProductParams
		var body productUpdateRequest

		if err := c.ShouldBindUri(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if err := product.Update(params.ID, &models.Product{
			Price: body.Price,
			Code:  body.Code,
		}); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		c.AbortWithStatus(http.StatusOK)
	}
}
