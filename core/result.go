package core

import (
	"context"
	"errors"
	"net/http"

	"github.com/sony/gobreaker"
)

type Result struct {
	V interface{}
	E error
}

func (r *Result) ToHttpStatus() int {
	switch {
	case errors.Is(r.E, ErrBadRequest):
		return http.StatusBadRequest
	case errors.Is(r.E, ErrConflict):
		return http.StatusConflict
	case errors.Is(r.E, ErrForbiddenAcccess):
		return http.StatusForbidden
	case errors.Is(r.E, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(r.E, context.DeadlineExceeded):
		return http.StatusGatewayTimeout
	case errors.Is(r.E, gobreaker.ErrOpenState):
		return http.StatusServiceUnavailable
	case errors.Is(r.E, gobreaker.ErrTooManyRequests):
		return http.StatusTooManyRequests
	case r.E != nil:
		return http.StatusInternalServerError
	case r.V != nil:
		return http.StatusOK
	default:
		return http.StatusNoContent
	}
}
