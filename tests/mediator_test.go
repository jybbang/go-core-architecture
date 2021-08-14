package tests

import (
	"reflect"
	"testing"

	"github.com/jybbang/go-core-architecture/core"
	"github.com/jybbang/go-core-architecture/middlewares"
)

type testCommand struct {
	expect int
}

func testCommandHandler(request interface{}) interface{} {
	return request.(*testCommand).expect
}

type testNotification struct {
}

func testNotificationHandler(notification interface{}) {
}

func Test_mediator_Send(t *testing.T) {
	m := core.GetMediator().AddHandler(new(testCommand), testCommandHandler)
	m.AddMiddleware(middlewares.NewLogMiddleware())

	type args struct {
		request core.Request
	}
	tests := []struct {
		name    string
		m       *core.Mediator
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "1",
			m:    m,
			args: args{
				request: &testCommand{
					expect: 123,
				},
			},
			want:    123,
			wantErr: false,
		},
		{
			name: "2",
			m:    m,
			args: args{
				request: &testCommand{
					expect: 1234,
				},
			},
			want:    1234,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Send(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("mediator.Send() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mediator.Send() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mediator_Publish(t *testing.T) {
	m := core.GetMediator().AddNotificationHandler(new(testNotification), testNotificationHandler)

	type args struct {
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
				notification: &testNotification{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Publish(tt.args.notification); (err != nil) != tt.wantErr {
				t.Errorf("mediator.Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
