package core

import (
	"reflect"

	cmap "github.com/orcaman/concurrent-map"
	"go.uber.org/zap"
)

// Builder Object for Mediator
type mediatorBuilder struct {
	requestHandlers      cmap.ConcurrentMap
	notificationHandlers cmap.ConcurrentMap
	log                  *zap.Logger
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

func (b *mediatorBuilder) AddPerformanceMeasure(logger *zap.Logger) *mediatorBuilder {
	b.log = logger
	return b
}

func (b *mediatorBuilder) AddHandler(request Request, handler RequestHandler) *mediatorBuilder {
	typeOf := reflect.TypeOf(request)
	typeName := typeOf.Elem().Name()

	if typeName == "" {
		panic("typeName is required")
	}

	b.requestHandlers.Set(typeName, handler)
	return b
}

func (b *mediatorBuilder) AddNotificationHandler(notification Notification, handler NotificationHandler) *mediatorBuilder {
	typeOf := reflect.TypeOf(notification)
	typeName := typeOf.Elem().Name()

	if typeName == "" {
		panic("typeName is required")
	}

	b.notificationHandlers.Set(typeName, handler)
	return b
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
		log:                  b.log,
	}
	instance.initialize()

	return instance
}
