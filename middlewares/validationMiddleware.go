package middlewares

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/jybbang/go-core-architecture/core"
)

type validationMiddleware struct {
	core.Middleware
	validate *validator.Validate
}

func NewValidationMiddleware() *validationMiddleware {
	return &validationMiddleware{
		validate: validator.New(),
	}
}

func (m *validationMiddleware) Run(ctx context.Context, request core.Request) core.Result {
	if err := m.validate.Struct(request); err != nil {
		return core.Result{E: core.ErrBadRequest}
	}
	return m.Next()
}
