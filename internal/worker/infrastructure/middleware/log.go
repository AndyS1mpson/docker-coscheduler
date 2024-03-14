package middleware

import (
	"context"

	"google.golang.org/grpc"

	"github.com/AndyS1mpson/docker-coscheduler/internal/utils/log"
)

// MiddlewareLog логирование ошибок запросов
func MiddlewareLog(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	h, err := handler(ctx, req)
	if err != nil {
		log.Error(err, log.Data{})
	}

	return h, err
}
