package core

import "context"

type stateAdapter interface {
	Close()
	Has(ctx context.Context, key string) (ok bool, err error)
	Get(ctx context.Context, key string, dest interface{}) (err error)
	Set(ctx context.Context, key string, value interface{}) error
	Delete(ctx context.Context, key string) error
}
