package middlewares

import (
	"context"

	"github.com/jybbang/go-core-architecture/core"
	"go.uber.org/zap"
)

type zapLogMiddleware struct {
	core.Middleware
	log *zap.Logger
}

func NewZapLogMiddleware(logger *zap.Logger) *zapLogMiddleware {
	return &zapLogMiddleware{
		log: logger,
	}
}

func (m *zapLogMiddleware) Run(ctx context.Context, request core.Request) (ok bool, err error) {
	m.log.Info("mediator request log",
		zap.Reflect("request", request))
	return true, nil
}
