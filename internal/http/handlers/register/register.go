package register

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-contrib/requestid"
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

func New(authService auth.Service, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := log.With(slog.String("req_id", requestid.Get(c)))
		var body registerRequest

		if err := c.ShouldBindJSON(&body); err != nil {
			var ve validator.ValidationErrors
			if errors.As(err, &ve) {
				log.Warn(err.Error())
				c.JSON(http.StatusBadRequest, response.ValidationError(ve))
				return
			}
			log.Error(err.Error())
			c.JSON(http.StatusBadRequest, response.Error(err.Error()))
			return
		}

		if err := authService.RegisterUser(body.Name, body.Email, body.Password); err != nil {
			// if user already exists
			if errors.Is(err, auth.ErrUserAlreadyExists) {
				log.Warn(err.Error())
				c.JSON(http.StatusBadRequest, response.Error(err.Error()))
				return
			}
			log.Error(err.Error())
			c.JSON(http.StatusServiceUnavailable, response.Error(err.Error()))
			return
		}

		c.AbortWithStatus(http.StatusCreated)
	}
}
