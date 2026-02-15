package middleware

import (
	"context"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RateLimiterInterceptor struct {
	limiter *rate.Limiter
}

func NewRateLimiterInterceptor(r rate.Limit, b int) *RateLimiterInterceptor {
	return &RateLimiterInterceptor{
		limiter: rate.NewLimiter(r, b),
	}
}

func (i *RateLimiterInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if !i.limiter.Allow() {
			return nil, status.Errorf(codes.ResourceExhausted, "cok fazla istek gonderdiniz, lutfen bekleyin")
		}
		return handler(ctx, req)
	}
}

func (i *RateLimiterInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if !i.limiter.Allow() {
			return status.Errorf(codes.ResourceExhausted, "cok fazla istek gonderdiniz")
		}
		return handler(srv, ss)
	}
}
