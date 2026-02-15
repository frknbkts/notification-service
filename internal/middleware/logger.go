package middleware

import (
	"context"
	"log/slog"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func LoggerInterceptor() grpc.UnaryServerInterceptor {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start)

		code := status.Code(err)

		level := slog.LevelInfo
		msg := "GRPC Request"
		if code != codes.OK {
			level = slog.LevelError
			msg = "GRPC Request Failed"
		}

		logger.Log(ctx, level, msg,
			slog.String("method", info.FullMethod),
			slog.String("duration", duration.String()),
			slog.String("code", code.String()),
			slog.Any("error", err),
		)

		return resp, err
	}
}
