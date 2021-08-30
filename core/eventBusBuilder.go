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
	settings   EventBusSettings
}

// Constructor for EventBusBuilder
func NewEventBusBuilder() *eventBusBuilder {
	o := new(eventBusBuilder)
	o.cbSettings = CircuitBreakerSettings{
		AllowedRequestInHalfOpen: 1,
		DurationOfBreak:          time.Duration(60 * time.Second),
		SamplingDuration:         time.Duration(60 * time.Second),
		SamplingFailureCount:     5,
	}
	o.settings = EventBusSettings{
		BufferedEventBufferCount: 1000,
		BufferedEventBufferTime:  time.Duration(1 * time.Second),
		BufferedEventTimeout:     time.Duration(5 * time.Second),
		ConnectionTimeout:        time.Duration(10 * time.Second),
	}

	return o
}

// Build Method which creates EventBus
func (b *eventBusBuilder) Build() *eventBus {
	if eventBusInstance != nil {
		panic("eventBus already created")
	}

	eventBusInstance = b.Create()

	return eventBusInstance
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
		settings:     b.settings,
	}

	instance.cb = b.cbSettings.ToCircuitBreaker("eventbus", instance.onCircuitOpen)

	instance.initialize()

	return instance
}

// Builder method to set the field messaging in EventBusBuilder
func (b *eventBusBuilder) Settings(settings EventBusSettings) *eventBusBuilder {
	err := model.Copy(&b.settings, settings)

	if err != nil {
		panic(fmt.Errorf("settings mapping errors occurred: %v", err))
	}

	return b
}

// Builder method to set the field messaging in EventBusBuilder
func (b *eventBusBuilder) CircuitBreaker(settings CircuitBreakerSettings) *eventBusBuilder {
	err := model.Copy(&b.cbSettings, settings)

	if err != nil {
		panic(fmt.Errorf("cb settings mapping errors occurred: %v", err))
	}

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
