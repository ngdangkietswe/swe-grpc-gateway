package main

import (
	"github.com/ngdangkietswe/swe-gateway-service/servers"
	"github.com/ngdangkietswe/swe-go-common-shared/config"
)

func main() {
	// Initialize the configuration
	config.Init()
	// Create and start the server
	ginServer := servers.NewServer()
	ginServer.Init()
	ginServer.Serve()
}
