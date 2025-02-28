package main

import (
	"github.com/ngdangkietswe/swe-gateway-service/logger"
	"github.com/ngdangkietswe/swe-gateway-service/server"
	"github.com/ngdangkietswe/swe-go-common-shared/config"
	"go.uber.org/fx"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config.Init()

	app := fx.New(
		logger.Module,
		server.Module,
		fx.Invoke(func(s *server.Server) {
			// Initialize the server
			s.Init()

			// Start the server in a goroutine
			go s.Serve()

			// Handle graceful shutdown
			stop := make(chan os.Signal, 1)
			signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
			<-stop
		}),
	)

	app.Run()
}
