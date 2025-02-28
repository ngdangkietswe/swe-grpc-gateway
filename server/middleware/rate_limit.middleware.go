package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ngdangkietswe/swe-go-common-shared/config"
	"github.com/ngdangkietswe/swe-go-common-shared/logger"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"go.uber.org/zap"
)

type RateLimitMiddleware struct {
	logger *logger.Logger
}

func (r RateLimitMiddleware) ShouldSkip(ctx *gin.Context) bool {
	//TODO implement me
	panic("implement me")
}

func (r RateLimitMiddleware) Handle(ctx *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (r RateLimitMiddleware) AsMiddleware() gin.HandlerFunc {
	rate, err := limiter.NewRateFromFormatted(config.GetString("RATE_LIMIT", "100-M")) // e.g., "100-H" for 100 requests per hour
	if err != nil {
		r.logger.Error("Invalid rate limit configuration: %v", zap.String("error", err.Error()))
	}

	store := memory.NewStore() // In-memory store; consider Redis for production
	instance := limiter.New(store, rate)

	return mgin.NewMiddleware(instance)
}

func NewRateLimitMiddleware(logger *logger.Logger) Middleware {
	return &RateLimitMiddleware{
		logger: logger,
	}
}
