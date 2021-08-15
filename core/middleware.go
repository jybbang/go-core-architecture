package core

import "context"

type Middleware struct {
	next Middlewarer
}

type Middlewarer interface {
	AddMiddleware(Middlewarer) Middlewarer
	Next(ctx context.Context, request Request, handler RequestHandler) Result
	Run(ctx context.Context, request Request) (ok bool, err error)
}

func (m *Middleware) AddMiddleware(middleware Middlewarer) Middlewarer {
	m.next = middleware
	return m.next
}

func (m *Middleware) Next(ctx context.Context, request Request, handler RequestHandler) Result {
	// Check context cancellation
	if err := ctx.Err(); err != nil {
		return Result{E: err}
	}

	if m.next != nil {
		ok, err := m.next.Run(ctx, request)
		if err != nil {
			return Result{
				V: nil,
				E: err,
			}
		}
		if !ok {
			return Result{
				V: nil,
				E: ErrForbiddenAcccess,
			}
		}
		return m.next.Next(ctx, request, handler)
	} else {
		return handler(ctx, request)
	}
}
