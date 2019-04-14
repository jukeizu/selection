package startup

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

func LoggingInterceptor(logger zerolog.Logger) grpc.ServerOption {
	return grpc.UnaryInterceptor(makeLoggingInterceptor(logger))
}

func makeLoggingInterceptor(logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		begin := time.Now()

		resp, err := handler(ctx, req)

		logger := logger.With().
			Str("method", info.FullMethod).
			Str("took", time.Since(begin).String()).
			Logger()

		if err != nil {
			logger.Error().Err(err).Msg("")
			return resp, err
		}

		logger.Info().Msg("called")

		return resp, err
	}
}
