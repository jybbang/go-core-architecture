package core

type ReplyHandler func(interface{})

type MessagingAdapter interface {
	Publish(DomainEventer) error
	Subscribe(string, ReplyHandler) error
	Unsubscribe(string) error
}
