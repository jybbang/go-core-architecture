package core

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/jybbang/go-core-architecture/core"
	"github.com/jybbang/go-core-architecture/infrastructure/mocks"
	"github.com/sony/gobreaker"
)

func Test_eventBus_AddDomainEvents(t *testing.T) {
	expect := 1000000
	mock := mocks.NewMockAdapter()

	m := core.NewMediatorBuilder().
		Create()
	e := core.NewEventBusBuilder().
		CircuitBreaker(core.CircuitBreakerSettings{
			Name:                 "t1",
			SamplingFailureCount: expect,
		}).
		MessaingAdapter(mock).
		CustomMediator(m).
		Create()

	for i := 0; i < expect; i++ {
		go e.AddDomainEvent(&core.DomainEvent{
			Topic: "t",
		})
	}

	ch := make(chan bool)
	go func() {
		time.Sleep(1 * time.Second)
		ch <- true
	}()

	then := 0
loop:
	for then != expect {
		select {
		case <-ch:
			break loop
		default:
			then = e.GetDomainEventsQueueCount()
			time.Sleep(10 * time.Millisecond)
		}
	}

	if then != expect {
		t.Errorf("Test_eventBus_AddDomainEvents() count = %v, expect %v", then, expect)
	}
}

func Test_eventBus_PublishDomainEvents(t *testing.T) {
	expect := 10000
	mock := mocks.NewMockAdapter()

	m := core.NewMediatorBuilder().
		AddNotificationHandler(new(okNotification), okNotificationHandler).
		Create()
	e := core.NewEventBusBuilder().
		CircuitBreaker(core.CircuitBreakerSettings{
			Name:                 "t2",
			SamplingFailureCount: expect,
		}).
		MessaingAdapter(mock).
		CustomMediator(m).
		Create()

	for i := 0; i < expect; i++ {
		event := new(okNotification)
		event.Topic = strconv.Itoa(i)
		e.AddDomainEvent(event)
	}

	when := e.GetDomainEventsQueueCount()
	if when != expect {
		t.Errorf("Test_eventBus_PublishDomainEvents() count = %v, expect %v", when, expect)
	}

	ctx := context.Background()
	err := e.PublishDomainEvents(ctx)
	if err != nil {
		t.Errorf("Test_eventBus_PublishDomainEvents() err = %v", err)
	}

	then := e.GetDomainEventsQueueCount()
	if then != 0 {
		t.Errorf("Test_eventBus_PublishDomainEvents() count = %v, expect %v", then, 0)
	}

	then2 := mock.GetPublishedCount()
	if then2 != uint32(expect) {
		t.Errorf("Test_eventBus_PublishDomainEvents() count = %v, expect %v", then2, expect)
	}
}

func Test_eventBus_PublishDomainEventsContextTimeoutShouldBeDeadlineExceeded(t *testing.T) {
	timeout := time.Duration(500 * time.Millisecond)
	expect := 1000
	mock := mocks.NewMockAdapter()

	m := core.NewMediatorBuilder().
		AddNotificationHandler(new(okNotification), okNotificationHandler).
		Create()
	e := core.NewEventBusBuilder().
		CircuitBreaker(core.CircuitBreakerSettings{
			Name:                 "t3",
			SamplingFailureCount: expect,
		}).
		MessaingAdapter(mock).
		CustomMediator(m).
		Create()

	for i := 0; i < expect; i++ {
		event := new(okNotification)
		event.Topic = strconv.Itoa(i)
		e.AddDomainEvent(event)
	}

	ctx, c := context.WithTimeout(context.TODO(), timeout)
	defer c()

	time.Sleep(timeout * 2)
	err := e.PublishDomainEvents(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Test_eventBus_PublishDomainEventsContextTimeoutShouldBeDeadlineExceeded() err = %v, expect %v", err, context.DeadlineExceeded)
	}
}

func Test_eventBus_PublishDomainEventsMediatorErrShouldBeError(t *testing.T) {
	expect := 1000
	mock := mocks.NewMockAdapter()

	m := core.NewMediatorBuilder().
		AddNotificationHandler(new(errNotification), errNotificationHandler).
		Create()
	e := core.NewEventBusBuilder().
		CircuitBreaker(core.CircuitBreakerSettings{
			Name:                 "t4",
			SamplingFailureCount: expect,
		}).
		MessaingAdapter(mock).
		CustomMediator(m).
		Create()

	for i := 0; i < expect; i++ {
		event := new(errNotification)
		event.Topic = strconv.Itoa(i)
		e.AddDomainEvent(event)
	}

	ctx := context.Background()
	err := e.PublishDomainEvents(ctx)
	if !errors.Is(err, core.ErrForbiddenAcccess) {
		t.Errorf("Test_eventBus_PublishDomainEventsMediatorErrShouldBeError() err = %v, expect %v", err, core.ErrForbiddenAcccess)
	}
}

func Test_eventBus_PublishDomainEventsCanNotPublishOptionShouldBeWorking(t *testing.T) {
	expect := 1000
	mock := mocks.NewMockAdapter()

	m := core.NewMediatorBuilder().
		AddNotificationHandler(new(okNotification), okNotificationHandler).
		Create()
	e := core.NewEventBusBuilder().
		CircuitBreaker(core.CircuitBreakerSettings{
			Name:                 "t5",
			SamplingFailureCount: expect,
		}).
		MessaingAdapter(mock).
		CustomMediator(m).
		Create()

	for i := 0; i < expect; i++ {
		event := new(okNotification)
		event.Topic = strconv.Itoa(i)
		event.CanNotPublishToEventsource = true
		e.AddDomainEvent(event)
	}

	ctx := context.Background()
	err := e.PublishDomainEvents(ctx)
	if err != nil {
		t.Errorf("Test_eventBus_PublishDomainEventsCanNotPublishOptionShouldBeWorking() err = %v", err)
	}

	then := e.GetDomainEventsQueueCount()
	if then != 0 {
		t.Errorf("Test_eventBus_PublishDomainEventsCanNotPublishOptionShouldBeWorking() count = %v, expect %v", then, 0)
	}

	then2 := mock.GetPublishedCount()
	if then2 != 0 {
		t.Errorf("Test_eventBus_PublishDomainEventsCanNotPublishOptionShouldBeWorking() count = %v, expect %v", then2, 0)
	}
}

func Test_eventBus_PublishDomainEventsCircuitBrakerShouldBeWorking(t *testing.T) {
	timeout := time.Duration(500 * time.Millisecond)
	expect := 10
	mock := mocks.NewMockAdapter()

	m := core.NewMediatorBuilder().
		AddNotificationHandler(new(okNotification), okNotificationHandler).
		AddNotificationHandler(new(errNotification), errNotificationHandler).
		Create()
	e := core.NewEventBusBuilder().
		CircuitBreaker(core.CircuitBreakerSettings{
			Name:                 "t6",
			DurationOfBreak:      timeout,
			SamplingFailureCount: expect,
		}).
		MessaingAdapter(mock).
		CustomMediator(m).
		Create()

	ctx := context.Background()
	for i := 0; i < expect; i++ {
		event := new(errNotification)
		event.Topic = strconv.Itoa(i)
		e.AddDomainEvent(event)
	}

	// expectRequest + 1
	event := new(okNotification)
	event.Topic = "ok"
	e.AddDomainEvent(event)

	err := e.PublishDomainEvents(ctx)
	if !errors.Is(err, gobreaker.ErrOpenState) {
		t.Errorf("Test_eventBus_PublishDomainEventsCircuitBrakerShouldBeWorking() err = %v, expect %v", err, gobreaker.ErrOpenState)
	}

	then := mock.GetPublishedCount()
	if then != 0 {
		t.Errorf("Test_eventBus_PublishDomainEventsCircuitBrakerShouldBeWorking() count = %v, expect %v", then, 0)
	}

	time.Sleep(timeout)

	event = new(okNotification)
	event.Topic = "ok2"
	e.AddDomainEvent(event)

	err = e.PublishDomainEvents(ctx)
	if err != nil {
		t.Errorf("Test_eventBus_PublishDomainEventsCircuitBrakerShouldBeWorking() err = %v", err)
	}

	then = mock.GetPublishedCount()
	if then != 1 {
		t.Errorf("Test_eventBus_PublishDomainEventsCircuitBrakerShouldBeWorking() count = %v, expect %v", then, 1)
	}
}

func Test_eventBus_PublishDomainEventsBufferedEventShouldBeWorking(t *testing.T) {
	timeout := time.Duration(500 * time.Millisecond)
	expect := 10
	mock := mocks.NewMockAdapter()

	m := core.NewMediatorBuilder().
		AddNotificationHandler(new(okNotification), okNotificationHandler).
		Create()
	e := core.NewEventBusBuilder().
		Settings(core.EventBusSettings{
			BufferedEventBufferTime: timeout,
		}).
		CircuitBreaker(core.CircuitBreakerSettings{
			Name:                 "t7",
			SamplingFailureCount: expect,
		}).
		MessaingAdapter(mock).
		CustomMediator(m).
		Create()

	ctx := context.Background()
	for i := 0; i < expect; i++ {
		event := new(okNotification)
		event.Topic = strconv.Itoa(i)
		event.CanBuffered = true
		e.AddDomainEvent(event)
	}

	e.PublishDomainEvents(ctx)
	then := mock.GetPublishedCount()
	if then != 0 {
		t.Errorf("Test_eventBus_PublishDomainEventsCircuitBrakerShouldBeWorking() count = %v, expect %v", then, 0)
	}

	time.Sleep(timeout * 2)

	then = mock.GetPublishedCount()
	if then != 1 {
		t.Errorf("Test_eventBus_PublishDomainEventsCircuitBrakerShouldBeWorking() count = %v, expect %v", then, 1)
	}
}
