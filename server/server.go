package server

import (
	"fmt"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ngdangkietswe/swe-gateway-service/server/middleware"
	"github.com/ngdangkietswe/swe-gateway-service/server/route"
	"github.com/ngdangkietswe/swe-go-common-shared/cache"
	"github.com/ngdangkietswe/swe-go-common-shared/config"
	"github.com/ngdangkietswe/swe-go-common-shared/logger"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

type Server struct {
	router *gin.Engine
	logger *logger.Logger
}

func NewServer(logger *logger.Logger) *Server {
	return &Server{
		logger: logger,
		router: gin.Default(),
	}
}

// authConnInstance is a function that initializes the auth connection
func (s *Server) authConnInstance() (*grpc.ClientConn, error) {
	dialOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	return grpc.NewClient(
		fmt.Sprintf("%s:%d",
			config.GetString("GRPC_AUTH_HOST", "localhost"),
			config.GetInt("GRPC_AUTH_PORT", 7020)),
		dialOptions...)
}

// Init is a function that initializes the server
func (s *Server) Init() {
	// Init auth grpc connection
	authConn, err := s.authConnInstance()
	if err != nil {
		log.Fatalf("failed to initialize auth connection: %v", err)
	}

	redisCache := cache.NewRedisCache(cache.WithTimeout(3 * time.Second))

	s.router.Use(
		gin.Recovery(),
		requestid.New(
			requestid.WithGenerator(func() string {
				id, _ := uuid.NewV6()
				return id.String()
			}),
		),
	)

	s.router.Use(
		middleware.NewCORSMiddleware(s.logger).AsMiddleware(),
		middleware.NewRequestTimeoutMiddleware(s.logger).AsMiddleware(),
		middleware.NewRequestLoggingMiddleware(s.logger).AsMiddleware(),
		middleware.NewRateLimitMiddleware(s.logger).AsMiddleware(),
		middleware.NewAuthMiddleware(s.logger, authConn, redisCache).AsMiddleware(),
	)

	route.RegisterGrpcGateway(s.router)
	route.RegisterSwagger(s.router)
}

// Serve is a function that starts the server
func (s *Server) Serve() {
	err := s.router.Run(fmt.Sprintf(":%d", config.GetInt("PORT", 7777)))
	if err != nil {
		return
	}
}

var Module = fx.Provide(NewServer)
