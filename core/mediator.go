package core

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

type Mediator struct {
	middleware           Middlewarer
	requestHandlers      cmap.ConcurrentMap
	notificationHandlers cmap.ConcurrentMap
}

func (m *Mediator) Setup() *Mediator {
	return m
}

func (m *Mediator) AddMiddleware(middleware Middlewarer) Middlewarer {
	m.middleware = middleware
	return m.middleware
}

func (m *Mediator) AddHandler(request Request, handler RequestHandler) *Mediator {
	valueOf := reflect.ValueOf(request)
	typeName := valueOf.Type().Name()

	m.requestHandlers.Set(typeName, handler)
	return m
}

func (m *Mediator) AddNotificationHandler(notification Notification, handler NotificationHandler) *Mediator {
	valueOf := reflect.ValueOf(notification)
	typeName := valueOf.Type().Name()

	m.notificationHandlers.Set(typeName, handler)
	return m
}

func (m *Mediator) Send(request Request) (interface{}, error) {
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

	return m.next(request, handlerFn)
}

func (m *Mediator) Publish(notification Notification) error {
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

func (m *Mediator) next(request Request, handler RequestHandler) (interface{}, error) {
	if m.middleware != nil {
		ok, err := m.middleware.Run(request)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, fmt.Errorf("middleware block this request")
		}
		return m.middleware.Next(request, handler)
	} else {
		return handler(request), nil
	}
}

func timeMeasurement(start time.Time, typeName string) {
	elapsed := time.Since(start)
	if elapsed > time.Duration(500*time.Millisecond) {
		Log.Warn("long process time - ", typeName, elapsed)
	}
}
