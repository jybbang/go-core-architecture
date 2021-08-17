package core

import (
	"context"

	"github.com/jybbang/go-core-architecture/core"
	"github.com/jybbang/go-core-architecture/infrastructure/mocks"
)

type testModel struct {
	core.Entity
	Expect int `bson:"expect,omitempty"`
}

type okCommand struct {
	Expect int
}

type errCommand struct {
	Expect int
}

func okCommandHandler(ctx context.Context, request interface{}) core.Result {
	return core.Result{V: request.(*okCommand).Expect}
}

func errCommandHandler(ctx context.Context, request interface{}) core.Result {
	return core.Result{E: core.ErrForbiddenAcccess}
}

type okNotification struct {
	core.DomainEvent
}

type errNotification struct {
	core.DomainEvent
}

func okNotificationHandler(ctx context.Context, notification interface{}) error {
	return nil
}

func errNotificationHandler(ctx context.Context, notification interface{}) error {
	return core.ErrForbiddenAcccess
}

func init() {
	core.NewMediatorBuilder().
		AddHandler(new(okCommand), okCommandHandler).
		AddHandler(new(errCommand), errCommandHandler).
		AddNotificationHandler(new(okNotification), okNotificationHandler).
		AddNotificationHandler(new(errNotification), errNotificationHandler).
		Build()

	core.NewEventbusBuilder().MessaingAdapter(mocks.NewMockAdapter()).Build()
}
