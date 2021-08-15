package core

import "context"

type messagingAdapter interface {
	Publish(ctx context.Context, event DomainEventer) error
	Subscribe(ctx context.Context, topic string, handler ReplyHandler) error
	Unsubscribe(ctx context.Context, topic string) error
}
