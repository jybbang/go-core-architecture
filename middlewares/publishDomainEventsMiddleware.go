package middlewares

import (
	"context"

	"github.com/jybbang/go-core-architecture/core"
)

type publishDomainEventsMiddleware struct {
	core.Middleware
}

func NewPublishDomainEventsMiddleware() *publishDomainEventsMiddleware {
	return &publishDomainEventsMiddleware{}
}

func (m *publishDomainEventsMiddleware) Run(ctx context.Context, request core.Request) core.Result {
	result := m.Next()
	core.GetEventBus().PublishDomainEvents(ctx)
	return result
}
