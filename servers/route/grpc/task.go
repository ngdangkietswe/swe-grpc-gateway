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
	err := gen.RegisterTaskServiceHandlerFromEndpoint(
		context.Background(),
		mux,
		fmt.Sprintf("%s:%d", config.GetString("GRPC_TASK_HOST", "localhost"), config.GetInt("GRPC_TASK_PORT", 7010)),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)

	if err != nil {
		log.Fatalf("Can't register Task GRPC Handler: %v", err)
	}
}
