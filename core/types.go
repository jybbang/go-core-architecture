package core

import "context"

type Request interface{}
type RequestHandler func(ctx context.Context, request interface{}) Result

type Notification interface{}
type NotificationHandler func(ctx context.Context, notification interface{}) error

type ReplyHandler func(receivedData interface{})
