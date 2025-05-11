package grpc

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/ngdangkietswe/swe-go-common-shared/config"
	gen "github.com/ngdangkietswe/swe-protobuf-shared/generated/workflow"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func RegisterWorkflowGrpcHandler(mux *runtime.ServeMux) {
	grpcWorkflowAddress := fmt.Sprintf("%s:%d",
		config.GetString("GRPC_WORKFLOW_HOST", "localhost"),
		config.GetInt("GRPC_WORKFLOW_PORT", 7060),
	)

	dialOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	services := []struct {
		registerFunc func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
		name         string
	}{
		{gen.RegisterWorkflowServiceHandlerFromEndpoint, "WorkflowService"},
	}

	for _, service := range services {
		if err := service.registerFunc(context.Background(), mux, grpcWorkflowAddress, dialOptions); err != nil {
			log.Fatalf("Failed to register %s: %v", service.name, err)
		}
	}
}
