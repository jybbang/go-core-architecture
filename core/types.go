package core

import "context"

type Services struct {
	eventbus *eventbus
	states   *stateService
}

type Request interface{}
type RequestHandler func(
	ctx context.Context,
	services Services,
	request interface{}) Result

type Notification interface{}
type NotificationHandler func(
	ctx context.Context,
	services Services,
	notification interface{}) error

type ReplyHandler func(receivedData interface{})
