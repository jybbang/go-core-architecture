package core

import "context"

type KV struct {
	K string
	V interface{}
}

type stateAdapter interface {
	OnCircuitOpen()
	Close()
	Has(ctx context.Context, key string) bool
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}) error
	BatchSet(ctx context.Context, kvs []KV) error
	Delete(ctx context.Context, key string) error
}
