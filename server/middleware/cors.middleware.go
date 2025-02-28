package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ngdangkietswe/swe-go-common-shared/logger"
)

type CORSMiddleware struct {
	logger *logger.Logger
}

func (c CORSMiddleware) ShouldSkip(ctx *gin.Context) bool {
	//TODO implement me
	panic("implement me")
}

func (c CORSMiddleware) Handle(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*") // You can specify specific origins instead of "*"
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
	ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	// Handle preflight request (OPTIONS)
	if ctx.Request.Method == "OPTIONS" {
		ctx.AbortWithStatus(204)
		return
	}

	ctx.Next()
	return
}

func (c CORSMiddleware) AsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c.Handle(ctx)
	}
}

func NewCORSMiddleware(logger *logger.Logger) Middleware {
	return &CORSMiddleware{
		logger: logger,
	}
}
