package middleware

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
)

var (
	opsProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "notification_service_requests_total",
		Help: "The total number of processed requests",
	}, []string{"method"})
)

func MetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	resp, err := handler(ctx, req)

	opsProcessed.WithLabelValues(info.FullMethod).Inc()

	_ = start
	return resp, err
}
