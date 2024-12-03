package grpc

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/ngdangkietswe/swe-go-common-shared/config"
	gen "github.com/ngdangkietswe/swe-protobuf-shared/generated/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func RegisterStorageGrpcHandler(mux *runtime.ServeMux) {
	err := gen.RegisterStorageServiceHandlerFromEndpoint(
		context.Background(),
		mux,
		fmt.Sprintf("%s:%d", config.GetString("GRPC_STORAGE_HOST", "localhost"), config.GetInt("GRPC_STORAGE_PORT", 7030)),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)

	if err != nil {
		log.Fatalf("Can't register Storage GRPC Handler: %v", err)
	}
}
