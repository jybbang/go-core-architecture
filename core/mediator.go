package core

import (
	"context"
	"reflect"
	"sync/atomic"
	"time"

	cmap "github.com/orcaman/concurrent-map"
	"go.uber.org/zap"
)

type mediator struct {
	Middleware
	requestHandlers      cmap.ConcurrentMap
	notificationHandlers cmap.ConcurrentMap
	log                  *zap.Logger
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

	if m.log != nil {
		defer m.timeMeasurement(time.Now(), typeName)
	}

	item, ok := m.requestHandlers.Get(typeName)
	if !ok {
		panic("request handler not found, you should register handler before use it")
	}

	handler := item.(RequestHandler)

	result := m.nextRun(ctx, request, handler)

	if result.E != nil {
		return result
	}

	GetEventbus().PublishDomainEvents(ctx)
	atomic.AddUint32(&m.sentCount, 1)

	return result
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

func (m *mediator) timeMeasurement(start time.Time, typeName string) {
	elapsed := time.Since(start)
	if elapsed > time.Duration(500*time.Millisecond) {
		m.log.Warn("send request long running", zap.String("request", typeName), zap.Duration("measure", elapsed))
	}
}
