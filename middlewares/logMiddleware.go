package middlewares

import (
	"context"

	"github.com/jybbang/go-core-architecture/core"
)

type logMiddleware struct {
	core.Middleware
}

func NewLogMiddleware() *logMiddleware {
	return &logMiddleware{}
}

func (m *logMiddleware) Run(ctx context.Context, request core.Request) (ok bool, err error) {
	core.Log.Infow("mediator request log", "request", request)
	return true, nil
}
