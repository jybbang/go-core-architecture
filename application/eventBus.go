package application

import (
	"fmt"
	"sync"
	"time"

	"github.com/jybbang/core-architecture/application/contracts"
	"github.com/jybbang/core-architecture/domain"
	"github.com/sony/gobreaker"
)

type eventBus struct {
	mediator     *mediator
	messaging    contracts.MessagingAdapter
	domainEvents []*domain.DomainEvent
	cb           *gobreaker.CircuitBreaker
	sync.RWMutex
}

func (e *eventBus) SetMessaingAdapter(messageService contracts.MessagingAdapter) *eventBus {
	e.messaging = messageService
	return e
}

func (e *eventBus) AddDomainEvent(domainEvent *domain.DomainEvent) {
	e.Lock()
	defer e.Unlock()

	domainEvent.AddingEvent()

	e.domainEvents = append(e.domainEvents, domainEvent)
}

func (e *eventBus) PublishDomainEvents() error {
	e.Lock()
	defer e.Unlock()

	now := time.Now()
	_, err := e.cb.Execute(func() (interface{}, error) {

		for !e.empty() {
			event, eventErr := e.dequeueDomainEvent()
			if eventErr != nil {
				return nil, eventErr
			}

			event.PublishingEvent(now)

			mediatorErr := e.mediator.Publish(event)
			if mediatorErr != nil {
				return nil, mediatorErr
			}

			if !event.CanPublishToEventsource {
				continue
			}

			msgErr := e.messaging.Publish(event)
			if msgErr != nil {
				return nil, msgErr
			}
		}

		return nil, nil
	})
	return err
}

func (e *eventBus) dequeueDomainEvent() (*domain.DomainEvent, error) {
	if len(e.domainEvents) > 0 {
		result := e.domainEvents[0]
		e.domainEvents = e.domainEvents[1:]
		return result, nil
	}

	return nil, fmt.Errorf("domainEvents empty exception")
}

func (e *eventBus) empty() bool {
	return len(e.domainEvents) == 0
}
