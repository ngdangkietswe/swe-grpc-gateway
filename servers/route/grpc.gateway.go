package route

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/ngdangkietswe/swe-gateway-service/servers/route/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
)

// RegisterGrpcGateway register gRPC gateway. This will multiplex or route request different gRPC service
func RegisterGrpcGateway(router *gin.Engine) {
	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, writer http.ResponseWriter, request *http.Request, err error) {
			// Creating custom error type to add HTTP status code
			newError := runtime.HTTPStatusError{
				HTTPStatus: 400,
				Err:        err,
			}
			// Using default HTTP error handler to write error response
			runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, writer, request, &newError)
		}),
		// Using JSONPb marshaler to marshal and unmarshal JSON
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
			Marshaler: &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					EmitUnpopulated: true,
					UseEnumNumbers:  true,
					UseProtoNames:   true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			},
		}))

	// Register grpc handler
	grpc.RegisterTaskGrpcHandler(mux)
	grpc.RegisterAuthGrpcHandler(mux)

	router.Any("/api/v1/*any", gin.WrapH(mux))
}
