package core

import (
	"context"

	"github.com/google/uuid"
)

type queryRepositoryAdapter interface {
	Close()
	SetModel(model Entitier, tableName string)
	Find(ctx context.Context, id uuid.UUID, dest Entitier) (err error)
	Any(ctx context.Context) (ok bool, err error)
	AnyWithFilter(ctx context.Context, query interface{}, args interface{}) (ok bool, err error)
	Count(ctx context.Context) (count int64, err error)
	CountWithFilter(ctx context.Context, query interface{}, args interface{}) (count int64, err error)
	List(ctx context.Context, dest interface{}) (err error)
	ListWithFilter(ctx context.Context, query interface{}, args interface{}, dest interface{}) (err error)
}
