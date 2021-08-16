package core

import (
	"context"
	"testing"
	"time"

	"github.com/jybbang/go-core-architecture/core"
)

func Test_mediator_Send(t *testing.T) {
	expect := 1000000
	m := core.GetMediator()
	ctx := context.Background()

	sumExpect := 0
	sum := 0
	for i := 0; i < expect; i++ {
		result := m.Send(ctx, &okCommand{
			Expect: i,
		})
		sumExpect += i
		sum += result.V.(int)
	}

	if sum != sumExpect {
		t.Errorf("Test_mediator_Send() sum = %v, expect %v", sum, sumExpect)
	}

	for i := 0; i < expect; i++ {
		m.Send(ctx, &errCommand{})
	}

	then := int(m.GetSentCount())
	if then != expect {
		t.Errorf("Test_mediator_Send() count = %v, expect %v", then, expect)
	}
}

func Test_mediator_Publish(t *testing.T) {
	expect := 1000000
	m := core.GetMediator()
	ctx := context.Background()

	for i := 0; i < expect; i++ {
		go m.Publish(ctx, &okNotification{})
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
			then = int(m.GetPublishedCount())
			time.Sleep(10 * time.Millisecond)
		}
	}

	if then != expect {
		t.Errorf("Test_mediator_Publish() count = %v, expect %v", then, expect)
	}
}
