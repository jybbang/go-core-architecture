package core

import (
	"context"
	"reflect"
	"testing"

	"github.com/jybbang/go-core-architecture/core"
	"github.com/jybbang/go-core-architecture/middlewares"
)

type testCommand struct {
	expect int
}

func testCommandHandler(ctx context.Context, request interface{}) core.Result {
	return core.Result{V: request.(*testCommand).expect}
}

func Test_mediator_Send(t *testing.T) {
	m := core.GetMediator().AddHandler(new(testCommand), testCommandHandler)
	m.AddMiddleware(middlewares.NewLogMiddleware())

	type args struct {
		ctx     context.Context
		request core.Request
	}
	tests := []struct {
		name string
		m    *core.Mediator
		args args
		want core.Result
	}{
		{
			name: "1",
			m:    m,
			args: args{
				ctx: context.Background(),
				request: &testCommand{
					expect: 123,
				},
			},
			want: core.Result{V: 123},
		},
		{
			name: "2",
			m:    m,
			args: args{
				ctx: context.Background(),
				request: &testCommand{
					expect: 1234,
				},
			},
			want: core.Result{V: 1234},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Send(tt.args.ctx, tt.args.request); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Mediator.Send() = %v, want %v", got, tt.want)
			}
		})
	}
}

type testNotification struct {
}

func testNotificationHandler(ctx context.Context, notification interface{}) error {
	return nil
}

func Test_mediator_Publish(t *testing.T) {
	m := core.GetMediator().AddNotificationHandler(new(testNotification), testNotificationHandler)

	type args struct {
		ctx          context.Context
		notification core.Notification
	}
	tests := []struct {
		name    string
		m       *core.Mediator
		args    args
		wantErr bool
	}{
		{
			name: "1",
			m:    m,
			args: args{
				ctx:          context.Background(),
				notification: &testNotification{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Publish(tt.args.ctx, tt.args.notification); (err != nil) != tt.wantErr {
				t.Errorf("Mediator.Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}