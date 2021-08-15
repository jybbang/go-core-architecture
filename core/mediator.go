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

func (m *mediator) AddHandler(request Request, handler RequestHandler) *mediator {
	valueOf := reflect.ValueOf(request)
	typeName := valueOf.Type().Name()

	m.requestHandlers.Set(typeName, handler)
	return m
}

func (m *mediator) AddNotificationHandler(notification Notification, handler NotificationHandler) *mediator {
	valueOf := reflect.ValueOf(notification)
	typeName := valueOf.Type().Name()

	m.notificationHandlers.Set(typeName, handler)
	return m
}

func (m *mediator) Send(ctx context.Context, request Request) Result {
	valueOf := reflect.ValueOf(request)
	typeName := valueOf.Type().Name()

	item, ok := m.requestHandlers.Get(typeName)
	if !ok {
		return Result{E: errors.New("request handler not found")}
	}

	handler := item.(RequestHandler)

	defer timeMeasurement(time.Now(), typeName)

	return m.Next(ctx, request, handler)
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

func (m *mediator) Run(ctx context.Context, request Request) (ok bool, err error) {
	return true, nil
}

func timeMeasurement(start time.Time, typeName string) {
	elapsed := time.Since(start)
	if elapsed > time.Duration(500*time.Millisecond) {
		Log.Warnw("send request long running", "request", typeName, "measure", elapsed)
	}
}
