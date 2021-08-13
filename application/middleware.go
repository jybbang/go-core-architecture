package application

import (
	"fmt"
)

type Middleware struct {
	next Middlewarer
}

type Middlewarer interface {
	AddMiddleware(Middlewarer) Middlewarer
	Next(Request, RequestHandler) (interface{}, error)
	Run(Request) (bool, error)
}

func (m *Middleware) AddMiddleware(middleware Middlewarer) Middlewarer {
	m.next = middleware
	return m.next
}

func (m *Middleware) Next(request Request, handler RequestHandler) (interface{}, error) {
	ok, err := m.next.Run(request)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("middleware block this request")
	}

	if m.next != nil {
		return m.Next(request, handler)
	} else {
		return handler(request), nil
	}
}
