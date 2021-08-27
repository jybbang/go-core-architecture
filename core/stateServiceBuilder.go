package core

import "github.com/sony/gobreaker"

// Builder Object for StateService
type stateServiceBuilder struct {
	state stateAdapter
	cb    *gobreaker.CircuitBreaker
}

// Constructor for StateServiceBuilder
func NewStateServiceBuilder() *stateServiceBuilder {
	o := new(stateServiceBuilder)

	st := gobreaker.Settings{
		Name: "state service",
	}
	o.cb = gobreaker.NewCircuitBreaker(st)

	return o
}

// Build Method which creates StateService
func (b *stateServiceBuilder) Build() *stateService {
	if statesInstance != nil {
		panic("state service already created")
	}

	statesInstance = b.Create()

	return statesInstance
}

// Build Method which creates EventBus
func (b *stateServiceBuilder) Create() *stateService {
	if b.state == nil {
		panic("state adapter is required")
	}

	instance := &stateService{
		state: b.state,
		cb:    b.cb,
	}
	instance.initialize()

	return instance
}

// Builder method to set the field state in StateServiceBuilder
func (b *stateServiceBuilder) StateAdapter(adapter stateAdapter) *stateServiceBuilder {
	if adapter == nil {
		panic("adapter is required")
	}

	b.state = adapter
	return b
}

// Builder method to set the field state in StateServiceBuilder
func (b *stateServiceBuilder) UseCache(settings CacheSettings) *stateServiceBuilder {
	if b.state == nil {
		panic("adapter is required")
	}

	b.state = newCache(b.state, settings)
	return b
}

// Builder method to set the field cb in StateServiceBuilder
func (b *stateServiceBuilder) CircuitBreaker(setting CircuitBreakerSettings) *stateServiceBuilder {
	b.cb = gobreaker.NewCircuitBreaker(setting.ToGobreakerSettings(b.cb.Name()))
	return b
}
