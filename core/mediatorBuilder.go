package core

import cmap "github.com/orcaman/concurrent-map"

// Builder Object for Mediator
type mediatorBuilder struct {
}

// Constructor for MediatorBuilder
func NewMediatorBuilder() *mediatorBuilder {
	o := new(mediatorBuilder)
	return o
}

// Build Method which creates Mediator
func (b *mediatorBuilder) Build() *mediator {
	if mediatorInstance != nil {
		panic("mediator already created")
	}

	mediatorInstance = &mediator{
		requestHandlers:      cmap.New(),
		notificationHandlers: cmap.New(),
	}
	mediatorInstance.initialize()

	return mediatorInstance
}
