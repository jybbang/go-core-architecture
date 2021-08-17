package core

import "context"

type behavior interface {
	AddNext(next behavior) behavior
	Run(ctx context.Context, request Request) Result
	Next() Result
	setParameters(ctx context.Context, request Request, handler RequestHandler)
}

type Middleware struct {
	ctx     context.Context
	request Request
	handler RequestHandler
	next    behavior
}

func (m *Middleware) AddNext(next behavior) behavior {
	m.next = next
	return m.next
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

func (m *Middleware) setParameters(ctx context.Context, request Request, handler RequestHandler) {
	m.ctx = ctx
	m.request = request
	m.handler = handler
}
