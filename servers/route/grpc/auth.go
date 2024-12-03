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
	err := gen.RegisterAuthServiceHandlerFromEndpoint(
		context.Background(),
		mux,
		fmt.Sprintf("%s:%d", config.GetString("GRPC_AUTH_HOST", "localhost"), config.GetInt("GRPC_AUTH_PORT", 7020)),
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)

	if err != nil {
		log.Fatalf("Can't register Auth GRPC Handler: %v", err)
	}
}
