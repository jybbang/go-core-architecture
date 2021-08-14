package application

import (
	"fmt"
	"reflect"
	"time"

	cmap "github.com/orcaman/concurrent-map"
)

type Request interface{}
type RequestHandler func(interface{}) interface{}

type Notification interface{}
type NotificationHandler func(interface{})

type mediator struct {
	next                 Middlewarer
	requestHandlers      cmap.ConcurrentMap
	notificationHandlers cmap.ConcurrentMap
}

func (m *mediator) AddMiddleware(middleware Middlewarer) Middlewarer {
	m.next = middleware
	return m.next
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

func (m *mediator) Send(request Request) (interface{}, error) {
	valueOf := reflect.ValueOf(request)
	typeName := valueOf.Type().Name()

	handler, ok := m.requestHandlers.Get(typeName)
	if !ok {
		return nil, fmt.Errorf("handler not found exception")
	}

	handlerFn, ok := handler.(RequestHandler)
	if !ok {
		return nil, fmt.Errorf("handler not func exception")
	}

	defer timeMeasurement(time.Now(), typeName)

	if m.next != nil {
		return m.next.Next(request, handlerFn)
	} else {
		return handlerFn(request), nil
	}
}

func (m *mediator) Publish(notification Notification) error {
	valueOf := reflect.ValueOf(notification)
	typeName := valueOf.Type().Name()

	handler, ok := m.notificationHandlers.Get(typeName)
	if !ok {
		return fmt.Errorf("handler not found exception")
	}

	handlerFn, ok := handler.(NotificationHandler)
	if !ok {
		return fmt.Errorf("handler not func exception")
	}

	handlerFn(notification)
	return nil
}

func timeMeasurement(start time.Time, typeName string) {
	elapsed := time.Duration(time.Since(start))
	if elapsed > 500 {
		Log.Warn("long process time", typeName, elapsed, "ms")
	}
}
