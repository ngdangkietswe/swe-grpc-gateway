package middleware

import "github.com/gin-gonic/gin"

type Middleware interface {
	ShouldSkip(ctx *gin.Context) bool
	Handle(ctx *gin.Context)
	AsMiddleware() gin.HandlerFunc
}
