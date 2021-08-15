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

func (m *validationMiddleware) Run(ctx context.Context, request core.Request) (ok bool, err error) {
	if err = m.validate.Struct(request); err != nil {
		return false, err
	}
	return true, nil
}
