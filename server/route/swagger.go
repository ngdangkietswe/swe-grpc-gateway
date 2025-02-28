package route

import "github.com/gin-gonic/gin"

// RegisterSwagger is a function that registers the swagger route
func RegisterSwagger(router *gin.Engine) {
	router.Static("/swagger", "./swagger")
}
