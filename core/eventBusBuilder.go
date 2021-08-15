package core

import (
	"github.com/reactivex/rxgo/v2"
	"github.com/sony/gobreaker"
)

// Builder Object for EventBus
type eventbusBuilder struct {
	messaging messagingAdapter
	cb        *gobreaker.CircuitBreaker
}

// Constructor for EventBusBuilder
func NewEventbusBuilder() *eventbusBuilder {
	o := new(eventbusBuilder)

	st := gobreaker.Settings{
		Name:          "eventbus",
		OnStateChange: onCbStateChange,
		Timeout:       cbDefaultTimeout,
		MaxRequests:   cbDefaultAllowedRequests,
	}
	o.cb = gobreaker.NewCircuitBreaker(st)

	return o
}

// Build Method which creates EventBus
func (b *eventbusBuilder) Build() *eventbus {
	if eventBusInstance != nil {
		panic("eventbus already created")
	}

	eventBusInstance = &eventbus{
		mediator:     GetMediator(),
		domainEvents: make([]DomainEventer, 0),
		ch:           make(chan rxgo.Item, 1),
		messaging:    b.messaging,
		cb:           b.cb,
	}
	eventBusInstance.initialize()

	return eventBusInstance
}

// Builder method to set the field messaging in EventBusBuilder
func (b *eventbusBuilder) MessaingAdapter(adapter messagingAdapter) *eventbusBuilder {
	b.messaging = adapter
	return b
}

// Builder method to set the field messaging in EventBusBuilder
func (b *eventbusBuilder) CircuitBreaker(setting gobreaker.Settings) *eventbusBuilder {
	setting.Name = b.cb.Name()
	if setting.OnStateChange == nil {
		setting.OnStateChange = onCbStateChange
	}
	b.cb = gobreaker.NewCircuitBreaker(setting)
	return b
}
