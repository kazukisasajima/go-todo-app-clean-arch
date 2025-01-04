package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"

	"go-todo-app-clean-arch/adapter/controller/gin/presenter"
)

func TimeoutMiddleware(duration time.Duration) gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(duration),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(func(c *gin.Context) {
			c.JSON(http.StatusRequestTimeout, &presenter.ErrorResponse{Message: "timeout"})
			c.Abort()
		}),
	)
}
