package core

import (
	"github.com/reactivex/rxgo/v2"
	"github.com/sony/gobreaker"
)

// Builder Object for EventBus
type eventBusBuilder struct {
	messaging messagingAdapter
	cb        *gobreaker.CircuitBreaker
}

// Constructor for EventBusBuilder
func NewEventBusBuilder() *eventBusBuilder {
	o := new(eventBusBuilder)

	st := gobreaker.Settings{
		Name:          "eventbus",
		OnStateChange: OnCbStateChange,
	}
	o.cb = gobreaker.NewCircuitBreaker(st)

	return o
}

// Build Method which creates EventBus
func (b *eventBusBuilder) Build() *eventbus {
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
func (b *eventBusBuilder) MessaingAdapter(adapter messagingAdapter) *eventBusBuilder {
	b.messaging = adapter
	return b
}

// Builder method to set the field messaging in EventBusBuilder
func (b *eventBusBuilder) CircuitBreaker(setting gobreaker.Settings) *eventBusBuilder {
	setting.Name = b.cb.Name()
	setting.OnStateChange = OnCbStateChange
	b.cb = gobreaker.NewCircuitBreaker(setting)
	return b
}
