package core

import (
	"fmt"
	"time"

	"gopkg.in/jeevatkm/go-model.v1"
)

// Builder Object for StateService
type stateServiceBuilder struct {
	state      stateAdapter
	cbSettings CircuitBreakerSettings
	settings   StateServiceSettings
}

// Constructor for StateServiceBuilder
func NewStateServiceBuilder() *stateServiceBuilder {
	o := new(stateServiceBuilder)
	o.cbSettings = CircuitBreakerSettings{
		AllowedRequestInHalfOpen: 1,
		DurationOfBreak:          time.Duration(60 * time.Second),
		SamplingDuration:         time.Duration(60 * time.Second),
		SamplingFailureCount:     5,
	}
	o.settings = StateServiceSettings{
		ConnectionTimeout: time.Duration(10 * time.Second),
	}

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
		state:    b.state,
		settings: b.settings,
	}

	instance.cb = b.cbSettings.ToCircuitBreaker("state service", instance.onCircuitOpen)

	instance.initialize()

	return instance
}

// Builder method to set the field messaging in EventBusBuilder
func (b *stateServiceBuilder) Settings(settings StateServiceSettings) *stateServiceBuilder {
	err := model.Copy(&b.settings, settings)

	if err != nil {
		panic(fmt.Errorf("settings mapping errors occurred: %v", err))
	}

	return b
}

// Builder method to set the field cb in StateServiceBuilder
func (b *stateServiceBuilder) CircuitBreaker(settings CircuitBreakerSettings) *stateServiceBuilder {
	err := model.Copy(&b.cbSettings, settings)

	if err != nil {
		panic(fmt.Errorf("cb settings mapping errors occurred: %v", err))
	}

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

// Builder method to set the field state in StateServiceBuilder
func (b *stateServiceBuilder) StateAdapter(adapter stateAdapter) *stateServiceBuilder {
	if adapter == nil {
		panic("adapter is required")
	}

	b.state = adapter
	return b
}
