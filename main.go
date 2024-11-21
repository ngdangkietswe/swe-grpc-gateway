package main

import "github.com/ngdangkietswe/swe-gateway-service/servers"

func main() {
	ginServer := servers.NewServer()
	ginServer.Init()
	ginServer.Serve()
}
