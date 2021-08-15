package core

import (
	"reflect"

	cmap "github.com/orcaman/concurrent-map"
)

// Builder Object for Mediator
type mediatorBuilder struct {
	requestHandlers      cmap.ConcurrentMap
	notificationHandlers cmap.ConcurrentMap
}

// Constructor for MediatorBuilder
func NewMediatorBuilder() *mediatorBuilder {
	o := new(mediatorBuilder)
	o.requestHandlers = cmap.New()
	o.notificationHandlers = cmap.New()
	return o
}

func (b *mediatorBuilder) AddHandler(request Request, handler RequestHandler) *mediatorBuilder {
	valueOf := reflect.ValueOf(request)
	typeName := valueOf.Type().Name()

	b.requestHandlers.Set(typeName, handler)
	return b
}

func (b *mediatorBuilder) AddNotificationHandler(notification Notification, handler NotificationHandler) *mediatorBuilder {
	valueOf := reflect.ValueOf(notification)
	typeName := valueOf.Type().Name()

	b.notificationHandlers.Set(typeName, handler)
	return b
}

// Build Method which creates Mediator
func (b *mediatorBuilder) Build() *mediator {
	if mediatorInstance != nil {
		panic("mediator already created")
	}

	mediatorInstance = &mediator{
		requestHandlers:      b.requestHandlers,
		notificationHandlers: b.notificationHandlers,
	}
	mediatorInstance.initialize()

	return mediatorInstance
}
