package create

import (
	"local/gorm-example/internal/services/product"
	"local/gorm-example/internal/services/product/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type productCreateRequest struct {
	Code  string `json:"code" binding:"required"`
	Price uint   `json:"price" binding:"required"`
}

func New(product product.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body productCreateRequest

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if err := product.Create(&models.Product{
			Price: body.Price,
			Code:  body.Code,
		}); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": err.Error()})
			return
		}

		c.AbortWithStatus(http.StatusCreated)
	}
}
