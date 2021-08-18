package core

import (
	"time"

	"github.com/enriquebris/goconcurrentqueue"
	"github.com/reactivex/rxgo/v2"
	"github.com/sony/gobreaker"
	"gopkg.in/jeevatkm/go-model.v1"
)

// Builder Object for EventBus
type eventbusBuilder struct {
	mediator  *mediator
	messaging messagingAdapter
	cb        *gobreaker.CircuitBreaker
	setting   EventbusSettings
}

// Constructor for EventBusBuilder
func NewEventbusBuilder() *eventbusBuilder {
	o := new(eventbusBuilder)

	st := gobreaker.Settings{
		Name: "eventbus",
	}

	o.cb = gobreaker.NewCircuitBreaker(st)
	return o
}

// Builder method to set the field messaging in EventBusBuilder
func (b *eventbusBuilder) Settings(settings EventbusSettings) *eventbusBuilder {
	s := &EventbusSettings{
		BufferedEventBufferCount: 1000,
		BufferedEventBufferTime:  time.Duration(1 * time.Second),
		BufferedEventTimeout:     time.Duration(5 * time.Second),
	}

	err := model.Copy(s, settings)
	if err != nil {
		panic(err)
	}

	b.setting = *s
	return b
}

// Builder method to set the field messaging in EventBusBuilder
func (b *eventbusBuilder) CustomMediator(mediator *mediator) *eventbusBuilder {
	if mediator == nil {
		panic("mediator is required")
	}

	b.mediator = mediator
	return b
}

// Builder method to set the field messaging in EventBusBuilder
func (b *eventbusBuilder) MessaingAdapter(adapter messagingAdapter) *eventbusBuilder {
	if adapter == nil {
		panic("adapter is required")
	}

	b.messaging = adapter
	return b
}

// Builder method to set the field messaging in EventBusBuilder
func (b *eventbusBuilder) CircuitBreaker(setting CircuitBreakerSettings) *eventbusBuilder {
	b.cb = gobreaker.NewCircuitBreaker(setting.ToGobreakerSettings(b.cb.Name()))
	return b
}

// Build Method which creates EventBus
func (b *eventbusBuilder) Create() *eventbus {
	if b.messaging == nil {
		panic("messaging adapter is required")
	}
	if b.mediator == nil {
		b.mediator = GetMediator()
	}

	instance := &eventbus{
		mediator:     b.mediator,
		domainEvents: goconcurrentqueue.NewFIFO(),
		ch:           make(chan rxgo.Item, 1),
		messaging:    b.messaging,
		cb:           b.cb,
		settings:     b.setting,
	}
	instance.initialize()

	return instance
}

// Build Method which creates EventBus
func (b *eventbusBuilder) Build() *eventbus {
	if eventBusInstance != nil {
		panic("eventbus already created")
	}

	eventBusInstance = b.Create()

	return eventBusInstance
}
