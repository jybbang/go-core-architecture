package core

import (
	"context"
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
	sync.Mutex
}

type BufferedEvent struct {
	DomainEvent
	BufferedEvents []DomainEventer
}

func (e *EventBus) Initialize() *EventBus {
	observable := rxgo.FromChannel(e.ch).
		BufferWithTimeOrCount(rxgo.WithDuration(1*time.Second), 1000)

	go e.SubscribeBufferedEvent(observable)

	return e
}

func (e *EventBus) SetupCb(setting gobreaker.Settings) *EventBus {
	setting.Name = e.cb.Name()
	setting.OnStateChange = OnCbStateChange
	e.cb = gobreaker.NewCircuitBreaker(setting)
	return e
}

func (e *EventBus) SubscribeBufferedEvent(observable rxgo.Observable) {
	ch := observable.Observe()

	for {
		vals := <-ch
		event := &BufferedEvent{
			BufferedEvents: vals.V.([]DomainEventer),
		}
		timeout, _ := context.WithTimeout(context.Background(), time.Duration(2*time.Second))
		e.AddDomainEvent(event)
		e.PublishDomainEvents(timeout)
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

func (e *EventBus) PublishDomainEvents(ctx context.Context) error {
	e.Lock()
	defer e.Unlock()

	now := time.Now()
	_, err := e.cb.Execute(func() (interface{}, error) {

		for !e.empty() {
			event := e.dequeueDomainEvent()
			event.PublishingEvent(ctx, now)

			mediatorErr := e.mediator.Publish(ctx, event)
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

			msgErr := e.messaging.Publish(ctx, event)
			if msgErr != nil {
				return nil, msgErr
			}
		}

		return nil, nil
	})

	return err
}

func (e *EventBus) dequeueDomainEvent() DomainEventer {
	result := e.domainEvents[0]
	e.domainEvents = e.domainEvents[1:]
	return result
}

func (e *EventBus) empty() bool {
	return len(e.domainEvents) == 0
}
