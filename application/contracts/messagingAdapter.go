package contracts

import "github.com/jybbang/go-core-architecture/domain"

type ReplyHandler func(string)

type MessagingAdapter interface {
	Publish(domain.DomainEventer) error
	Subscribe(string, ReplyHandler) error
	Unsubscribe(string) error
}
