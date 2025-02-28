package middleware

import (
	"fmt"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/ngdangkietswe/swe-go-common-shared/logger"
)

type RequestLoggingMiddleware struct {
	logger *logger.Logger
}

func (r *RequestLoggingMiddleware) ShouldSkip(ctx *gin.Context) bool {
	//TODO implement me
	panic("implement me")
}

func (r *RequestLoggingMiddleware) Handle(ctx *gin.Context) {
	start := time.Now()
	requestID := requestid.Get(ctx)

	// Capture additional request details
	clientIP := ctx.ClientIP()
	userAgent := ctx.Request.UserAgent()
	referer := ctx.Request.Referer()

	// Proceed with the request
	ctx.Next()

	// Calculate latency and status after the request
	latency := time.Since(start)
	status := ctx.Writer.Status()

	// Fancy log message without statusSymbol
	r.logger.Info(fmt.Sprintf("[%s] %s %s → %d (%s) | Latency: %s | IP: %s | UA: %s | Referer: %s",
		requestID,
		ctx.Request.Method,
		ctx.Request.URL.Path,
		status,
		httpStatusText(status), // Human-readable status
		formatLatency(latency), // Prettify latency
		clientIP,
		userAgent,
		referer),
	)
}

// Helper function to get human-readable HTTP status text
func httpStatusText(status int) string {
	switch status {
	case 200:
		return "OK"
	case 201:
		return "Created"
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized"
	case 403:
		return "Forbidden"
	case 404:
		return "Not Found"
	case 500:
		return "Internal Server Error"
	case 502:
		return "Bad Gateway"
	case 504:
		return "Gateway Timeout"
	default:
		return fmt.Sprintf("Status %d", status)
	}
}

// Helper function to format latency in a more readable way
func formatLatency(latency time.Duration) string {
	switch {
	case latency < time.Millisecond:
		return fmt.Sprintf("%.3fµs", float64(latency.Microseconds()))
	case latency < time.Second:
		return fmt.Sprintf("%.3fms", float64(latency.Milliseconds()))
	default:
		return fmt.Sprintf("%.3fs", latency.Seconds())
	}
}

func (r *RequestLoggingMiddleware) AsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		r.Handle(c)
	}
}

func NewRequestLoggingMiddleware(logger *logger.Logger) Middleware {
	return &RequestLoggingMiddleware{
		logger: logger,
	}
}
