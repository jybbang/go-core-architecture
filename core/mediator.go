package core

import (
	"context"
	"errors"
	"reflect"
	"time"

	cmap "github.com/orcaman/concurrent-map"
	"go.uber.org/zap"
)

type mediator struct {
	Middleware
	requestHandlers      cmap.ConcurrentMap
	notificationHandlers cmap.ConcurrentMap
	log                  *zap.Logger
}

func (m *mediator) initialize() *mediator {
	return m
}

func (m *mediator) Send(ctx context.Context, request Request) Result {
	valueOf := reflect.ValueOf(request)
	typeName := valueOf.Type().Name()

	if m.log != nil {
		defer m.timeMeasurement(time.Now(), typeName)
	}

	if openTracer != nil {
		span := openTracer.StartSpan(typeName)
		defer span.Finish()
	}

	item, ok := m.requestHandlers.Get(typeName)
	if !ok {
		return Result{E: errors.New("request handler not found")}
	}
	handler := item.(RequestHandler)

	eventbus := GetEventbus()

	result := m.nextRun(ctx, request, handler)

	if result.E != nil {
		return result
	}

	eventbus.PublishDomainEvents(ctx)

	return result
}

func (m *mediator) Publish(ctx context.Context, notification Notification) error {
	valueOf := reflect.ValueOf(notification)
	typeName := valueOf.Type().Name()

	if openTracer != nil {
		span := openTracer.StartSpan(typeName)
		defer span.Finish()
	}

	item, ok := m.notificationHandlers.Get(typeName)
	if !ok {
		return errors.New("request handler not found")
	}

	handler := item.(NotificationHandler)

	return handler(ctx, notification)
}

func (m *mediator) timeMeasurement(start time.Time, typeName string) {
	elapsed := time.Since(start)
	if elapsed > time.Duration(500*time.Millisecond) {
		m.log.Warn("send request long running", zap.String("request", typeName), zap.Duration("measure", elapsed))
	}
}
