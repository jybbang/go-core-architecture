package core

import "context"

type stateAdapter interface {
	Close()
	Has(ctx context.Context, key string) bool
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}) error
	Delete(ctx context.Context, key string) error
}
