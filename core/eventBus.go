package core

import (
	"context"
	"fmt"
	"time"

	"github.com/enriquebris/goconcurrentqueue"
	"github.com/reactivex/rxgo/v2"
	"github.com/sony/gobreaker"
)

type eventBus struct {
	mediator     *mediator
	messaging    messagingAdapter
	domainEvents *goconcurrentqueue.FIFO
	ch           chan rxgo.Item
	cb           *gobreaker.CircuitBreaker
	settings     EventBusSettings
}

type bufferedEvent struct {
	DomainEvent
	BufferedEvents []DomainEventer
}

// dummy
func bufferedEventHandler(ctx context.Context, notification interface{}) error {
	return nil
}

func (e *eventBus) initialize() *eventBus {
	if err := e.connect(); err != nil {
		panic(err)
	}

	observable := rxgo.FromChannel(e.ch).
		BufferWithTimeOrCount(rxgo.WithDuration(e.settings.BufferedEventBufferTime), e.settings.BufferedEventBufferCount)

	go e.subscribeBufferedEvent(observable)

	return e
}

func (e *eventBus) connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), e.settings.ConnectionTimeout)
	defer cancel()

	return e.messaging.Connect(ctx)
}

func (e *eventBus) onCircuitOpen() {
	e.messaging.Disconnect()

	if !e.messaging.IsConnected() {
		e.connect()
	}
}

func (e *eventBus) subscribeBufferedEvent(observable rxgo.Observable) {
	ch := observable.Observe()

	for {
		items := <-ch

		events, ok := items.V.([]interface{})
		if !ok || len(events) == 0 {
			continue
		}

		bufferedEvent := &bufferedEvent{}
		for _, v := range events {
			if event, ok := v.(DomainEventer); ok {
				bufferedEvent.BufferedEvents = append(bufferedEvent.BufferedEvents, event)
			}
		}

		bufferedEvent.Topic = "BufferedEvents"
		timeout, cancel := context.WithTimeout(context.Background(), e.settings.BufferedEventTimeout)
		e.AddDomainEvent(bufferedEvent)
		e.PublishDomainEvents(timeout)
		cancel()
	}
}

func (e *eventBus) empty() bool {
	return e.GetDomainEventsQueueCount() == 0
}

func (e *eventBus) GetDomainEventsQueueCount() int {
	return e.domainEvents.GetLen()
}

func (e *eventBus) AddDomainEvent(domainEvent DomainEventer) {
	if domainEvent.GetTopic() == "" {
		panic("topic is required")
	}

	domainEvent.SetAddingEvent()
	e.domainEvents.Enqueue(domainEvent)
}

func (e *eventBus) PublishDomainEvents(ctx context.Context) error {
	defer publishEventsPanicRecover()

	var err error
	now := time.Now()

	for !e.empty() {
		item, _ := e.domainEvents.Dequeue()
		event := item.(DomainEventer)

		_, err = e.cb.Execute(func() (interface{}, error) {
			err = e.mediator.Publish(ctx, event)
			if err != nil {
				return nil, err
			}

			if event.GetCanNotPublishToEventsource() {
				return nil, nil
			}

			event.SetPublishingEvent(ctx, now)

			if event.GetCanBuffered() {
				e.ch <- rxgo.Item{
					V: event,
				}
				return nil, nil
			}

			err = e.Publish(ctx, event)
			if err != nil {
				return nil, err
			}

			return nil, nil
		})
	}

	return err
}

func (e *eventBus) Publish(ctx context.Context, event DomainEventer) error {
	return e.messaging.Publish(ctx, event)
}

func (e *eventBus) Subscribe(ctx context.Context, topic string, handler ReplyHandler) error {
	return e.messaging.Subscribe(ctx, topic, handler)
}

func (e *eventBus) Unsubscribe(ctx context.Context, topic string) error {
	return e.messaging.Unsubscribe(ctx, topic)
}

func publishEventsPanicRecover() {
	if r := recover(); r != nil {
		fmt.Println("publish domain events recovering from panic:", r)
	}
}
