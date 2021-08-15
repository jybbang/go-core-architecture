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
		Name:        "state service",
		Timeout:     cbDefaultTimeout,
		MaxRequests: cbDefaultAllowedRequests,
	}
	o.cb = gobreaker.NewCircuitBreaker(st)

	return o
}

// Build Method which creates StateService
func (b *stateServiceBuilder) Build() *stateService {
	if statesInstance != nil {
		panic("state service already created")
	}

	statesInstance = &stateService{
		state: b.state,
		cb:    b.cb,
	}
	statesInstance.initialize()

	return statesInstance
}

// Builder method to set the field state in StateServiceBuilder
func (b *stateServiceBuilder) StateAdapter(adapter stateAdapter) *stateServiceBuilder {
	b.state = adapter
	return b
}

// Builder method to set the field cb in StateServiceBuilder
func (b *stateServiceBuilder) CircuitBreaker(setting gobreaker.Settings) *stateServiceBuilder {
	setting.Name = b.cb.Name()
	b.cb = gobreaker.NewCircuitBreaker(setting)
	return b
}
