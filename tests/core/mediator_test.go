package core

import (
	"context"
	"reflect"
	"testing"

	"github.com/jybbang/go-core-architecture/core"
	"github.com/jybbang/go-core-architecture/middlewares"
	"go.uber.org/zap"
)

type testCommand struct {
	Expect int `validate:"eq=100"`
}

func testCommandHandler(ctx context.Context, request interface{}) core.Result {
	return core.Result{V: request.(*testCommand).Expect}
}

func Test_mediator_Send(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	m := core.NewMediatorBuilder().
		AddHandler(new(testCommand), testCommandHandler).
		Build()

	m.AddMiddleware(middlewares.NewLogMiddleware(logger)).
		AddMiddleware(middlewares.NewValidationMiddleware())

	core.NewEventbusBuilder().Build()
	core.NewStateServiceBuilder().Build()

	type args struct {
		ctx     context.Context
		request core.Request
	}
	tests := []struct {
		name string
		args args
		want core.Result
	}{
		{
			"1",
			args{
				ctx: context.Background(),
				request: &testCommand{
					Expect: 100,
				},
			},
			core.Result{V: 100},
		},
		{
			"2",
			args{
				ctx: context.Background(),
				request: &testCommand{
					Expect: 99,
				},
			},
			core.Result{E: core.ErrBadRequest},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := m.Send(tt.args.ctx, tt.args.request); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Mediator.Send() = %v, want %v", got, tt.want)
			}
		})
	}
}

type testNotification struct {
	core.DomainEvent
}

func testNotificationHandler(ctx context.Context, notification interface{}) error {
	return nil
}

func Test_mediator_Publish(t *testing.T) {
	m := core.NewMediatorBuilder().
		AddNotificationHandler(new(testNotification), testNotificationHandler).
		Build()

	core.NewEventbusBuilder().Build()
	core.NewStateServiceBuilder().Build()

	type args struct {
		ctx          context.Context
		notification core.Notification
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"1",
			args{
				ctx:          context.Background(),
				notification: &testNotification{},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := m.Publish(tt.args.ctx, tt.args.notification); (err != nil) != tt.wantErr {
				t.Errorf("Mediator.Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
