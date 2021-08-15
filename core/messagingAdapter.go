package core

import "context"

type ReplyHandler func(receivedData interface{})

type MessagingAdapter interface {
	Publish(ctx context.Context, event DomainEventer) error
	Subscribe(ctx context.Context, topic string, handler ReplyHandler) error
	Unsubscribe(ctx context.Context, topic string) error
}
