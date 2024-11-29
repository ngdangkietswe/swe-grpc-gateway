package grpc

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/ngdangkietswe/swe-gateway-service/configs"
	gen "github.com/ngdangkietswe/swe-protobuf-shared/generated/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func RegisterStorageGrpcHandler(mux *runtime.ServeMux) {
	err := gen.RegisterStorageServiceHandlerFromEndpoint(
		context.Background(),
		mux,
		fmt.Sprintf("%s:%d", configs.GlobalConfig.GrpcStorageHost, configs.GlobalConfig.GrpcStoragePort),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)

	if err != nil {
		log.Fatalf("Can't register Storage GRPC Handler: %v", err)
	}
}
