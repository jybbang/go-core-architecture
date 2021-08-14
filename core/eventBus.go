package core

import (
	"fmt"
	"sync"
	"time"

	"github.com/reactivex/rxgo/v2"
	"github.com/sony/gobreaker"
)

type EventBus struct {
	mediator     *Mediator
	messaging    MessagingAdapter
	domainEvents []DomainEventer
	ch           chan rxgo.Item
	cb           *gobreaker.CircuitBreaker
	sync.RWMutex
}

type BufferedEvent struct {
	DomainEvent
	Events []DomainEventer
}

func (e *EventBus) Setup() *EventBus {
	observable := rxgo.FromChannel(e.ch).
		BufferWithTimeOrCount(rxgo.WithDuration(1*time.Second), 1000)

	go e.SubscribeBufferedEvent(observable)

	return e
}

func (e *EventBus) SubscribeBufferedEvent(observable rxgo.Observable) {
	ch := observable.Observe()

	for {
		vals := <-ch
		event := &BufferedEvent{
			Events: vals.V.([]DomainEventer),
		}
		e.AddDomainEvent(event)
		e.PublishDomainEvents()
	}
}

func (e *EventBus) SetMessaingAdapter(messageService MessagingAdapter) *EventBus {
	e.messaging = messageService
	return e
}

func (e *EventBus) AddDomainEvent(domainEvent DomainEventer) {
	e.Lock()
	defer e.Unlock()

	domainEvent.AddingEvent()

	e.domainEvents = append(e.domainEvents, domainEvent)
}

func (e *EventBus) PublishDomainEvents() error {
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

			if event.GetCanNotPublishToEventsource() {
				continue
			}

			if event.GetCanBuffered() {
				e.ch <- rxgo.Item{
					V: event,
				}
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

func (e *EventBus) dequeueDomainEvent() (DomainEventer, error) {
	if len(e.domainEvents) > 0 {
		result := e.domainEvents[0]
		e.domainEvents = e.domainEvents[1:]
		return result, nil
	}

	return nil, fmt.Errorf("domainEvents empty exception")
}

func (e *EventBus) empty() bool {
	return len(e.domainEvents) == 0
}
