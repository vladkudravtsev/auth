package register

import (
	"errors"
	"net/http"

	"github.com/vladkudravtsev/auth/internal/lib/api/response"
	"github.com/vladkudravtsev/auth/internal/services/auth"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type registerRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func New(authService auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body registerRequest

		if err := c.ShouldBindJSON(&body); err != nil {
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				c.JSON(http.StatusBadRequest, response.ValidationError(ve))
				return
			}

			c.JSON(http.StatusBadRequest, response.Error(err.Error()))
			return
		}

		if err := authService.RegisterUser(body.Name, body.Email, body.Password); err != nil {
			// if user already exists
			if errors.Is(err, auth.ErrUserAlreadyExists) {
				c.JSON(http.StatusBadRequest, response.Error(err.Error()))
				return
			}

			c.JSON(http.StatusServiceUnavailable, response.Error(err.Error()))
			return
		}

		c.AbortWithStatus(http.StatusCreated)
	}
}
