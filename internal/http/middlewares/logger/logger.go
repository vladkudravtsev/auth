package logger

import (
	"log/slog"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

func New(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		entry := log.With(
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("req_id", requestid.Get(c)),
			slog.String("user_agent", c.Request.UserAgent()),
		)

		defer entry.Info("request completed")
	}
}
