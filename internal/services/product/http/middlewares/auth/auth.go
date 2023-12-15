package auth

import (
	"local/gorm-example/internal/lib/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func New(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := jwt.ValidateToken(c.Request.Header.Get("Authorization"), secret); err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
