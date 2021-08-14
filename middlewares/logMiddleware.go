package middlewares

import (
	"github.com/jybbang/go-core-architecture/core"
)

type logMiddleware struct {
	core.Middleware
}

func NewLogMiddleware() *logMiddleware {
	return &logMiddleware{}
}

func (m *logMiddleware) Run(request core.Request) (bool, error) {
	core.Log.Info("mediator request", request)
	return true, nil
}
