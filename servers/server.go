package servers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ngdangkietswe/swe-gateway-service/configs"
	"github.com/ngdangkietswe/swe-gateway-service/servers/middleware"
	"github.com/ngdangkietswe/swe-gateway-service/servers/route"
)

type Server struct {
	router *gin.Engine
}

func NewServer() Server {
	return Server{
		router: gin.Default(),
	}
}

// Init is a function that initializes the server
func (server *Server) Init() {
	server.router.Use(middleware.NewAuthMiddleware().AsMiddleware())
	route.RegisterGrpcGateway(server.router)
	route.RegisterSwagger(server.router)
}

// Serve is a function that starts the server
func (server *Server) Serve() {
	err := server.router.Run(fmt.Sprintf(":%d", configs.GlobalConfig.Port))
	if err != nil {
		return
	}
}
