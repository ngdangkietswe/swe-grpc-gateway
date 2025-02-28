package grpc

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/ngdangkietswe/swe-go-common-shared/config"
	gen "github.com/ngdangkietswe/swe-protobuf-shared/generated/integration"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func RegisterIntegrationGrpcHandler(mux *runtime.ServeMux) {
	grpcIntegrationAddress := fmt.Sprintf("%s:%d",
		config.GetString("GRPC_INTEGRATION_HOST", "localhost"),
		config.GetInt("GRPC_INTEGRATION_PORT", 7040),
	)

	dialOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	services := []struct {
		registerFunc func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
		name         string
	}{
		{gen.RegisterStravaServiceHandlerFromEndpoint, "StravaService"},
	}

	for _, service := range services {
		if err := service.registerFunc(context.Background(), mux, grpcIntegrationAddress, dialOptions); err != nil {
			log.Fatalf("Failed to register %s: %v", service.name, err)
		}
	}
}
