package core

import (
	"context"

	"github.com/google/uuid"
)

type QueryRepositoryAdapter interface {
	SetModel(model Entitier)
	Find(ctx context.Context, dest Entitier, id uuid.UUID) (ok bool, err error)
	Any(ctx context.Context) (ok bool, err error)
	AnyWithFilter(ctx context.Context, query interface{}, args interface{}) (ok bool, err error)
	Count(ctx context.Context) (count int64, err error)
	CountWithFilter(ctx context.Context, query interface{}, args interface{}) (count int64, err error)
	List(ctx context.Context, dest []Entitier) error
	ListWithFilter(ctx context.Context, dest []Entitier, query interface{}, args interface{}) error
}
