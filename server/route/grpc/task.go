package grpc

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/ngdangkietswe/swe-go-common-shared/config"
	gen "github.com/ngdangkietswe/swe-protobuf-shared/generated/task"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func RegisterTaskGrpcHandler(mux *runtime.ServeMux) {
	grpcTaskAddress := fmt.Sprintf("%s:%d",
		config.GetString("GRPC_TASK_HOST", "localhost"),
		config.GetInt("GRPC_TASK_PORT", 7010),
	)

	dialOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	services := []struct {
		registerFunc func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
		name         string
	}{
		{gen.RegisterTaskServiceHandlerFromEndpoint, "TaskService"},
		{gen.RegisterCommentServiceHandlerFromEndpoint, "CommentService"},
	}

	for _, service := range services {
		if err := service.registerFunc(context.Background(), mux, grpcTaskAddress, dialOptions); err != nil {
			log.Fatalf("Failed to register %s: %v", service.name, err)
		}
	}
}
