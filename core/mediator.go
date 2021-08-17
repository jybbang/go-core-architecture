package core

import (
	"context"
	"reflect"
	"sync/atomic"

	cmap "github.com/orcaman/concurrent-map"
)

type mediator struct {
	middleware           behavior
	requestHandlers      cmap.ConcurrentMap
	notificationHandlers cmap.ConcurrentMap
	sentCount            uint32
	publishedCount       uint32
}

func (m *mediator) initialize() *mediator {
	return m
}

func (m *mediator) GetSentCount() uint32 {
	return m.sentCount
}

func (m *mediator) GetPublishedCount() uint32 {
	return m.publishedCount
}

func (m *mediator) AddMiddleware(next behavior) behavior {
	m.middleware = next
	return m.middleware
}

func (m *mediator) Send(ctx context.Context, request Request) Result {
	if err := ctx.Err(); err != nil {
		return Result{E: err}
	}

	typeOf := reflect.TypeOf(request)
	typeName := typeOf.Elem().Name()

	if typeName == "" {
		panic("typeName is required")
	}

	if openTracer != nil {
		span := openTracer.StartSpan(typeName)
		defer span.Finish()
	}

	item, ok := m.requestHandlers.Get(typeName)
	if !ok {
		panic("request handler not found, you should register handler before use it")
	}

	handler := item.(RequestHandler)

	result := m.next(ctx, request, handler)

	if result.E != nil {
		return result
	}

	GetEventbus().PublishDomainEvents(ctx)
	atomic.AddUint32(&m.sentCount, 1)

	return result
}

func (m *mediator) next(ctx context.Context, request Request, handler RequestHandler) Result {
	if m.middleware != nil {
		m.middleware.setParameters(ctx, request, handler)
		return m.middleware.Run(ctx, request)
	} else {
		return handler(ctx, request)
	}
}

func (m *mediator) Publish(ctx context.Context, notification Notification) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	typeOf := reflect.TypeOf(notification)
	typeName := typeOf.Elem().Name()

	if typeName == "" {
		panic("typeName is required")
	}

	if openTracer != nil {
		span := openTracer.StartSpan(typeName)
		defer span.Finish()
	}

	item, ok := m.notificationHandlers.Get(typeName)
	if !ok {
		panic("request handler not found, you should register handler before use it")
	}

	handler := item.(NotificationHandler)

	err := handler(ctx, notification)
	if err != nil {
		return err
	}

	atomic.AddUint32(&m.publishedCount, 1)

	return nil
}
