package core

import "context"

type stateAdapter interface {
	Has(ctx context.Context, key string) (ok bool, err error)
	Get(ctx context.Context, key string, dest interface{}) (ok bool, err error)
	Set(ctx context.Context, key string, value interface{}) error
}
