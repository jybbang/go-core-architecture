package core

type ReplyHandler func(string)

type MessagingAdapter interface {
	Publish(DomainEventer) error
	Subscribe(string, ReplyHandler) error
	Unsubscribe(string) error
}
