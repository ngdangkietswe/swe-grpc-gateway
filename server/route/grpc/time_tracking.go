package grpc

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/ngdangkietswe/swe-go-common-shared/config"
	gen "github.com/ngdangkietswe/swe-protobuf-shared/generated/timetracking"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func RegisterTimeTrackingGrpcHandler(mux *runtime.ServeMux) {
	grpcTimeTrackingAddress := fmt.Sprintf("%s:%d",
		config.GetString("GRPC_TIME_TRACKING_HOST", "localhost"),
		config.GetInt("GRPC_TIME_TRACKING_PORT", 7050),
	)

	dialOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	services := []struct {
		registerFunc func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
		name         string
	}{
		{gen.RegisterTimeTrackingServiceHandlerFromEndpoint, "TimeTrackingService"},
	}

	for _, service := range services {
		if err := service.registerFunc(context.Background(), mux, grpcTimeTrackingAddress, dialOptions); err != nil {
			log.Fatalf("Failed to register %s: %v", service.name, err)
		}
	}
}
