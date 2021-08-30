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

	// default
	o.AddNotificationHandler(new(bufferedEvent), bufferedEventHandler)

	return o
}

// Build Method which creates Mediator
func (b *mediatorBuilder) Build() *mediator {
	if mediatorInstance != nil {
		panic("mediator already created")
	}

	mediatorInstance = b.Create()

	return mediatorInstance
}

// Build Method which creates Mediator
func (b *mediatorBuilder) Create() *mediator {
	instance := &mediator{
		requestHandlers:      b.requestHandlers,
		notificationHandlers: b.notificationHandlers,
	}

	instance.initialize()

	return instance
}

func (b *mediatorBuilder) AddHandler(request Request, handler RequestHandler) *mediatorBuilder {
	if request == nil {
		panic("request is required")
	}

	if handler == nil {
		panic("handler is required")
	}

	typeOf := reflect.TypeOf(request)

	typeName := typeOf.Elem().Name()

	b.requestHandlers.Set(typeName, handler)

	return b
}

func (b *mediatorBuilder) AddNotificationHandler(notification Notification, handler NotificationHandler) *mediatorBuilder {
	if notification == nil {
		panic("notification is required")
	}

	if handler == nil {
		panic("notification handler is required")
	}

	typeOf := reflect.TypeOf(notification)

	typeName := typeOf.Elem().Name()

	b.notificationHandlers.Set(typeName, handler)

	return b
}
