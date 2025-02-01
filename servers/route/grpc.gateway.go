package route

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/ngdangkietswe/swe-gateway-service/servers/route/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"log"
	"net/http"
	"strings"
)

// RegisterGrpcGateway register gRPC gateway. This will multiplex or route request different gRPC service
func RegisterGrpcGateway(router *gin.Engine) {
	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, writer http.ResponseWriter, request *http.Request, err error) {
			newError := runtime.HTTPStatusError{
				HTTPStatus: 400,
				Err:        err,
			}
			runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, writer, request, &newError)
		}),
		runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
			if strings.HasPrefix(strings.ToLower(key), "grpc-") {
				return key, true
			}
			return "", false
		}),
		runtime.WithOutgoingHeaderMatcher(func(key string) (string, bool) {
			return key, true
		}),
		runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
			md := metadata.Pairs()
			for key, values := range req.Header {
				lowerKey := strings.ToLower(key)
				if lowerKey == "authorization" || lowerKey == "x-api-key" || strings.HasPrefix(lowerKey, "grpc-") {
					md.Append(lowerKey, values...)
				}
			}
			return md
		}),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
			Marshaler: &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					EmitUnpopulated: true,
					UseProtoNames:   true,
					UseEnumNumbers:  true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			},
		}),
		runtime.WithForwardResponseOption(func(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
			log.Printf("Response: %+v", resp)
			return nil
		}),
	)

	// Register grpc handler
	grpc.RegisterTaskGrpcHandler(mux)
	grpc.RegisterAuthGrpcHandler(mux)
	grpc.RegisterStorageGrpcHandler(mux)
	grpc.RegisterIntegrationGrpcHandler(mux)

	router.Any("/api/v1/*any", gin.WrapH(mux))
}
