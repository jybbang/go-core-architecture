package middlewares

import (
	"context"

	"github.com/jybbang/go-core-architecture/core"
	"go.uber.org/zap"
)

type logMiddleware struct {
	core.Middleware
	log *zap.Logger
}

func NewLogMiddleware(logger *zap.Logger) *logMiddleware {
	return &logMiddleware{
		log: logger,
	}
}

func (m *logMiddleware) Run(ctx context.Context, request core.Request) core.Result {
	m.log.Info("send request log", zap.Reflect("request", request))
	return m.Next()
}
