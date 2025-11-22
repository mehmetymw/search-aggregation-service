package grpc

import (
	"context"

	"github.com/mehmetymw/search-aggregation-service/backend/domain/entity"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RateLimitInterceptor struct {
	limiter *rate.Limiter
}

func NewRateLimitInterceptor(config entity.RateLimitConfig) *RateLimitInterceptor {
	return &RateLimitInterceptor{
		limiter: rate.NewLimiter(rate.Limit(config.RPS), config.Burst),
	}
}

func (i *RateLimitInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if !i.limiter.Allow() {
			return nil, status.Error(codes.ResourceExhausted, "rate limit exceeded")
		}
		return handler(ctx, req)
	}
}
