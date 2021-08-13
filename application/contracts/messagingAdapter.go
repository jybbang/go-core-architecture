package contracts

import "github.com/jybbang/core-architecture/domain"

type ReplyHandler func(string)

type MessagingAdapter interface {
	Publish(*domain.DomainEvent) error
	Subscribe(string, ReplyHandler) error
	Unsubscribe(string) error
}
