package core

import (
	"context"
	"sync"
	"time"

	"github.com/reactivex/rxgo/v2"
	"github.com/sony/gobreaker"
)

type eventbus struct {
	mediator     *mediator
	messaging    messagingAdapter
	domainEvents []DomainEventer
	ch           chan rxgo.Item
	cb           *gobreaker.CircuitBreaker
	mutex        sync.Mutex
	settings     EventbusSettings
}

type EventbusSettings struct {
	BufferedEventBufferTime  time.Duration
	BufferedEventBufferCount int
	BufferedEventTimeout     time.Duration
}

type bufferedEvent struct {
	DomainEvent
	BufferedEvents []DomainEventer
}

func bufferedEventHandler(ctx context.Context, notification interface{}) error {
	return nil
}

func (e *eventbus) initialize() *eventbus {
	observable := rxgo.FromChannel(e.ch).
		BufferWithTimeOrCount(rxgo.WithDuration(e.settings.BufferedEventBufferTime), e.settings.BufferedEventBufferCount)

	go e.subscribeBufferedEvent(observable)

	return e
}

func (e *eventbus) subscribeBufferedEvent(observable rxgo.Observable) {
	ch := observable.Observe()

	for {
		items := <-ch
		if items.V == nil {
			continue
		}

		events := &bufferedEvent{}
		for _, v := range items.V.([]interface{}) {
			if event, ok := v.(DomainEventer); ok {
				events.BufferedEvents = append(events.BufferedEvents, event)
			}
		}

		events.Topic = "BufferedEvents"
		timeout, cancel := context.WithTimeout(context.Background(), e.settings.BufferedEventTimeout)
		e.AddDomainEvent(events)
		e.PublishDomainEvents(timeout)
		cancel()
	}
}

func (e *eventbus) dequeueDomainEvent() DomainEventer {
	result := e.domainEvents[0]
	e.domainEvents = e.domainEvents[1:]
	return result
}

func (e *eventbus) empty() bool {
	return e.GetDomainEventsQueueCount() == 0
}

func (e *eventbus) GetDomainEventsQueueCount() int {
	return len(e.domainEvents)
}

func (e *eventbus) AddDomainEvent(domainEvent DomainEventer) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	domainEvent.AddingEvent()

	if domainEvent.GetTopic() == "" {
		panic("topic is required")
	}

	e.domainEvents = append(e.domainEvents, domainEvent)
}

func (e *eventbus) PublishDomainEvents(ctx context.Context) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	var err error
	now := time.Now()

	for !e.empty() {
		event := e.dequeueDomainEvent()

		_, err = e.cb.Execute(func() (interface{}, error) {
			if err := ctx.Err(); err != nil {
				return nil, err
			}

			event.PublishingEvent(ctx, now)

			mediatorErr := e.mediator.Publish(ctx, event)
			if mediatorErr != nil {
				return nil, mediatorErr
			}

			if event.GetCanNotPublishToEventsource() {
				return nil, nil
			}

			if event.GetCanBuffered() {
				e.ch <- rxgo.Item{
					V: event,
				}
				return nil, nil
			}

			msgErr := e.messaging.Publish(ctx, event)
			if msgErr != nil {
				return nil, msgErr
			}

			return nil, nil
		})
	}

	return err
}
