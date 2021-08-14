package middlewares

import (
	"github.com/jybbang/go-core-architecture/application"
)

type logMiddleware struct {
	application.Middleware
}

func NewLogMiddleware() *logMiddleware {
	return &logMiddleware{}
}

func (m *logMiddleware) Run(request application.Request) (bool, error) {
	application.Log.Info("mediator log")
	return true, nil
}
