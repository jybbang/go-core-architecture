package core

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/jybbang/go-core-architecture/core"
	"github.com/jybbang/go-core-architecture/infrastructure/mocks"
	"github.com/sony/gobreaker"
)

func TestEventBus_AddDomainEvents(t *testing.T) {
	expect := 1000000
	mock := mocks.NewMockAdapter()
	e := core.NewEventbusBuilder().
		CircuitBreaker(core.CircuitBreakerSettings{
			Name:                 "t1",
			SamplingFailureCount: expect,
		}).MessaingAdapter(mock).Create()

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
		t.Errorf("TestEventBus_AddDomainEvents() count = %v, expect %v", then, expect)
	}
}

func TestEventBus_PublishDomainEvents(t *testing.T) {
	expect := 10000
	mock := mocks.NewMockAdapter()
	e := core.NewEventbusBuilder().
		CircuitBreaker(core.CircuitBreakerSettings{
			Name:                 "t2",
			SamplingFailureCount: expect,
		}).MessaingAdapter(mock).Create()

	for i := 0; i < expect; i++ {
		event := new(okNotification)
		event.Topic = strconv.Itoa(i)
		e.AddDomainEvent(event)
	}

	when := e.GetDomainEventsQueueCount()
	if when != expect {
		t.Errorf("TestEventBus_PublishDomainEvents() count = %v, expect %v", when, expect)
	}

	ctx := context.Background()
	err := e.PublishDomainEvents(ctx)
	if err != nil {
		t.Errorf("TestEventBus_PublishDomainEvents() err = %v", err)
	}

	then := e.GetDomainEventsQueueCount()
	if then != 0 {
		t.Errorf("TestEventBus_PublishDomainEvents() count = %v, expect %v", then, 0)
	}

	then2 := mock.GetPublishedCount()
	if then2 != uint32(expect) {
		t.Errorf("TestEventBus_PublishDomainEvents() count = %v, expect %v", then2, expect)
	}
}

func TestEventBus_PublishDomainEventsContextTimeoutShouldBeError(t *testing.T) {
	expect := 1000
	mock := mocks.NewMockAdapter()
	e := core.NewEventbusBuilder().
		CircuitBreaker(core.CircuitBreakerSettings{
			Name:                 "t3",
			SamplingFailureCount: expect,
		}).MessaingAdapter(mock).Create()

	for i := 0; i < expect; i++ {
		event := new(okNotification)
		event.Topic = strconv.Itoa(i)
		e.AddDomainEvent(event)
	}

	ctx, c := context.WithTimeout(context.TODO(), time.Duration(1*time.Second))
	defer c()

	time.Sleep(1 * time.Second)
	err := e.PublishDomainEvents(ctx)
	if err != context.DeadlineExceeded {
		t.Errorf("TestEventBus_PublishDomainEventsContextTimeoutShouldBeError() err = %v, expect %v", err, context.DeadlineExceeded)
	}
}

func TestEventBus_PublishDomainEventsMediatorErrShouldBeError(t *testing.T) {
	expect := 1000
	mock := mocks.NewMockAdapter()
	e := core.NewEventbusBuilder().
		CircuitBreaker(core.CircuitBreakerSettings{
			Name:                 "t4",
			SamplingFailureCount: expect,
		}).MessaingAdapter(mock).Create()

	for i := 0; i < expect; i++ {
		event := new(errNotification)
		event.Topic = strconv.Itoa(i)
		e.AddDomainEvent(event)
	}

	ctx := context.Background()
	err := e.PublishDomainEvents(ctx)
	if err != core.ErrForbiddenAcccess {
		t.Errorf("TestEventBus_PublishDomainEventsMediatorErrShouldBeError() err = %v, expect %v", err, core.ErrForbiddenAcccess)
	}
}

func TestEventBus_PublishDomainEventsCanNotPublishOptionShouldBeWorking(t *testing.T) {
	expect := 1000
	mock := mocks.NewMockAdapter()
	e := core.NewEventbusBuilder().
		CircuitBreaker(core.CircuitBreakerSettings{
			Name:                 "t5",
			SamplingFailureCount: expect,
		}).MessaingAdapter(mock).Create()

	for i := 0; i < expect; i++ {
		event := new(okNotification)
		event.Topic = strconv.Itoa(i)
		event.CanNotPublishToEventsource = true
		e.AddDomainEvent(event)
	}

	ctx := context.Background()
	err := e.PublishDomainEvents(ctx)
	if err != nil {
		t.Errorf("TestEventBus_PublishDomainEventsCanNotPublishOptionShouldBeWorking() err = %v", err)
	}

	then := e.GetDomainEventsQueueCount()
	if then != 0 {
		t.Errorf("TestEventBus_PublishDomainEventsCanNotPublishOptionShouldBeWorking() count = %v, expect %v", then, 0)
	}

	then2 := mock.GetPublishedCount()
	if then2 != 0 {
		t.Errorf("TestEventBus_PublishDomainEventsCanNotPublishOptionShouldBeWorking() count = %v, expect %v", then2, 0)
	}
}

func TestEventBus_PublishDomainEventsCircuitBrakerShouldBeWorking(t *testing.T) {
	var timeout = time.Duration(500 * time.Millisecond)
	var expect = 10

	mock := mocks.NewMockAdapter()
	e := core.NewEventbusBuilder().
		CircuitBreaker(core.CircuitBreakerSettings{
			Name:                 "t6",
			DurationOfBreak:      timeout,
			SamplingFailureCount: expect,
		}).MessaingAdapter(mock).Create()

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
	if err != gobreaker.ErrOpenState {
		t.Errorf("TestEventBus_PublishDomainEventsMediatorErrShouldBeError() err = %v, expect %v", err, gobreaker.ErrOpenState)
	}

	then := mock.GetPublishedCount()
	if then != 0 {
		t.Errorf("TestEventBus_PublishDomainEventsCircuitBrakerShouldBeWorking() count = %v, expect %v", then, 0)
	}

	time.Sleep(timeout)

	event = new(okNotification)
	event.Topic = "ok2"
	e.AddDomainEvent(event)

	err = e.PublishDomainEvents(ctx)
	if err != nil {
		t.Errorf("TestEventBus_PublishDomainEventsCircuitBrakerShouldBeWorking() err = %v", err)
	}

	then = mock.GetPublishedCount()
	if then != 1 {
		t.Errorf("TestEventBus_PublishDomainEventsCircuitBrakerShouldBeWorking() count = %v, expect %v", then, 1)
	}
}

func TestEventBus_PublishDomainEventsBufferedEventShouldBeWorking(t *testing.T) {
	var timeout = time.Duration(500 * time.Millisecond)
	var expect = 10

	mock := mocks.NewMockAdapter()
	e := core.NewEventbusBuilder().
		Setting(core.EventbusSettings{
			BufferedEventBufferTime: timeout,
		}).
		CircuitBreaker(core.CircuitBreakerSettings{
			Name:                 "t7",
			SamplingFailureCount: expect,
		}).MessaingAdapter(mock).Create()

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
		t.Errorf("TestEventBus_PublishDomainEventsCircuitBrakerShouldBeWorking() count = %v, expect %v", then, 0)
	}

	time.Sleep(timeout * 2)

	then = mock.GetPublishedCount()
	if then != 1 {
		t.Errorf("TestEventBus_PublishDomainEventsCircuitBrakerShouldBeWorking() count = %v, expect %v", then, 1)
	}
}
