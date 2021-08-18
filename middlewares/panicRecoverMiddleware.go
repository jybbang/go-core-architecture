package middlewares

import (
	"context"

	"github.com/jybbang/go-core-architecture/core"
)

type panicRecoverMiddleware struct {
	core.Middleware
	callback func(interface{})
}

func NewPanicRecoverMiddleware(callback func(interface{})) *panicRecoverMiddleware {
	return &panicRecoverMiddleware{
		callback: callback,
	}
}

func (m *panicRecoverMiddleware) Run(ctx context.Context, request core.Request) core.Result {
	defer m.panicRecover()
	return m.Next()
}

func (m *panicRecoverMiddleware) panicRecover() {
	if r := recover(); r != nil {
		if m.callback != nil {
			m.callback(r)
		}
	}
}
