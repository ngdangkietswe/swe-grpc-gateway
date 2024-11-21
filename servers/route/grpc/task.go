package grpc

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/ngdangkietswe/swe-gateway-service/configs"
	gen "github.com/ngdangkietswe/swe-protobuf-shared/generated/task"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func RegisterTaskGrpcHandler(mux *runtime.ServeMux) {
	err := gen.RegisterTaskServiceHandlerFromEndpoint(
		context.Background(),
		mux,
		fmt.Sprintf("%s:%d", configs.GlobalConfig.GrpcTaskHost, configs.GlobalConfig.GrpcTaskPort),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)

	if err != nil {
		log.Fatalf("Can't register Task GRPC Handler: %v", err)
	}
}
