package core

import "context"

type Behavior interface {
	AddNext(next Behavior) Behavior
	Run(ctx context.Context, request Request) Result
	Next() Result
	setParameters(ctx context.Context, request Request, handler RequestHandler)
}

type Middleware struct {
	ctx     context.Context
	request Request
	handler RequestHandler
	next    Behavior
}

func (m *Middleware) AddNext(next Behavior) Behavior {
	m.next = next
	return m.next
}

func (m *Middleware) setParameters(ctx context.Context, request Request, handler RequestHandler) {
	m.ctx = ctx
	m.request = request
	m.handler = handler
}

func (m *Middleware) Next() Result {
	if err := m.ctx.Err(); err != nil {
		return Result{E: err}
	}

	if m.next != nil {
		m.next.setParameters(m.ctx, m.request, m.handler)
		return m.next.Run(m.ctx, m.request)
	} else {
		return m.handler(m.ctx, m.request)
	}
}
