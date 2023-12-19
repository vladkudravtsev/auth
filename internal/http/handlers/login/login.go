package login

import (
	"errors"
	"net/http"

	"github.com/vladkudravtsev/auth/internal/lib/api/response"
	"github.com/vladkudravtsev/auth/internal/services/auth"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Email    string `json:"email" binding:"email,required" name:"email"`
	Password string `json:"password" binding:"required" name:"password"`
	AppID    uint   `json:"app_id" binding:"required" name:"app_id"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func New(authService auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body loginRequest

		if err := c.ShouldBindJSON(&body); err != nil {
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				c.JSON(http.StatusBadRequest, response.ValidationError(ve))
				return
			}

			c.JSON(http.StatusServiceUnavailable, response.Error(err.Error()))
			return
		}

		token, err := authService.Login(body.Email, body.Password, body.AppID)

		if err != nil {
			// if invalid creds
			if errors.Is(err, auth.ErrInvalidCredentials) {
				c.JSON(http.StatusUnauthorized, response.Error(err.Error()))
				return
			}

			c.JSON(http.StatusServiceUnavailable, response.Error(err.Error()))
			return
		}

		c.JSON(http.StatusOK, &loginResponse{Token: token})
	}
}
