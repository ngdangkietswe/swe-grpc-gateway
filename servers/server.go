package servers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ngdangkietswe/swe-gateway-service/servers/middleware"
	"github.com/ngdangkietswe/swe-gateway-service/servers/route"
	"github.com/ngdangkietswe/swe-go-common-shared/cache"
	"github.com/ngdangkietswe/swe-go-common-shared/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
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
	// Init auth grpc connection
	dialOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	authConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d",
			config.GetString("GRPC_AUTH_HOST", "localhost"),
			config.GetInt("GRPC_AUTH_PORT", 7020)),
		dialOptions...)
	if err != nil {
		log.Fatal("Can't connect to Auth GRPC Server")
	}

	redisCache := cache.NewRedisCache(cache.WithTimeout(3 * time.Second))

	server.router.Use(middleware.NewAuthMiddleware(authConn, redisCache).AsMiddleware())
	route.RegisterGrpcGateway(server.router)
	route.RegisterSwagger(server.router)
}

// Serve is a function that starts the server
func (server *Server) Serve() {
	err := server.router.Run(fmt.Sprintf(":%d", config.GetInt("PORT", 7777)))
	if err != nil {
		return
	}
}
