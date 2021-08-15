package core

import (
	"context"
	"errors"
	"reflect"
	"time"

	cmap "github.com/orcaman/concurrent-map"
)

type mediator struct {
	Middleware
	requestHandlers      cmap.ConcurrentMap
	notificationHandlers cmap.ConcurrentMap
}

func (m *mediator) initialize() *mediator {
	return m
}

func (m *mediator) Send(ctx context.Context, request Request) Result {
	valueOf := reflect.ValueOf(request)
	typeName := valueOf.Type().Name()
	defer timeMeasurement(time.Now(), typeName)

	item, ok := m.requestHandlers.Get(typeName)
	if !ok {
		return Result{E: errors.New("request handler not found")}
	}
	handler := item.(RequestHandler)

	services := Services{
		Eventbus: GetEventbus(),
		States:   GetStateService(),
	}

	result := m.nextRun(ctx, services, request, handler)

	services.Eventbus.PublishDomainEvents(ctx)

	return result
}

func (m *mediator) Publish(ctx context.Context, notification Notification) error {
	valueOf := reflect.ValueOf(notification)
	typeName := valueOf.Type().Name()

	item, ok := m.notificationHandlers.Get(typeName)
	if !ok {
		return errors.New("request handler not found")
	}

	handler := item.(NotificationHandler)

	return handler(ctx, notification)
}

func timeMeasurement(start time.Time, typeName string) {
	elapsed := time.Since(start)
	if elapsed > time.Duration(500*time.Millisecond) {
		Log.Warnw("send request long running", "request", typeName, "measure", elapsed)
	}
}
