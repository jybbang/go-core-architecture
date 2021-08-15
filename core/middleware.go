package core

import "context"

type Middleware struct {
	next middlewarer
}

type middlewarer interface {
	AddMiddleware(middlewarer) middlewarer
	Run(ctx context.Context, request Request) (ok bool, err error)
	nextRun(ctx context.Context, request Request, handler RequestHandler) Result
}

func (m *Middleware) AddMiddleware(middleware middlewarer) middlewarer {
	m.next = middleware
	return m.next
}

func (m *Middleware) nextRun(ctx context.Context, request Request, handler RequestHandler) Result {
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
		return m.next.nextRun(ctx, request, handler)
	} else {
		return handler(ctx, request)
	}
}
