package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ngdangkietswe/swe-gateway-service/utils"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
}

// ShouldSkip is a middleware function that checks if the request should skip the auth middleware
func (a AuthMiddleware) ShouldSkip(ctx *gin.Context) bool {
	publicRoutes := []string{
		"/api/v1/auth/login",
		"/api/v1/auth/register",
	}

	for _, route := range publicRoutes {
		if strings.HasPrefix(ctx.Request.URL.Path, route) {
			return true
		}
	}

	return false
}

// Handle is a middleware function that checks if the request has a valid token
func (a AuthMiddleware) Handle(ctx *gin.Context) {
	token := utils.GetTokenFromReq(ctx)
	if token == "" {
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Token is required"})
		ctx.Abort()
		return
	}

	_, err := utils.ValidateToken(token)
	if err != nil {
		ctx.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		ctx.Abort()
		return
	}

	ctx.Next()
	return
}

func (a AuthMiddleware) AsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if a.ShouldSkip(ctx) {
			ctx.Next()
			return
		}

		a.Handle(ctx)
	}
}

func NewAuthMiddleware() Middleware {
	return &AuthMiddleware{}
}
