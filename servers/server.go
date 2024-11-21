package servers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ngdangkietswe/swe-gateway-service/configs"
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

func (server *Server) Init() {
	route.RegisterGrpcGateway(server.router)
}

func (server *Server) Serve() {
	err := server.router.Run(fmt.Sprintf(":%d", configs.GlobalConfig.Port))
	if err != nil {
		return
	}
}
