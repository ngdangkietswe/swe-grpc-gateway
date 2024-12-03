package utils

import (
	"github.com/gin-gonic/gin"
	"strings"
)

// GetTokenFromReq get token from request. It will get token from Authorization header
func GetTokenFromReq(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}
