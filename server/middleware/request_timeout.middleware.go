package middleware

import (
	"fmt"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/ngdangkietswe/swe-go-common-shared/config"
	"github.com/ngdangkietswe/swe-go-common-shared/logger"
	"time"
)

type RequestTimeoutMiddleware struct {
	logger *logger.Logger
}

func (r RequestTimeoutMiddleware) ShouldSkip(ctx *gin.Context) bool {
	//TODO implement me
	panic("implement me")
}

func (r RequestTimeoutMiddleware) Handle(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (r RequestTimeoutMiddleware) AsMiddleware() gin.HandlerFunc {
	duration := time.Duration(config.GetInt("REQUEST_TIMEOUT_SECONDS", 30))

	return timeout.New(
		timeout.WithTimeout(duration*time.Second),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(func(c *gin.Context) {
			r.logger.Info(fmt.Sprintf("Request to %s timed out after %d seconds", c.Request.URL.Path, duration))
			c.JSON(504, gin.H{"error": "Request timed out"})
		}),
	)
}

func NewRequestTimeoutMiddleware(logger *logger.Logger) Middleware {
	return &RequestTimeoutMiddleware{
		logger: logger,
	}
}
