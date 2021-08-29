package core

import (
	"fmt"
	"time"

	"github.com/enriquebris/goconcurrentqueue"
	"github.com/reactivex/rxgo/v2"
	"gopkg.in/jeevatkm/go-model.v1"
)

// Builder Object for EventBus
type eventBusBuilder struct {
	mediator   *mediator
	messaging  messagingAdapter
	cbSettings CircuitBreakerSettings
	setting    EventBusSettings
}

// Constructor for EventBusBuilder
func NewEventBusBuilder() *eventBusBuilder {
	o := new(eventBusBuilder)
	o.cbSettings = CircuitBreakerSettings{Name: "eventbus"}
	return o
}

// Builder method to set the field messaging in EventBusBuilder
func (b *eventBusBuilder) Settings(settings EventBusSettings) *eventBusBuilder {
	s := &EventBusSettings{
		BufferedEventBufferCount: 1000,
		BufferedEventBufferTime:  time.Duration(1 * time.Second),
		BufferedEventTimeout:     time.Duration(5 * time.Second),
	}

	err := model.Copy(s, settings)
	if err != nil {
		panic(fmt.Errorf("mapping errors occurred: %v", err))
	}

	b.setting = *s
	return b
}

// Builder method to set the field messaging in EventBusBuilder
func (b *eventBusBuilder) CustomMediator(mediator *mediator) *eventBusBuilder {
	if mediator == nil {
		panic("mediator is required")
	}

	b.mediator = mediator
	return b
}

// Builder method to set the field messaging in EventBusBuilder
func (b *eventBusBuilder) MessaingAdapter(adapter messagingAdapter) *eventBusBuilder {
	if adapter == nil {
		panic("adapter is required")
	}

	b.messaging = adapter
	return b
}

// Builder method to set the field messaging in EventBusBuilder
func (b *eventBusBuilder) CircuitBreaker(setting CircuitBreakerSettings) *eventBusBuilder {
	b.cbSettings = setting
	return b
}

// Build Method which creates EventBus
func (b *eventBusBuilder) Create() *eventBus {
	if b.messaging == nil {
		panic("messaging adapter is required")
	}
	if b.mediator == nil {
		b.mediator = GetMediator()
	}

	instance := &eventBus{
		mediator:     b.mediator,
		domainEvents: goconcurrentqueue.NewFIFO(),
		ch:           make(chan rxgo.Item, 1),
		messaging:    b.messaging,
		settings:     b.setting,
	}
	instance.initialize(b.cbSettings)

	return instance
}

// Build Method which creates EventBus
func (b *eventBusBuilder) Build() *eventBus {
	if eventBusInstance != nil {
		panic("eventBus already created")
	}

	eventBusInstance = b.Create()

	return eventBusInstance
}
