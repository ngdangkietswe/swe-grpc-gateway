package grpc

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/ngdangkietswe/swe-go-common-shared/config"
	gen "github.com/ngdangkietswe/swe-protobuf-shared/generated/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func RegisterAuthGrpcHandler(mux *runtime.ServeMux) {
	grpcAuthAddress := fmt.Sprintf("%s:%d",
		config.GetString("GRPC_AUTH_HOST", "localhost"),
		config.GetInt("GRPC_AUTH_PORT", 7020),
	)

	dialOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	services := []struct {
		registerFunc func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
		name         string
	}{
		{gen.RegisterAuthServiceHandlerFromEndpoint, "AuthService"},
		{gen.RegisterPermissionServiceHandlerFromEndpoint, "PermissionService"},
	}

	for _, service := range services {
		if err := service.registerFunc(context.Background(), mux, grpcAuthAddress, dialOptions); err != nil {
			log.Fatalf("Failed to register %s: %v", service.name, err)
		}
	}
}
