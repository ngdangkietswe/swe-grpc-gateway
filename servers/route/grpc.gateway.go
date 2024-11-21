package route

import (
	"context"
	"github.com/felixge/httpsnoop"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/ngdangkietswe/swe-gateway-service/servers/route/grpc"
	"log"
	"net/http"
)

// WithLogger is a middleware to log request and response. This will log HTTP status code, duration and request path
func WithLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		m := httpsnoop.CaptureMetrics(handler, writer, request)
		log.Printf("http[%d]-- %s -- %s\n", m.Code, m.Duration, request.URL.Path)
	})
}

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
		}))

	grpc.RegisterTaskGrpcHandler(mux)

	router.Any("/api/v1/*any", gin.WrapH(WithLogger(mux)))
}
